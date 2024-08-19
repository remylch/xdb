package p2p

import (
	"fmt"
	"net"
	"sync"
	"xdb/shared"
)

type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
	IsClient() bool
	DefineType(shared.PeerType) error
}

type TCPPeer struct {
	net.Conn
	outbound bool
	peerType shared.PeerType
	Wg       *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		peerType: shared.UndefinedPeer,
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

func (n *TCPPeer) IsClient() bool {
	return n.peerType == shared.ClientPeer
}

func (n *TCPPeer) DefineType(t shared.PeerType) error {
	if n.peerType != shared.UndefinedPeer {
		return fmt.Errorf("cannot update peer type when it has already been defined earlier")
	}
	n.peerType = t
	return nil
}
