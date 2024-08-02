package p2p

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(addr string) error
	Addr() string
}
