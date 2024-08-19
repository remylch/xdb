package p2p

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"net"
	"xdb/shared"
)

// message sent from a client to one of the servers
func SendTestMessage(s *Server, collection string) {
	go func() {
		conn, err := net.Dial("tcp", s.Transport.Addr())
		if err != nil {
			log.Fatal(err)
			return
		}
		defer conn.Close()

		handshakeMsg := shared.Message{
			Payload: HandshakeMessage{
				Type: shared.ClientPeer,
			},
		}

		msg := shared.Message{
			Payload: MessageStoreFile{
				Collection: collection,
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

		if _, err := conn.Write(append(lengthBuf, messageBytes...)); err != nil {
			log.Fatal(err)
		}

		select {}
	}()
}

func MakeServer(dataDir, listenAddress string, nodes ...string) *Server {
	tcpOpts := TCPTransportOptions{
		ListenAddr: listenAddress,
		ShakeHands: DefaultHandshake,
	}
	tcpTransport := NewTCPTransport(tcpOpts)

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
