package store

import (
	"bytes"
	"testing"
)

func tearDown(s *XDBStore) {
	s.Clear()
}

func TestXDBStore(t *testing.T) {
	collection := "test"
	s := NewXDBStore("./specificDataDir")
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
