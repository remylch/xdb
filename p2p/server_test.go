package p2p

import (
	"os"
	"testing"
	"time"
)

func TestSaveAndRetrieveCollectionData(t *testing.T) {
	os.Setenv("HASH_KEY", "your-32-byte-secret-key-here!!!!")
	collection := "test"
	s1 := MakeServer("./datadir", ":3000")
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
