package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":8080"
	opts := TCPTransportOptions{
		ListenAddr: listenAddr,
		ShakeHands: NOPHandshakeFunc,
		Decoder:    DefaultDecoder{},
	}
	transport := NewTCPTransport(opts)
	assert.Equal(t, transport.ListenAddr, listenAddr)
	assert.Nil(t, transport.ListenAndAccept())
}
