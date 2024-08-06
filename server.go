package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"xdb/p2p"
	"xdb/store"

	"github.com/gorilla/websocket"
)

type ServerOpts struct {
	DataDir        string
	Transport      p2p.Transport
	BootstrapNodes []string
}

type Server struct {
	ServerOpts
	peerLock   sync.Mutex
	peers      map[string]p2p.Peer
	store      *store.XDBStore
	quitch     chan struct{}
	httpServer *http.Server
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		store:      store.NewXDBStore(opts.DataDir),
		quitch:     make(chan struct{}),
		peers:      make(map[string]p2p.Peer),
		httpServer: &http.Server{
			Addr:    opts.Transport.Addr(),
			Handler: http.DefaultServeMux,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("[%s] Starting server", s.Transport.Addr())
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	go func() {
		http.HandleFunc("/ws", s.handleWebSocket)
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	s.bootstrapNetwork()
	s.loop()
	return nil
}

func (s *Server) Close() {
	close(s.quitch)
	//Stop http server
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}
}

func (s *Server) GetPeerGraph() string {
	var peerServerAddresses []string
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	for _, peer := range s.peers {
		peerServerAddresses = append(peerServerAddresses, peer.RemoteAddr().String())
	}
	strPeersServer := strings.Join(peerServerAddresses, ", ")
	return fmt.Sprintf("\n-----------\n[%s] : \nServers: %s\n-----------\n", s.Transport.Addr(), strPeersServer)
}

func (s *Server) OnPeer(peer p2p.Peer) {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	peerAddr := peer.RemoteAddr().String()
	s.peers[peerAddr] = peer
	log.Printf("[%s] New peer server connected: %s", time.Now().Format(time.RFC3339), peerAddr)
	log.Printf("[%s] Total peers: %d", time.Now().Format(time.RFC3339), len(s.peers))
}

func (s *Server) Store(collection string, r io.Reader) error {
	fileBuffer := new(bytes.Buffer)
	if _, err := io.Copy(fileBuffer, r); err != nil {
		return err
	}

	hasSavedNewData, err := s.store.Save(collection, fileBuffer.Bytes())

	if err != nil {
		return err
	}

	if !hasSavedNewData {
		return nil
	}

	return nil
}

func (s *Server) Retrieve(collection string) ([]byte, error) {
	data, err := s.store.Get(collection)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Server) bootstrapNetwork() {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			maxRetries := 5
			for i := 0; i < maxRetries; i++ {
				log.Printf("[%s] attempting to connect with remote: %s (attempt %d/%d)\n", s.Transport.Addr(), addr, i+1, maxRetries)
				if err := s.Transport.Dial(addr); err != nil {
					log.Printf("Error dialing bootstrap node %s: %v", addr, err)
					time.Sleep(2 * time.Second)
				} else {
					log.Printf("Successfully connected to bootstrap node %s", addr)
					return
				}
			}
			log.Printf("Failed to connect to bootstrap node %s after %d attempts", addr, maxRetries)
		}(addr)
	}
}

func (s *Server) loop() {
	defer func() {
		log.Println("Server stopped")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			fmt.Printf("RPC received from peer %s \n", rpc.From)
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("decoding error: ", err)
			}

			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("handle message error: ", err)
			}
		case <-s.quitch:
			return
		}
	}
}

/*
handleMessage is the main function to handle the message from the peer or unknown source
- When a message is received from a known peer, it is stored in the store
- When a message is received from an unknown source, it is stored in the store and broadcasted to all the peers
*/
func (s *Server) handleMessage(from string, msg *Message) error {
	s.peerLock.Lock()
	_, isPeer := s.peers[from]
	s.peerLock.Unlock()

	if err := s.Store(msg.Collection, bytes.NewReader(msg.Data)); err != nil {
		log.Printf("[%s] Error storing peer message: %v", time.Now().Format(time.RFC3339), err)
		return err
	}

	if !isPeer {
		s.broadcast(msg)
	}

	log.Printf("[%s] Message handled successfully", time.Now().Format(time.RFC3339))
	return nil
}

func (s *Server) broadcast(msg *Message) {
	log.Printf("[%s] Broadcasting message %s to peers", time.Now().Format(time.RFC3339), msg)

	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(msg)

	for _, peer := range s.peers {
		if err := peer.Send(buffer.Bytes()); err != nil {
			fmt.Printf("Error sending to peer %s: %v\n", peer.RemoteAddr(), err)
			//TODO: add retry logic later
		}
	}
}

func (s *Server) handleClientRequest(from string, msg *Message) error {
	var response []byte
	var err error

	switch msg.Operation {
	case Operation(OPERATION_READ):
		fmt.Println("write msg : ", msg)
		response, err = s.Retrieve(msg.Collection)
	case Operation(OPERATION_WRITE):
		fmt.Println("write msg : ", msg)
		err = s.Store(msg.Collection, bytes.NewReader(msg.Data))
		if err != nil {
			return err
		}
		//s.broadcast(peerMsg)
		return nil
	default:
		err = fmt.Errorf("unknown operation: %s", msg.Operation)
	}

	if err != nil {
		response = []byte(err.Error())
	}

	// Send response back to the client
	rpc := p2p.RPC{
		Payload: response,
	}

	fmt.Println("response : ", rpc)

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(&rpc); err != nil {
		return err
	}

	s.peerLock.Lock()
	peer, ok := s.peers[from]
	s.peerLock.Unlock()

	if !ok {
		return fmt.Errorf("peer not found: %s", from)
	}

	return peer.Send(buf.Bytes())
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for this example
	},
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}

		log.Println("WebSocket message received:", messageType, string(p))

		var msg Message
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		if msg.Collection != "" {
			response, err := s.handleWebSocketMessage(msg)
			if err != nil {
				log.Println("Handle websocket message error:", err)
				response = []byte(err.Error())
			}
			if err := conn.WriteMessage(messageType, response); err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
		}
	}
}

func (s *Server) handleWebSocketMessage(msg Message) ([]byte, error) {
	switch msg.Operation {
	case "write":
		if err := s.Store(msg.Collection, bytes.NewReader(msg.Data)); err != nil {
			return nil, fmt.Errorf("error storing WebSocket message: %v", err)
		}
		s.broadcast(&msg)
		return []byte("Write operation successful"), nil
	case "read":
		data, err := s.Retrieve(msg.Collection)
		if err != nil {
			return nil, fmt.Errorf("error retrieving data: %v", err)
		}
		return data, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", msg.Operation)
	}
}
