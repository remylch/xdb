package main

// TODO: Change to /opt/data/xdb/
var DEFAULT_DATA_DIR = "./data"

type Message struct {
	Collection string
	Data       []byte
}
