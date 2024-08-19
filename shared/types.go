package shared

type Message struct {
	Payload any
}

type PeerType string

// Overrided base types
type ByteSlice []byte
type String string
