package store

import (
	"bytes"
	"github.com/andybalholm/brotli"
)

func CompressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := brotli.NewWriter(&buf)
	_, err := w.Write(data)

	if err != nil {
		return nil, err
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecompressData(compressedData []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(compressedData))
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
