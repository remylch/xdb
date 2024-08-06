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
		From:    "client_source",
		Payload: buffer.Bytes(),
	}

	// Simulate receiving a message from an unknown source
	s.Transport.(*p2p.TCPTransport).HandleRPC(rpc)

	return nil
}

func makeServer(dataDir, listenAddress string, nodes ...string) *Server {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr: listenAddress,
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

// Base server is running on :6789 port
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	s1 := makeServer("./data/s1ddir", ":6789")

	s2 := makeServer("./data/s2ddir", ":4000", ":6789")
	s3 := makeServer("./data/s3ddir", ":6000", ":4000", ":6789")

	servers := []*Server{s1, s2, s3}

	for _, s := range servers {
		go func(s *Server) {
			if err := s.Start(); err != nil {
				log.Fatalln(err)
			}
		}(s)
	}

	time.Sleep(3 * time.Second)

	err = sendTestMessage(s3, Message{
		Collection: "test",
		Data:       []byte("Hello from unknown source"),
		Operation:  Operation(OPERATION_WRITE),
	})

	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	}

	select {}
}
