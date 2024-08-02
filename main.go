package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
	"xdb/p2p"
)

func sendTestMessage(s *Server, msg string) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(Message(msg)); err != nil {
		return err
	}

	rpc := p2p.RPC{
		Payload: buffer.Bytes(),
	}

	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	for _, peer := range s.peers {
		if err := peer.Send(rpc.Payload); err != nil {
			fmt.Printf("Error sending to peer %s: %v\n", peer.RemoteAddr(), err)
		} else {
			fmt.Printf("Message sent to peer %s\n", peer.RemoteAddr())
		}
	}

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
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	s1 := makeServer("./data/s1ddir", ":3000", ":4000", ":5000")
	s2 := makeServer("./data/s2ddir", ":4000", ":3000", ":5000")
	s3 := makeServer("./data/s3ddir", ":5000", ":3000", ":4000")

	go s1.Start()
	go s2.Start()
	go s3.Start()
	time.Sleep(3 * time.Second)

	// Test sending a message from s1 to s2
	err := sendTestMessage(s1, "Hello from s1")
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	}

	select {}
}
