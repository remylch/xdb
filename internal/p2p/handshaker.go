package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"xdb/internal/shared"
)

// HandshakeFunc is the function that handles the handshake between two nodes.
type HandshakeFunc func(Peer) error

type HandshakeMessage struct {
	Type shared.PeerType
}

func NOPHandshakeFunc(peer Peer) error {
	return nil
}

func ReadIncomingHandshake(peer Peer) (shared.PeerType, error) {
	messageBuf := ReadPrefixedLengthMessage(peer)

	rpc := RPC{
		From:    peer.RemoteAddr().String(),
		Payload: messageBuf,
	}

	reader := bytes.NewReader(rpc.Payload)
	var responseMsg shared.Message

	if err := gob.NewDecoder(reader).Decode(&responseMsg); err != nil {
		fmt.Println("decode error ", err, rpc.Payload)
		return "", err
	}

	if _, ok := responseMsg.Payload.(HandshakeMessage); !ok {
		return "", fmt.Errorf("invalid response message type, expected HandshakeMessage, got %T", responseMsg.Payload)
	}

	return responseMsg.Payload.(HandshakeMessage).Type, nil
}

func DefaultHandshake(peer Peer) error {
	handshakeMsg := shared.Message{
		Payload: HandshakeMessage{Type: shared.NodePeer},
	}

	lengthBuf, messageBytes, err := PrefixedLengthMessage(handshakeMsg)

	if err != nil {
		return err
	}

	if _, err := peer.Write(append(lengthBuf, messageBytes...)); err != nil {
		return err
	}

	peerType, err := ReadIncomingHandshake(peer)

	if err != nil {
		return err
	}

	if err := peer.DefineType(peerType); err != nil {
		fmt.Println("DefineType error ", err)
		return err
	}

	return nil
}
