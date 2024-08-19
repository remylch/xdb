package p2p

// RPC represents any message sent between two nodes.
type RPC struct {
	From    string
	Payload []byte
}
