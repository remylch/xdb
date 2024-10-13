package p2p

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
	"xdb/store"
)

func tearDown(s *Server) {
	if err := s.store.Clear(); err != nil {
		log.Fatal(err)
	}
}

func TestClientHandshake(t *testing.T) {
	st := store.NewXDBStore(store.DefaultTestDataDir, "your-32-byte-secret-key-here!!!!")
	s1 := MakeServer(DefaultServerAddr, st)

	go s1.Start()

	time.Sleep(200 * time.Millisecond)

	require.Len(t, s1.GetConnexions(true), 0, "No client should be connected initially")

	SendTestMessage(s1, "test")

	time.Sleep(200 * time.Millisecond)

	require.Len(t, s1.GetConnexions(true), 1, "Expected 1 client to be connected after sending a message")

	tearDown(s1)
}

/* Test pass using IDE test but not make command
func TestNodeHandshake(t *testing.T) {
	s1 := MakeServer("./datadir", ":3000")
	s2 := MakeServer("./datadir2", ":4000", ":3000")

	go s1.Start()
	require.Len(t, s1.GetConnexions(false), 0, "No node should be connected initially")

	go s2.Start()
	time.Sleep(500 * time.Millisecond)

	require.Len(t, s1.GetConnexions(false), 1, "Expected 1 node to be connected to s1")
	require.Len(t, s2.GetConnexions(false), 1, "Expected 1 node to be connected to s2")

	tearDown(s1)
	tearDown(s2)
}
*/
