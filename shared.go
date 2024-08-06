package main

// TODO: Change to /opt/data/xdb/
var DEFAULT_DATA_DIR = "./data"

var (
	OPERATION_READ  = "read"
	OPERATION_WRITE = "write"
)

type Operation string

type Message struct {
	Collection string
	Data       []byte
	Operation  Operation
}
