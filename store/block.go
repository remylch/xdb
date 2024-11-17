package store

import (
	"encoding/binary"
	"errors"

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
	blocks := make([]DataBlock, 0)

	if len(data) == 0 {
		return blocks, ErrorEmptyData
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

func (db DataBlock) data() []byte {
	return db[24:]
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

func (db DataBlock) compress() (DataBlock, error) {
	dataBlockDataCompressed, err := CompressData(db[24:])

	if err != nil {
		return nil, err
	}

	compressedDataBlock := make(DataBlock, 24+len(dataBlockDataCompressed))
	copy(compressedDataBlock[:24], db[:24])
	copy(compressedDataBlock[24:], dataBlockDataCompressed)

	return compressedDataBlock, nil
}

func (db DataBlock) decompress() (DataBlock, error) {
	decompressedData, err := DecompressData(db[24:])
	if err != nil {
		return nil, err
	}

	decompressedBlock := make(DataBlock, 24+len(decompressedData))
	copy(decompressedBlock[:24], db[:24])
	copy(decompressedBlock[24:], decompressedData)

	return decompressedBlock, nil
}

func decompressAll(dbs []DataBlock) ([]DataBlock, error) {
	decompressedDataBlocks := make([]DataBlock, 0)

	for _, db := range dbs {
		res, err := db.decompress()

		if err != nil {
			decompressedDataBlocks = append(decompressedDataBlocks, db)
			continue
		}

		decompressedDataBlocks = append(decompressedDataBlocks, res)
	}

	return decompressedDataBlocks, nil
}
