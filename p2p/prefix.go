package p2p

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"xdb/shared"
)

func PrefixedLengthMessage(msg shared.Message) (shared.ByteSlice, shared.ByteSlice, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(msg); err != nil {
		return nil, nil, fmt.Errorf("error encoding message: %v", err)
	}

	messageBytes := buffer.Bytes()
	length := uint32(len(messageBytes))
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, length)

	// Encode the handshake message
	if err := encoder.Encode(msg); err != nil {
		return nil, nil, fmt.Errorf("failed to encode handshake message: %v", err)
	}

	return lengthBuf, messageBytes, nil
}

func ReadPrefixedLengthMessage(conn net.Conn) shared.ByteSlice {
	// Read the length prefix (4 bytes)
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lengthBuf); err != nil {
		if err == io.EOF {
			return nil
		}
		log.Printf("[%s] error reading length from connection: %s \n", conn.RemoteAddr(), err)
		return nil
	}

	// Convert lengthBuf to an integer
	length := binary.BigEndian.Uint32(lengthBuf)

	if length == 0 {
		return nil
	}

	// Read the message based on the length
	messageBuf := make([]byte, length)
	if _, err := io.ReadFull(conn, messageBuf); err != nil {
		if err == io.EOF {
			return nil
		}
		log.Printf("[%s] error reading message from connection: %s \n", conn.RemoteAddr(), err)
		return nil
	}

	//fmt.Printf("message : len %v , buf: %s \n", length, messageBuf)

	return messageBuf

}
