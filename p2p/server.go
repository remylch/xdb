package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
	"xdb/shared"
	"xdb/store"
)

type ServerOpts struct {
	DataDir        string
	HashKey        string
	Transport      Transport
	BootstrapNodes []string
}

type Server struct {
	ServerOpts
	peerLock sync.Mutex
	peers    map[string]Peer
	store    *store.XDBStore
	quitch   chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		store:      store.NewXDBStore(opts.DataDir, opts.HashKey),
		quitch:     make(chan struct{}),
		peers:      make(map[string]Peer),
	}
}

func (s *Server) Start() error {
	log.Printf("[%s] Starting server... \n", s.Transport.Addr())
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

/*
TODO: get peers from each node to return the graph
*/
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

func (s *Server) OnPeer(peer Peer) {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	peerAddr := peer.RemoteAddr().String()
	s.peers[peerAddr] = peer
	log.Printf("[%s] New [IsClient: %v] peer connected: %s, total peers: %d", s.Transport.Addr(), peer.IsClient(), peerAddr, len(s.peers))
}

func (s *Server) OnPeerDisconnect(addr string) {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	delete(s.peers, addr)
	log.Printf("[%s] Peer disconnected: %s, total peers: %d", s.Transport.Addr(), addr, len(s.peers))
}

func (s *Server) Store(collection string, msg MessageStoreFile, shouldBroadcast bool) error {
	r := bytes.NewReader(msg.Data)
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

	log.Printf("Data saved in collection [%s]\n", collection)

	if shouldBroadcast {
		s.broadcast(&shared.Message{Payload: msg})
	}

	return nil
}

func (s *Server) Retrieve(collection string) []byte {
	isCollectionOnDisk := s.store.Has(collection)
	if !isCollectionOnDisk {
		//TODO: fetch it from peers if it's find, store it on the current peer disk too
		log.Printf("collection [%s] does not exist", collection)
		return nil
	}
	data, err := s.store.Get(collection)
	if err != nil {
		log.Printf("[%s] Error retrieving data from local collection [%s]: %v", s.Transport.Addr(), collection, err)
		return nil
	}
	return data
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
				} else {
					log.Printf("Successfully connected to bootstrap node %s", addr)
					return
				}
			}
			log.Printf("Failed to connect to bootstrap node %s after %d attempts", addr, maxRetries)
		}(addr)
	}
}

// loop is using streaming approach to read the rpc messages
func (s *Server) loop() {
	defer func() {
		log.Println("Server stopped")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg shared.Message

			reader := bytes.NewReader(rpc.Payload)
			if err := gob.NewDecoder(reader).Decode(&msg); err != nil {
				log.Println("rpc decoding error: ", err)
			}

			fmt.Println("RPC message received : ", msg)
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("handle rpc message error: ", err)
			}
		case <-s.quitch:
			return
		}
	}
}

func (s *Server) GetConnexions(client bool) []string {
	clients := make([]string, 0)
	for _, peer := range s.peers {
		if peer.IsClient() == client {
			clients = append(clients, peer.RemoteAddr().String())
		}
	}
	return clients
}

/*
handleMessage is the main function to handle the message from the peers (client or server)
- MessageStoreFile:
- When a message is received from a client, it is stored and broadcasted to all the server peers
- When a message is received from a server, it is only stored

- MessageGetFile:
- The data are retrieved from the peer's disk and sent to the client. If the data is not in the local disk, we find it from the other peers of the network.

- MessageHandshake:
- Handshake messages are used to confirm that the peer is connected to the network
*/
func (s *Server) handleMessage(from string, initialMsg *shared.Message) error {
	s.peerLock.Lock()
	peer, isPeer := s.peers[from]
	s.peerLock.Unlock()

	if !isPeer {
		return fmt.Errorf("peer not found : %s", from)
	}

	switch msg := initialMsg.Payload.(type) {
	case MessageStoreFile:
		if err := s.Store(msg.Collection, msg, peer.IsClient()); err != nil {
			log.Printf("[%s] Error storing peer message: %v", s.Transport.Addr(), err)
			return err
		}
	case MessageGetFile:
		return peer.Send(s.Retrieve(msg.Collection))
	case HandshakeMessage:
		return s.ConfirmHandshake(peer)
	}

	log.Printf("[%s] Message handled successfully", time.Now().Format(time.RFC3339))
	return nil
}

func (s *Server) broadcast(msg *shared.Message) {
	log.Printf("[%s] Broadcasting message %s to peers", s.Transport.Addr(), msg)

	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	var (
		lengthBuf    []byte
		messageBytes []byte
		err          error
	)

	if lengthBuf, messageBytes, err = PrefixedLengthMessage(*msg); err != nil {
		return
	}

	for _, peer := range s.peers {
		if !peer.IsClient() {
			if err := peer.Send(append(lengthBuf, messageBytes...)); err != nil {
				fmt.Printf("Error sending to peer %s: %v\n", peer.RemoteAddr(), err)
				//TODO: add retry logic later
			}
		}
	}
}

func (s *Server) ConfirmHandshake(peer Peer) error {
	msg := shared.Message{
		Payload: HandshakeMessage{
			Type: shared.NodePeer,
		},
	}

	lengthBuf, messageBytes, err := PrefixedLengthMessage(msg)

	if err != nil {
		return err
	}

	return peer.Send(append(lengthBuf, messageBytes...))
}

type MessageStoreFile struct {
	Collection string
	Data       []byte
}

type MessageGetFile struct {
	Collection string
}

var once sync.Once

func init() {
	once.Do(func() {
		gob.Register(HandshakeMessage{})
		gob.Register(MessageStoreFile{})
		gob.Register(MessageGetFile{})
	})
}
