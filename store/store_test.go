package store

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func tearDown(s *XDBStore) {
	s.Clear()
}

func TestXDBStore(t *testing.T) {
	collection := "test"
	s := NewXDBStore("./specificDataDir", "your-32-byte-secret-key-here!!!!")
	input := []byte("hello")
	dataChanged, err := s.Save(collection, input)

	if err != nil {
		t.Error(err)
	}

	if dataChanged == false {
		t.Error("data should be changed")
	}

	if !s.fileExists(collection) {
		t.Error("file should exists")
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
	datadir := "./specificDataDir"
	data := []byte("hello")

	//Given the datadir and an existing collection file
	err := os.MkdirAll(datadir, 0755)
	require.NoError(t, err, "failed to create data directory")
	err = os.WriteFile(datadir+"/UVZAFg==", data, 0644)
	require.NoError(t, err, "failed to write data file")

	collection := "test"
	s := NewXDBStore(datadir, "your-32-byte-secret-key-here!!!!")

	require.Len(t, s.collections, 1, "store should be initialized with one collection")
	require.Equal(t, s.collections[0].name, collection, "collection name should be equal to the one created in the test")

	tearDown(s)
}
