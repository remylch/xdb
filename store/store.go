package store

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
)

type Store interface {
	Save(string, []byte) error
}

type XDBStore struct {
	DefaultDataDir string
	hashKey        string
}

func NewXDBStore(dataDir string) *XDBStore {
	store := &XDBStore{
		DefaultDataDir: dataDir,
		//HashKey should be 32 bytes long
		hashKey: os.Getenv("HASH_KEY"),
	}

	//Write default data dir
	if err := os.MkdirAll(store.DefaultDataDir, os.ModePerm); err != nil {
		panic(err)
	}

	return store
}

func (s *XDBStore) Has(collection string) bool {
	_, err := os.Stat(s.DefaultDataDir + "/" + getCollectionHash(collection))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (s *XDBStore) Get(collection string) ([]byte, error) {
	file, err := os.Open(s.DefaultDataDir + "/" + getCollectionHash(collection))
	defer file.Close()

	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	data := make([]byte, fileSize)

	if _, err = file.Read(data); err != nil {
		return nil, err
	}

	data, err = Decrypt(s.hashKey, data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *XDBStore) Clear() error {
	return os.RemoveAll(s.DefaultDataDir)
}

func (s *XDBStore) Save(collection string, b []byte) (bool, error) {
	var err error
	fullPath := s.getFullPathWithHash(collection)

	if s.fileExists(collection) {
		decryptedData, err := s.Get(collection)

		if err != nil {
			return false, err
		}

		if isSameData(decryptedData, b) {
			return false, nil
		}

	} else {
		if _, err := os.Create(fullPath); err != nil {
			return false, err
		}
	}

	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()

	if err != nil {
		return false, err
	}

	b, err = Encrypt(s.hashKey, b)

	if err != nil {
		return false, err
	}

	if _, err := file.Write(b); err != nil {
		return false, fmt.Errorf("error writing to file: %s", err)
	}

	log.Printf("written %s to the disk", b)

	return true, nil
}

func (s *XDBStore) fileExists(collection string) bool {
	filePath := s.getFullPathWithHash(collection)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func getCollectionHash(collection string) string {
	hash := sha1.Sum([]byte(collection))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

func (s *XDBStore) getFullPathWithHash(collection string) string {
	return s.DefaultDataDir + "/" + getCollectionHash(collection)
}

func isSameData(a, b []byte) bool {
	return bytes.Equal(a, b)
}
