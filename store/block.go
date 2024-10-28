package store

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
)

const (
	DefaultDataBlockSize = 4096
	BlockPadding         = 0x00
)

type DataBlockHeader struct {
	dataID   uuid.UUID
	blockId  uint32
	totalIdx uint32
}

type DataBlock struct {
	header DataBlockHeader
	data   []byte
}

func newDataBlocks(data []byte) ([]DataBlock, error) {
	var (
		blocks       = make([]DataBlock, 0)
		err    error = nil
	)

	fmt.Println("INPUT", len(data), data)

	if data, err = CompressData(data); err != nil {
		return blocks, fmt.Errorf("error compressing data: %s", err)
	}

	dataLen := len(data)
	nbBlocks := (dataLen + DefaultDataBlockSize - 24 - 1) / (DefaultDataBlockSize - 24)
	dataID := uuid.New() // uuid = 16 bytes (128 bit)

	for i := 0; i < nbBlocks; i++ {
		dataBlock := DataBlock{
			header: DataBlockHeader{
				dataID:   dataID,
				blockId:  uint32(i),
				totalIdx: uint32(nbBlocks),
			},
			data: make([]byte, DefaultDataBlockSize),
		}

		//Fill Block header
		copy(dataBlock.data[:16], dataID[:])
		binary.LittleEndian.PutUint32(dataBlock.data[16:20], dataBlock.header.blockId)
		binary.LittleEndian.PutUint32(dataBlock.data[20:24], dataBlock.header.totalIdx)

		//Fill Block data
		start := i * (DefaultDataBlockSize - 24)
		end := start + (DefaultDataBlockSize - 24)

		if end > dataLen {
			end = dataLen
		}

		copy(dataBlock.data[24:], data[start:end])

		blocks = append(blocks, dataBlock)
	}
	return blocks, nil
}
