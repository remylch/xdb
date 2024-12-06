package store

import (
	"bytes"
	"compress/zlib"
	"io"
)

// TODO: Maybe improve this function to use a more efficient compression algorithm
// TODO: Maybe compress only data if > X bytes (compression algorithm is more efficient on large dataset))
func CompressData(data []byte) ([]byte, uint32, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, 0, err
	}
	err = w.Close()
	if err != nil {
		return nil, 0, err
	}
	return buf.Bytes(), uint32(len(buf.Bytes())), nil
}

func DecompressData(compressedData []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
