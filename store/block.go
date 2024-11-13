package store

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	DefaultDataBlockSize = 4096
	BlockPadding         = 0x00
)

var ErrorEmptyData = errors.New("datablock creation cannot be performed with empty data")

/*
Total header size = 24 bytes
| uuid = 16 bytes (128 bit)
| uint32 = 4 bytes
*/
type DataBlockHeader struct {
	dataID   uuid.UUID
	blockId  uint32
	totalIdx uint32
}

type DataBlock []byte

func newDataBlocks(data []byte) ([]DataBlock, error) {
	var (
		blocks       = make([]DataBlock, 0)
		err    error = nil
	)

	if len(data) == 0 {
		return blocks, ErrorEmptyData
	}

	if data, err = CompressData(data); err != nil {
		return blocks, fmt.Errorf("error compressing data: %s", err)
	}

	dataLen := len(data)
	nbBlocks := (dataLen + DefaultDataBlockSize - 24 - 1) / (DefaultDataBlockSize - 24)
	dataID := uuid.New()

	for i := 0; i < nbBlocks; i++ {
		dataBlock := make(DataBlock, DefaultDataBlockSize)
		blockId := uint32(i)

		//Fill Block header
		copy(dataBlock[:16], dataID[:])
		binary.LittleEndian.PutUint32(dataBlock[16:20], blockId)
		binary.LittleEndian.PutUint32(dataBlock[20:24], uint32(nbBlocks))

		//Fill Block data
		start := i * (DefaultDataBlockSize - 24)
		end := start + (DefaultDataBlockSize - 24)

		if end > dataLen {
			end = dataLen
		}

		copy(dataBlock[24:], data[start:end])

		blocks = append(blocks, dataBlock)
	}
	return blocks, nil
}

func (db DataBlock) header() DataBlockHeader {
	dataID, _ := uuid.FromBytes(db[:16])
	blockId := binary.LittleEndian.Uint32(db[16:20])
	nbBlocks := binary.LittleEndian.Uint32(db[20:24])

	return DataBlockHeader{
		dataID:   dataID,
		blockId:  blockId,
		totalIdx: nbBlocks,
	}
}
