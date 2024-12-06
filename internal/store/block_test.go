package store

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BlockFixture struct {
	input              []byte
	expectedBlockCount int
}

func TestBlockFunctions(t *testing.T) {
	smallInput := []byte("This is a test data string that we'll use to create and manipulate data blocks.")
	pseudoRandomInput := make([]byte, 6700)
	_, err := rand.Read(pseudoRandomInput)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	dataFixtures := []BlockFixture{
		{input: smallInput, expectedBlockCount: 1},
		{input: pseudoRandomInput, expectedBlockCount: 2},
	}

	// Test createBlocksFromBytes
	t.Run("createBlocksFromBytes", func(t *testing.T) {
		for _, fixture := range dataFixtures {
			blocks, err := createBlocksFromBytes(fixture.input)
			assert.NoError(t, err)
			assert.NotEmpty(t, blocks)

			assert.Equal(t, fixture.expectedBlockCount, len(blocks))

			for _, block := range blocks {
				assert.NotNil(t, block.header().blockId)
				assert.NotNil(t, block.header().dataID)
				assert.NotNil(t, block.header().totalIdx)
				assert.NotNil(t, block.header().compressedDataSize)
				assert.Equal(t, block.header().compressedDataSize, uint32(BlockPadding))
			}
		}
	})

	// Test readBlocksFromBytes
	t.Run("readBlocksFromBytes", func(t *testing.T) {
		for _, fixture := range dataFixtures {
			blocks, err := createBlocksFromBytes(fixture.input)
			assert.NoError(t, err)
			blocks, err = compressAll(blocks)
			assert.NoError(t, err)

			// Concatenate all blocks representing the total data extracted from file(s)
			var allBlocksData []byte
			for _, block := range blocks {
				allBlocksData = append(allBlocksData, block...)
			}

			// Read blocks from the concatenated data
			readBlocks, err := readBlocksFromBytes(allBlocksData)

			assert.NoError(t, err)
			assert.Equal(t, len(blocks), len(readBlocks))

			// Compare original blocks with read blocks
			for i := range blocks {
				// Compare everything except the compressed size in header
				assert.Equal(t, blocks[i][:24], readBlocks[i][:24], "Header mismatch (except compressed size)")
				assert.Equal(t, blocks[i][28:], readBlocks[i][28:], "Data mismatch")
			}
		}
	})

	// Test removePadding
	t.Run("removePadding", func(t *testing.T) {
		for _, fixture := range dataFixtures {
			blocks, err := createBlocksFromBytes(fixture.input)
			assert.NoError(t, err)

			unpaddedBlocks := removePadding(blocks...)
			assert.Equal(t, len(blocks), len(unpaddedBlocks))

			// Concatenate all unpadded data
			var reconstructedData []byte
			for _, block := range unpaddedBlocks {
				// Verify no padding remains at the end of each block
				data := block[BlockHeaderSize:]
				assert.False(t, bytes.HasSuffix(data, []byte{BlockPadding}), "Block still contains padding")
				reconstructedData = append(reconstructedData, data...)
			}

			// Verify the concatenated data matches the original input
			assert.Equal(t, fixture.input, reconstructedData, "Reconstructed data doesn't match original input")
		}
	})

}
