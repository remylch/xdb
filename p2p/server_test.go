package p2p

import (
	"testing"
	"time"
	"xdb/store"
)

func TestSaveAndRetrieveCollectionData(t *testing.T) {
	s := store.NewXDBStore(store.DefaultTestDataDir, "your-32-byte-secret-key-here!!!!")
	collection := "test"
	s.CreateCollection(collection)
	s1 := MakeServer(DefaultServerAddr, s)
	go s1.Start()

	var data []byte

	time.Sleep(200 * time.Millisecond)

	data = s1.Retrieve(collection)

	if len(data) > 0 {
		t.Errorf("should not have data initially but have %s ", data)
	}

	SendTestMessage(s1, collection)

	time.Sleep(500 * time.Millisecond)

	data = s1.Retrieve(collection)

	if len(data) == 0 {
		t.Errorf("data should exists but have : %s ", data)
	}

	tearDown(s1)
}
