package store

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func tearDown(s *XDBStore) {
	s.Clear()
}

// FIXME
func TestXDBStore(t *testing.T) {
	collection := "test"
	s := NewXDBStore(DefaultTestDataDir, "your-32-byte-secret-key-here!!!!")

	s.CreateCollection(collection)

	isCollectionStored := s.Has(collection)

	if !isCollectionStored {
		t.Error("collection should have it's directory in the store datadir")
	}

	input := []byte("hello")

	dataChanged, err := s.Save(collection, input)

	if err != nil {
		t.Error(err)
	}

	if dataChanged == false {
		t.Error("data should be changed")
	}

	dataChanged, err = s.Save(collection, input)

	if err != nil {
		t.Error(err)
	}

	if dataChanged == true {
		t.Error("data should not be changed")
	}

	data, err := s.Get(collection)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, input) {
		t.Error("data retrieved should be equal to input")
	}

	tearDown(s)
}

func TestStoreInitialization(t *testing.T) {
	data := []byte("hello")

	//Given the datadir and an existing collection file
	err := os.MkdirAll(DefaultTestDataDir, 0755)
	require.NoError(t, err, "failed to create data directory")
	err = os.WriteFile(DefaultTestDataDir+"/UVZAFg==", data, 0644)
	require.NoError(t, err, "failed to write data file")

	collection := "test"
	s := NewXDBStore(DefaultTestDataDir, "your-32-byte-secret-key-here!!!!")

	require.Len(t, s.collections, 1, "store should be initialized with one collection")
	require.Equal(t, s.collections[0].name, collection, "collection name should be equal to the one created in the test")

	tearDown(s)
}

func TestCryptoFilename(t *testing.T) {
	baseFileName := "test"
	s := NewXDBStore(DefaultTestDataDir, "your-32-byte-secret-key-here!!!!")

	encryptedFileName, err := encryptFilename(s.hashKey, baseFileName)

	if err != nil {
		t.Error(err)
	}

	res, err := decryptFilename(s.hashKey, encryptedFileName)

	if err != nil {
		t.Error(err)
	}

	if baseFileName != res {
		t.Errorf("decrypted filename is not equal to the original filename has: %s, expected: %s", res, baseFileName)
	}

	fmt.Println(res, baseFileName)

	tearDown(s)
}
