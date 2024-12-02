package store

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultDataBlockManager_ReadDataBlock(t *testing.T) {
	smallInput, _ := createBlocksFromBytes([]byte("small input"))
	smallInput = removePadding(smallInput...)

	pseudoRandomInput := make([]byte, 6700)
	_, err := rand.Read(pseudoRandomInput)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	sixThousandBytesBlocks, _ := createBlocksFromBytes(pseudoRandomInput)
	sixThousandBytesBlocks = removePadding(sixThousandBytesBlocks...)

	tests := []struct {
		name       string
		dataFile   string
		dataBlocks []DataBlock
		want       []DataBlock
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "One data block in the file",
			dataFile:   "data-1",
			dataBlocks: smallInput,
			want:       smallInput,
		},
		{
			name:       "Two data blocks in the file",
			dataFile:   "data-1",
			dataBlocks: sixThousandBytesBlocks,
			want:       sixThousandBytesBlocks,
		},
	}
	for _, tt := range tests {
		s := NewXDBStore(DefaultTestDataDir, "your-32-byte-secret-key-here!!!!")
		collection := "test"
		s.CreateCollection(collection)

		t.Run(tt.name, func(t *testing.T) {
			fullFilePath := s.getFullPathWithHash(collection) + "/" + tt.dataFile
			m := &DefaultDataBlockManager{}

			if err := m.WriteDataBlock(fullFilePath, tt.dataBlocks); err != nil {
				t.Errorf("Error writing data block: %v", err)
				return
			}

			got, err := m.ReadDataBlock(fullFilePath)

			if err != nil {
				t.Errorf("Error reading data block: %v", err)
				return
			}

			for i := range tt.want {
				// Compare everything except the compressed size in header
				if !assert.Equal(t, tt.want[i][:24], got[i][:24], "Header mismatch (except compressed size)") {
					t.Errorf("Header mismatch in block %d", i)
				}
				if !assert.Equal(t, tt.want[i][28:], got[i][28:], "Data mismatch") {
					t.Errorf("Data mismatch in block %d", i)
				}
			}
		})

		s.Clear()
	}
}

func TestDefaultDataBlockManager_mergeCorrelatedDataBlocks(t *testing.T) {

	pseudoRandomInput := make([]byte, 6700)
	_, err := rand.Read(pseudoRandomInput)
	if err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "Two block merge",
			input: pseudoRandomInput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultDataBlockManager{}

			sixThousandBytesBlocks, _ := createBlocksFromBytes(pseudoRandomInput)

			assert.Equalf(t, tt.input, m.mergeCorrelatedDataBlocks(sixThousandBytesBlocks), "mergeCorrelatedDataBlocks(%v)")
		})
	}
}
