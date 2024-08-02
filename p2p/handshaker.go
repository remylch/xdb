package p2p

import "errors"

// ErrInvalidHandshake Err invalid handshake error is returned if the handshake between two nodes fails.
var ErrInvalidHandshake = errors.New("invalid handshake")

// HandshakeFunc is the function that handles the handshake between two nodes.
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(node Peer) error {
	return nil
}
