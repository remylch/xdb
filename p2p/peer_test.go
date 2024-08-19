package p2p

import (
	"github.com/stretchr/testify/require"
	"testing"
	"xdb/shared"
)

func TestPeerDefinition(t *testing.T) {
	peer := NewTCPPeer(nil, false)
	require.NoError(t, peer.DefineType(shared.ClientPeer), "Should be able to define beer type when no type has been defined")
	require.Error(t, peer.DefineType(shared.NodePeer), "Should not be able to update peer type when it's already defined earlier")
}
