package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
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
	fmt.Printf("[%s] Starting server", s.Transport.Addr())
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

func (s *Server) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	peerAddr := peer.RemoteAddr().String()
	s.peers[peerAddr] = peer
	log.Printf("[%s] New peer connected: %s", time.Now().Format(time.RFC3339), peerAddr)
	log.Printf("[%s] Total peers: %d", time.Now().Format(time.RFC3339), len(s.peers))
	return nil
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
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("handle message error: ", err)
			}
		case <-s.quitch:
			return
		}
	}
}

func (s *Server) handleMessage(from string, msg *Message) error {
	log.Printf("[%s] Received message from %s: %s", time.Now().Format(time.RFC3339), from, string(*msg))

	if err := s.Store("test", bytes.NewReader(*msg)); err != nil {
		log.Printf("[%s] Error storing message: %v", time.Now().Format(time.RFC3339), err)
		return err
	}

	if from == "" {
		log.Printf("[%s] Broadcasting message to peers", time.Now().Format(time.RFC3339))
		return s.broadcast(msg)
	}

	log.Printf("[%s] Message handled successfully", time.Now().Format(time.RFC3339))
	return nil
}

func (s *Server) broadcast(msg *Message) error {
	fmt.Println("broadcasting message to peers ", msg)

	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	for _, peer := range s.peers {
		if err := peer.Send(*msg); err != nil {
			fmt.Printf("Error sending to peer %s: %v\n", peer.RemoteAddr(), err)
		}
	}

	return nil
}
