package store

import (
	"bytes"
	"testing"
)

func TestCompressor(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{"Small Data", []byte("Hello, World!")},
		{"Larger Data", bytes.Repeat([]byte("abcdefghijklmnopqrstuvwxyz"), 100)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			compressedData, compressedSize, err := CompressData(tc.data)

			if err != nil {
				t.Fatalf("Error compressing data: %v", err)
			}

			// Check if > 100 because in some compression algorithm if there is too few data, the compressed version is bigger than the uncompressed
			if len(tc.data) > 100 && int(compressedSize) > len(tc.data) {
				t.Errorf("Compressed data is larger than original data for input larger than 100 bytes: %d > %d", len(compressedData), len(tc.data))
			}

			decompressedData, err := DecompressData(compressedData)
			if err != nil {
				t.Fatalf("Error decompressing data: %v", err)
			}

			if !bytes.Equal(tc.data, decompressedData) {
				t.Errorf("Decompressed data is not equal to original data")
			}
		})
	}

}
