package p2p

import (
	"net"
	"sync"
)

type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

type TCPPeer struct {
	net.Conn
	outbound bool
	Wg       *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

func (n *TCPPeer) Close() error {
	return n.Conn.Close()
}

func (n *TCPPeer) Send(data []byte) error {
	_, err := n.Conn.Write(data)
	return err
}

func (n *TCPPeer) CloseStream() {
	n.Wg.Done()
}
