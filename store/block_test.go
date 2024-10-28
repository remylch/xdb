package store

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewDataBlocks(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectedBlocks int
		expectError    bool
	}{
		{
			name:           "Empty input",
			input:          []byte{},
			expectedBlocks: 1,
			expectError:    false,
		},
		{
			name:           "Small input",
			input:          []byte("Hello, World!"),
			expectedBlocks: 1,
			expectError:    false,
		},
		{
			name:           "Large input",
			input:          bytes.Repeat([]byte("A"), DefaultDataBlockSize*2),
			expectedBlocks: 3,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks, err := newDataBlocks(tt.input)

			fmt.Println("blocks:", blocks)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(blocks) != tt.expectedBlocks {
				t.Errorf("Expected %d blocks, but got %d", tt.expectedBlocks, len(blocks))
			}

			for i, block := range blocks {
				if block.header.blockId != uint32(i) {
					t.Errorf("Block %d: Expected blockId %d, but got %d", i, i, block.header.blockId)
				}

				if block.header.totalIdx != uint32(tt.expectedBlocks) {
					t.Errorf("Block %d: Expected totalIdx %d, but got %d", i, tt.expectedBlocks, block.header.totalIdx)
				}

				if len(block.data) != DefaultDataBlockSize {
					t.Errorf("Block %d: Expected data length %d, but got %d", i, DefaultDataBlockSize, len(block.data))
				}
			}

			// Check if the last block is padded correctly
			lastBlock := blocks[len(blocks)-1]
			lastBlockDataLen := len(bytes.TrimRight(lastBlock.data[24:], string(rune(BlockPadding))))
			if lastBlockDataLen == DefaultDataBlockSize-24 && len(blocks) > 1 {
				t.Errorf("Last block should be padded, but it's full")
			}
		})
	}
}
