package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	rpcch    chan RPC
}

type TCPTransportOptions struct {
	ListenAddr       string
	ShakeHands       HandshakeFunc
	OnPeer           func(Peer)
	OnPeerDisconnect func(string)
}

func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
		rpcch:               make(chan RPC, 1024),
	}
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

// Consume implements the transport interface, which will return a read-only channel of RPCs (incoming messages).
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.start()
	log.Println("TCP listening on ", t.ListenAddr)
	return nil
}

func (t *TCPTransport) start() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %v\n", err)
		}
		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		//Close connection gracefully
		if t.OnPeerDisconnect != nil {
			t.OnPeerDisconnect(conn.RemoteAddr().String())
		}
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.ShakeHands(peer); err != nil {
		log.Printf("Handshake failed: %v\n", err)
		return
	}

	if t.OnPeer != nil {
		t.OnPeer(peer)
	}

	for {
		messageBuf := ReadPrefixedLengthMessage(conn)

		if messageBuf.IsEmpty() {
			continue
		}

		//log.Printf("[%s] RPC RECEIVED from %s with message %s", t.Addr(), conn.RemoteAddr().String(), messageBuf)

		rpc := RPC{}
		rpc.Payload = messageBuf
		rpc.From = conn.RemoteAddr().String()

		t.rpcch <- rpc
	}
}
