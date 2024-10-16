package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var key = "your-32-byte-secret-key-here!!!!"

func TestEncrypt(t *testing.T) {
	data := []byte("hello world")
	encrypted, err := Encrypt(key, data)

	if err != nil {
		t.Error(err)
	}

	if len(encrypted) == 0 {
		t.Error("encrypted data is empty")
	}

	if bytes.Equal(encrypted, data) {
		t.Error("encrypted data should not be equal to original data")
	}

	decrypted, err := Decrypt(key, encrypted)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(decrypted, data) {
		t.Error("decrypted data is not equal to original data")
	}

	assert.Equal(t, data, decrypted, "decrypted data is not equal to original data")
}
