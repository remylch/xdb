package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
	"xdb/p2p"

	"github.com/joho/godotenv"
)

func sendTestMessage(s *Server, msg Message) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(msg); err != nil {
		return fmt.Errorf("encoding error: %v", err)
	}

	rpc := p2p.RPC{
		From:    "unknown_source",
		Payload: buffer.Bytes(),
	}

	// Simulate receiving a message from an unknown source
	s.Transport.(*p2p.TCPTransport).HandleRPC(rpc)

	return nil
}

func makeServer(dataDir, listenAddress string, nodes ...string) *Server {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr: listenAddress,
		ShakeHands: p2p.NOPHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpOpts)

	opts := ServerOpts{
		DataDir:        dataDir,
		Transport:      tcpTransport,
		BootstrapNodes: nodes,
	}

	s := NewServer(opts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	s1 := makeServer("./data/s1ddir", ":3000")
	s2 := makeServer("./data/s2ddir", ":4000", ":3000")
	s3 := makeServer("./data/s3ddir", ":6000", ":4000", ":3000")

	servers := []*Server{s1, s2, s3}

	for _, s := range servers {
		go func(s *Server) {
			if err := s.Start(); err != nil {
				log.Fatalln(err)
			}

		}(s)
	}

	time.Sleep(3 * time.Second)

	// Test sending a message from s1 to s2
	err = sendTestMessage(s3, Message{
		Collection: "test",
		Data:       []byte("Hello from unknown source"),
	})

	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	}

	time.Sleep(3 * time.Second)

	log.Println(s1.GetPeerGraph())
	log.Println(s2.GetPeerGraph())
	log.Println(s3.GetPeerGraph())

	b1, _ := s1.Retrieve("test")
	b2, _ := s2.Retrieve("test")
	b3, _ := s3.Retrieve("test")

	log.Println(bytes.Equal(b1, b2))
	log.Println(bytes.Equal(b2, b3))

	select {}
}
