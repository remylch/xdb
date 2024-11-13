package store

import (
	"fmt"
	"os"
)

// DataBlockManager defines the interface for managing data blocks in files.
type DataBlockManager interface {
	WriteDataBlock(filepath string, blocks []DataBlock) error
	ReadDataBlock(filepath string) ([]DataBlock, error)
}

// FileDataBlockManager is an implementation of DataBlockManager that uses the filesystem.
type FileDataBlockManager struct{}

// NewFileDataBlockManager creates a new FileDataBlockManager.
func NewFileDataBlockManager() *FileDataBlockManager {
	return &FileDataBlockManager{}
}

func (m *FileDataBlockManager) WriteDataBlock(filepath string, blocks []DataBlock) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	for _, block := range blocks {
		if _, err := file.Write(block); err != nil {
			return fmt.Errorf("error writing to file: %s", err)
		}
	}

	return nil
}

// TODO: Add query on datablocks, for now read all the blocks from the file
func (m *FileDataBlockManager) ReadDataBlock(filepath string) ([]DataBlock, error) {
	return make([]DataBlock, 0), nil
}
