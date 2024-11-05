package main

import (
	env "github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	node := NewNode(NodeOpts{
		peers:    strings.Split(os.Getenv("BOOTSTRAP_NODES"), ","),
		hashKey:  os.Getenv("HASH_KEY"),
		dataDir:  os.Getenv("DATA_DIR"),
		nodeAddr: os.Getenv("NODE_ADDR"),
		apiAddr:  os.Getenv("API_ADDR"),
	})

	node.run()

	select {}
}
