package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultDataBlockManager_ReadDataBlock(t *testing.T) {
	dbs, _ := newDataBlocks([]byte("Hello world"))

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
			dataBlocks: dbs,
			want:       dbs,
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

			assert.Equalf(t, tt.want, got, "ReadDataBlock(%v)", tt.want)
		})

		s.Clear()
	}
}

func TestDefaultDataBlockManager_WriteDataBlock(t *testing.T) {
	type args struct {
		filepath string
		blocks   []DataBlock
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultDataBlockManager{}
			tt.wantErr(t, m.WriteDataBlock(tt.args.filepath, tt.args.blocks), fmt.Sprintf("WriteDataBlock(%v, %v)", tt.args.filepath, tt.args.blocks))
		})
	}
}

func TestDefaultDataBlockManager_mergeCorrelatedDataBlocks(t *testing.T) {
	type args struct {
		blocks []DataBlock
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultDataBlockManager{}
			assert.Equalf(t, tt.want, m.mergeCorrelatedDataBlocks(tt.args.blocks), "mergeCorrelatedDataBlocks(%v)", tt.args.blocks)
		})
	}
}
