package store

import (
	"bytes"
	"fmt"
	"os"
	"testing"
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

	//dataChanged, err = s.Save(collection, input)
	//
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//if dataChanged == true {
	//	t.Error("data should not be changed")
	//}

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
	s := NewXDBStore(DefaultTestDataDir, "your-32-byte-secret-key-here!!!!")

	dir, err := os.ReadDir(DefaultTestDataDir)

	if err != nil {
		t.Error("store directory should exist. Error : ", err)
	}

	for _, entry := range dir {
		if entry != nil {
			t.Error("store directory should be empty. Found : ", entry.Name())
		}
	}

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
