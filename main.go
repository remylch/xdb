package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"net"
	"os"
	"time"
	"xdb/p2p"
	"xdb/shared"

	env "github.com/joho/godotenv"
)

// message sent from a client to one of the servers
func sendTestMessage(s *Server) {
	go func() {
		conn, err := net.Dial("tcp", s.Transport.Addr())
		if err != nil {
			log.Fatal(err)
			return
		}
		defer conn.Close()

		handshakeMsg := shared.Message{
			Payload: p2p.HandshakeMessage{
				Type: shared.ClientPeer,
			},
		}

		msg := shared.Message{
			Payload: MessageStoreFile{
				Collection: "test",
				Data:       []byte("Hello, World!"),
			},
		}

		lengthBufHandshake := make([]byte, 4)
		var bufferHandshake bytes.Buffer
		encoderHandshake := gob.NewEncoder(&bufferHandshake)
		if err := encoderHandshake.Encode(handshakeMsg); err != nil {
			log.Fatal(err)
		}

		messageBytesHandshake := bufferHandshake.Bytes()
		lengthHandshake := uint32(len(messageBytesHandshake))
		binary.BigEndian.PutUint32(lengthBufHandshake, lengthHandshake)

		//----------------

		lengthBuf := make([]byte, 4)
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		if err := encoder.Encode(msg); err != nil {
			log.Fatal(err)
		}

		messageBytes := buffer.Bytes()
		length := uint32(len(messageBytes))
		binary.BigEndian.PutUint32(lengthBuf, length)

		if _, err := conn.Write(append(lengthBufHandshake, messageBytesHandshake...)); err != nil {
			log.Fatal(err)
		}

		time.Sleep(3 * time.Second)

		if _, err := conn.Write(append(lengthBuf, messageBytes...)); err != nil {
			log.Fatal(err)
		}

		select {}
	}()
}

func makeServer(dataDir, listenAddress string, nodes ...string) *Server {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr: listenAddress,
		ShakeHands: p2p.DefaultHandshake,
	}
	tcpTransport := p2p.NewTCPTransport(tcpOpts)

	opts := ServerOpts{
		DataDir:        dataDir,
		Transport:      tcpTransport,
		BootstrapNodes: nodes,
	}

	s := NewServer(opts)
	tcpTransport.OnPeer = s.OnPeer
	tcpTransport.OnPeerDisconnect = s.OnPeerDisconnect
	return s
}

func main() {
	err := env.Load()
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

	time.Sleep(2 * time.Second)

	sendTestMessage(s3)

	time.Sleep(2 * time.Second)

	b1, _ := s1.Retrieve("test")
	b2, _ := s2.Retrieve("test")
	b3, _ := s3.Retrieve("test")

	log.Println(bytes.Equal(b1, b2))
	log.Println(bytes.Equal(b2, b3))

	select {}
}
