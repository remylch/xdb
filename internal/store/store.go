package store

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"xdb/internal/shared"
)

const (
	DefaultDataDir = shared.DefaultBaseXdbDirectory + "internal"
)

type XDBStore struct {
	dataBlockManager DataBlockManager
	queryExecutor    QueryExecutor

	DataDir     string
	hashKey     string
	collections []Collection

	mu sync.RWMutex
}

func NewXDBStore(dataDir string, hashKey string) *XDBStore {
	if len(hashKey) != 32 {
		panic("the hash key should be 32 bytes long")
	}

	dataDir = strings.TrimSpace(dataDir)

	if dataDir == "" {
		dataDir = DefaultDataDir
	}

	dataBlockManager := NewFileDataBlockManager()

	store := &XDBStore{
		DataDir:          dataDir,
		hashKey:          hashKey,
		dataBlockManager: dataBlockManager,
		collections:      make([]Collection, 0),
		mu:               sync.RWMutex{},
	}
	//TODO: Improve to remove cross dependencies
	queryExecutor := NewBaseExecutor(dataBlockManager, dataDir, store.getCollectionHash)

	store.queryExecutor = queryExecutor

	if !shared.DirExists(dataDir) {
		//Write default data dir
		if err := os.MkdirAll(store.DataDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	store.init()

	return store
}

// init permit a node to attach an existing dir as data store
func (s *XDBStore) init() {
	s.mu.Lock()
	defer s.mu.Unlock()

	dirEntries, err := os.ReadDir(s.DataDir)

	if err != nil {
		panic(fmt.Sprintf("Error reading collections directory: %v", err))
	}

	for _, dir := range dirEntries {
		if dir.IsDir() {
			collectionName, err := decryptFilename(s.hashKey, dir.Name())

			if err != nil {
				panic(err)
			}

			//TODO Init indexes from the file
			s.collections = append(s.collections, newCollection(collectionName))
		}
	}
}

// CreateCollection creates a new file for the collection with the given name.
func (s *XDBStore) CreateCollection(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if collection already exists
	for _, collection := range s.collections {
		if collection.Name == name {
			return fmt.Errorf("collection '%s' already exists", name)
		}
	}

	// Create collection directory and files
	fullPath := s.getFullPathWithHash(name)
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating collection file: %v", err)
	}

	if err := os.WriteFile(fullPath+"/data-1", nil, 0644); err != nil {
		return fmt.Errorf("error creating initial data file for collection %v with error : %v", name, err)
	}

	s.collections = append(s.collections, newCollection(name))
	return nil
}

func (s *XDBStore) GetCollections() []Collection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.collections
}

func (s *XDBStore) Has(collection string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := s.collectionExists(collection)
	return res
}

func (s *XDBStore) Get(query string) ([]byte, error) {
	//TODO: Parse query
	result := s.queryExecutor.Execute(ReadQuery(query))

	decryptedData, err := Decrypt(s.hashKey, result.Data)

	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %v", err)
	}

	return decryptedData, nil
}

func (s *XDBStore) Clear() error {
	return os.RemoveAll(s.DataDir)
}

func (s *XDBStore) Save(collection string, b []byte) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.collectionExists(collection) {
		return false, fmt.Errorf("collection [%s] does not exist", collection)
	}

	var (
		dataBlocks []DataBlock
		err        error = nil
	)

	b, err = Encrypt(s.hashKey, b)

	if err != nil {
		return false, err
	}

	if dataBlocks, err = createBlocksFromBytes(b); err != nil {
		return false, err
	}

	//TODO: Save the datablock inside a file (how to choose/create it ?) then assign that datablock to an index
	filepath := s.getFileToWriteIn(collection)

	if err = s.dataBlockManager.WriteDataBlock(filepath, dataBlocks); err != nil {
		return false, err
	}

	log.Printf("[%s : %v bytes written (%v blocks)]", collection, len(b), len(dataBlocks))
	return true, nil
}

func (s *XDBStore) collectionExists(collection string) bool {
	dirPath := s.getFullPathWithHash(collection)
	info, err := os.Stat(dirPath)
	return !os.IsNotExist(err) && info.IsDir()
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

// TODO: put X data blocks per file. => find last file to append to or create a new one
func (s *XDBStore) getFileToWriteIn(collection string) string {
	return s.getFullPathWithHash(collection) + "/data-1"
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
