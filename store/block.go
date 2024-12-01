package store

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	DataBlockSize   = 4096
	BlockHeaderSize = 28
	BlockPadding    = 0x00
)

var ErrorEmptyData = errors.New("datablock creation cannot be performed with empty data")
var ErrorTooSmallHeader = errors.New("data too small to contain header")

/*
Total header size = 28 bytes
| uuid = 16 bytes (128 bit)
| uint32 = 4 bytes
*/
type DataBlockHeader struct {
	dataID             uuid.UUID
	blockId            uint32
	totalIdx           uint32
	compressedDataSize uint32
}

type DataBlock []byte

func createBlocksFromBytes(data []byte) ([]DataBlock, error) {
	blocks := make([]DataBlock, 0)

	if len(data) == 0 {
		return blocks, ErrorEmptyData
	}

	dataLen := len(data)
	nbBlocks := (dataLen + DataBlockSize - BlockHeaderSize - 1) / (DataBlockSize - BlockHeaderSize)
	dataID := uuid.New()

	for i := 0; i < nbBlocks; i++ {
		dataBlock := make(DataBlock, DataBlockSize)
		blockId := uint32(i)

		//Fill Block header
		copy(dataBlock[:16], dataID[:])
		binary.LittleEndian.PutUint32(dataBlock[16:20], blockId)
		binary.LittleEndian.PutUint32(dataBlock[20:24], uint32(nbBlocks))
		binary.LittleEndian.PutUint32(dataBlock[24:BlockHeaderSize], uint32(BlockPadding))

		//Fill Block data
		start := i * (DataBlockSize - BlockHeaderSize)
		end := start + (DataBlockSize - BlockHeaderSize)

		if end > dataLen {
			end = dataLen
		}

		copy(dataBlock[BlockHeaderSize:], data[start:end])

		blocks = append(blocks, dataBlock)
	}
	return blocks, nil
}

// FIXME: Cannot read block if not compressed maybe add isCompressed boolean to header and add logic on reading ???
func readBlocksFromBytes(data []byte) ([]DataBlock, error) {
	blocks := make([]DataBlock, 0)

	if len(data) == 0 {
		return blocks, ErrorEmptyData
	}

	if len(data) < BlockHeaderSize {
		return blocks, ErrorTooSmallHeader
	}

	offset := 0
	for offset < len(data) {

		//Read header first to get compressed size
		if offset+BlockHeaderSize > len(data) {
			break
		}

		compressedDataSize := binary.LittleEndian.Uint32(data[offset+24 : offset+BlockHeaderSize])
		blockSize := BlockHeaderSize + int(compressedDataSize)

		fmt.Println("Compressed data size : ", compressedDataSize, " block size : ", blockSize)

		// If we don't have enough data for the full block, take what's left
		if offset+blockSize > len(data) {
			block := make(DataBlock, len(data)-offset)
			copy(block, data[offset:])
			blocks = append(blocks, block)
			break
		}

		block := make(DataBlock, blockSize)
		copy(block, data[offset:offset+blockSize])
		blocks = append(blocks, block)

		offset += blockSize
	}

	return blocks, nil
}

// Get the data of the DataBlock without header and padding
func (db DataBlock) data() []byte {
	return bytes.TrimRight(db[BlockHeaderSize:], string(rune(BlockPadding)))
}

// Get the header of the DataBlock (28 first bytes)
func (db DataBlock) header() DataBlockHeader {
	dataID, _ := uuid.FromBytes(db[:16])
	blockId := binary.LittleEndian.Uint32(db[16:20])
	nbBlocks := binary.LittleEndian.Uint32(db[20:24])
	compressedDataSize := binary.LittleEndian.Uint32(db[24:BlockHeaderSize])

	return DataBlockHeader{
		dataID:             dataID,
		blockId:            blockId,
		totalIdx:           nbBlocks,
		compressedDataSize: compressedDataSize,
	}
}

// Compress the data of the DataBlock but not the header
func (db DataBlock) compress() (DataBlock, error) {
	dataBlockDataCompressed, compressedSize, err := CompressData(db.data())

	if err != nil {
		return nil, err
	}

	compressedDataBlock := make(DataBlock, BlockHeaderSize+len(dataBlockDataCompressed))
	copy(compressedDataBlock[:24], db[:24])
	binary.LittleEndian.PutUint32(compressedDataBlock[24:BlockHeaderSize], compressedSize)
	copy(compressedDataBlock[BlockHeaderSize:], dataBlockDataCompressed)

	return compressedDataBlock, nil
}

// Decompress the data of the DataBlock but not the header
func (db DataBlock) decompress() (DataBlock, error) {
	decompressedData, err := DecompressData(db.data())
	if err != nil {
		return nil, err
	}

	decompressedBlock := make(DataBlock, BlockHeaderSize+len(decompressedData))
	copy(decompressedBlock[:BlockHeaderSize], db[:BlockHeaderSize])
	copy(decompressedBlock[BlockHeaderSize:], decompressedData)

	return decompressedBlock, nil
}

func compressAll(dbs []DataBlock) ([]DataBlock, error) {
	compressedDataBlocks := make([]DataBlock, 0)

	for _, db := range dbs {
		res, err := db.compress()

		if err != nil {
			compressedDataBlocks = append(compressedDataBlocks, db)
			continue
		}
		compressedDataBlocks = append(compressedDataBlocks, res)
	}

	return compressedDataBlocks, nil
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

func removePadding(blocks ...DataBlock) []DataBlock {
	dbs := make([]DataBlock, 0, len(blocks))

	for _, block := range blocks {
		// Remove padding from the data part of the block
		trimmedData := bytes.TrimRight(block.data(), string(rune(BlockPadding)))

		// Create a new block with the original header and trimmed data
		newBlock := make(DataBlock, BlockHeaderSize+len(trimmedData))
		copy(newBlock[:BlockHeaderSize], block[:BlockHeaderSize])
		copy(newBlock[BlockHeaderSize:], trimmedData)

		dbs = append(dbs, newBlock)
	}

	return dbs
}
