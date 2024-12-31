package store

import (
	"fmt"
	"os"
)

// DataBlockManager defines the interface for managing data blocks in files.
type DataBlockManager interface {
	WriteDataBlock(filepath string, blocks []DataBlock) error
	ReadDataBlock(filepath string) ([]DataBlock, error)
	mergeCorrelatedDataBlocks(blocks []DataBlock) []byte
}

// DefaultDataBlockManager is an implementation of DataBlockManager that uses the filesystem.
type DefaultDataBlockManager struct{}

// NewFileDataBlockManager creates a new DefaultDataBlockManager.
func NewFileDataBlockManager() *DefaultDataBlockManager {
	return &DefaultDataBlockManager{}
}

func (m *DefaultDataBlockManager) WriteDataBlock(filepath string, blocks []DataBlock) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	for _, block := range blocks {
		compressedBlock, err := block.compress()

		if err != nil {
			return err
		}

		if _, err := file.Write(compressedBlock); err != nil {
			return fmt.Errorf("error writing to file: %s", err)
		}
	}

	return nil
}

func (m *DefaultDataBlockManager) ReadDataBlock(filepath string) ([]DataBlock, error) {
	fileBytes, err := os.ReadFile(filepath)

	if err != nil {
		return make([]DataBlock, 0), err
	}

	dbs, err := readBlocksFromBytes(fileBytes)

	if err != nil {
		return make([]DataBlock, 0), err
	}

	dbs = removePadding(dbs...)

	decompressed, err := decompressAll(dbs)

	if err != nil {
		return make([]DataBlock, 0), err
	}

	return decompressed, nil
}

func (m *DefaultDataBlockManager) mergeCorrelatedDataBlocks(blocks []DataBlock) []byte {
	mergedData := make([]byte, 0)

	for _, block := range blocks {
		mergedData = append(mergedData, block.data()...)
	}

	return mergedData
}
