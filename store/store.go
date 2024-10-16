package store

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DefaultDataDir     = "/opt/xdb/data/"
	DefaultTestDataDir = "./test/data/"
)

type XDBStore struct {
	DataDir     string
	hashKey     string
	collections []Collection
}

func NewXDBStore(dataDir string, hashKey string) *XDBStore {
	if len(hashKey) != 32 {
		panic("the hash key should be 32 bytes long")
	}

	dataDir = strings.TrimSpace(dataDir)

	if dataDir == "" {
		dataDir = DefaultDataDir
	}

	store := &XDBStore{
		DataDir: dataDir,
		hashKey: hashKey,
	}

	if !dirExists(dataDir) {
		//Write default data dir
		if err := os.MkdirAll(store.DataDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	store.init()

	return store
}

/*
*
init permit a node to attach an existing dir as data store
*/
func (s *XDBStore) init() {
	collectionsFiles, err := os.ReadDir(s.DataDir)

	if err != nil {
		panic(fmt.Sprintf("Error reading collections directory: %v", err))
	}

	for _, file := range collectionsFiles {
		hash := file.Name()
		collectionName, err := decryptFilename(s.hashKey, hash)

		if err != nil {
			panic(err)
		}
		//TODO Init indexes from the file
		s.collections = append(s.collections, *newCollection(collectionName))
	}
}

// CreateCollection creates a new file for the collection with the given name.
func (s *XDBStore) CreateCollection(name string) {
	for _, collection := range s.collections {
		if collection.name == name {
			log.Fatalf("Collection '%s' already exists", name)
			return
		}
	}
	//TODO: create collection index files
	fullPath := s.getFullPathWithHash(name)
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating collection file: %v", err)
		return
	}
	collection := newCollection(name)
	s.collections = append(s.collections, *collection)
}

func (s *XDBStore) Has(collection string) bool {
	_, err := os.Stat(s.DataDir + "/" + s.getCollectionHash(collection))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (s *XDBStore) Get(collection string) ([]byte, error) {
	file, err := os.Open(s.DataDir + "/" + s.getCollectionHash(collection))
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
	return os.RemoveAll(s.DataDir)
}

func (s *XDBStore) Save(collection string, b []byte) (bool, error) {
	var err error

	fullPath := s.getFullPathWithHash(collection)

	if !s.fileExists(collection) {
		return false, fmt.Errorf("collection '%s' does not exist", collection)
	}

	//TODO: Should get the part with a given index to avoid decrypting the whole file
	decryptedData, err := s.Get(collection)

	if err != nil {
		return false, err
	}

	if isSameData(decryptedData, b) {
		return false, nil
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

	log.Printf("[%s : %v bytes written]", collection, len(b))

	return true, nil
}

func (s *XDBStore) fileExists(collection string) bool {
	filePath := s.getFullPathWithHash(collection)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (s *XDBStore) getCollectionHash(collection string) string {
	res, err := encryptFilename(s.hashKey, collection)
	if err != nil {
		log.Printf("error encrypting collection name: %v", err)
		return collection
	}
	return res
}

func (s *XDBStore) getFullPathWithHash(collection string) string {
	return s.DataDir + "/" + s.getCollectionHash(collection)
}

func (s *XDBStore) GetCollections() []string {
	collections := make([]string, len(s.collections))
	for i, collection := range s.collections {
		collections[i] = collection.name
	}
	return collections
}

func dirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	return info.IsDir()
}

func isSameData(a, b []byte) bool {
	return bytes.Equal(a, b)
}

// encrypt encrypts the given plaintext string using AES encryption with the provided key and a fixed IV.
func encryptFilename(key, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Use a fixed IV (initialization vector)
	iv := []byte(key[:16]) // 16 bytes for AES-128

	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	// Return the encrypted string as a base64 encoded string
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts the given ciphertext string using AES encryption with the provided key and a fixed IV.
func decryptFilename(key, ciphertext string) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Use the same fixed IV (initialization vector)
	iv := []byte(key[:16]) // 16 bytes for AES-128

	plaintext := make([]byte, len(ciphertextBytes))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertextBytes)

	return string(plaintext), nil
}
