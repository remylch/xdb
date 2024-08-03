package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
	"xdb/p2p"
	"xdb/store"
)

type ServerOpts struct {
	DataDir        string
	Transport      p2p.Transport
	BootstrapNodes []string
}

type Server struct {
	ServerOpts
	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *store.XDBStore
	quitch   chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		store:      store.NewXDBStore(opts.DataDir),
		quitch:     make(chan struct{}),
		peers:      make(map[string]p2p.Peer),
	}
}

func (s *Server) Start() error {
	log.Printf("[%s] Starting server", s.Transport.Addr())
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	s.loop()
	return nil
}

func (s *Server) Close() {
	close(s.quitch)
}

func (s *Server) GetPeerGraph() string {
	var peerAddresses []string
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	for _, peer := range s.peers {
		peerAddresses = append(peerAddresses, peer.RemoteAddr().String())
	}
	strPeers := strings.Join(peerAddresses, ", ")
	return fmt.Sprintf("\n-----------\n[%s] : %s\n-----------\n", s.Transport.Addr(), strPeers)
}

func (s *Server) OnPeer(peer p2p.Peer) {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	peerAddr := peer.RemoteAddr().String()
	s.peers[peerAddr] = peer
	log.Printf("[%s] New peer connected: %s", time.Now().Format(time.RFC3339), peerAddr)
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
		fmt.Printf("[%s] data already exists\n", collection)
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
			fmt.Println("RPC received : ", rpc)
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("decoding error: ", err)
			}

			fmt.Println("RPC MESSAGE DECODED : ", msg.Collection, msg.Data)
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
