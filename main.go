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

	hashKey := os.Getenv("HASH_KEY")

	//TODO: add validation on
	if hashKey == "" {
		panic("HASH_KEY environment variable is required")
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	node := newNode(NodeOpts{
		peers:    strings.Split(os.Getenv("BOOTSTRAP_NODES"), ","),
		hashKey:  os.Getenv("HASH_KEY"),
		dataDir:  os.Getenv("DATA_DIR"),
		nodeAddr: os.Getenv("NODE_ADDR"),
		apiAddr:  os.Getenv("API_ADDR"),
	})

	node.run()

	//s2 := p2p.MakeServer("./data/s2ddir", ":4000", ":3000")
	//s3 := p2p.MakeServer("./data/s3ddir", ":6000", ":4000", ":3000")
	//
	//servers := []*p2p.Server{s1, s2, s3}
	//
	//for _, s := range servers {
	//	go func(s *p2p.Server) {
	//		if err := s.Start(); err != nil {
	//			log.Fatalln(err)
	//		}
	//	}(s)
	//}
	//
	//time.Sleep(2 * time.Second)
	//
	//p2p.SendTestMessage(s3, "test")
	//
	//time.Sleep(2 * time.Second)
	//
	//b1 := s1.Retrieve("test")
	//b2 := s2.Retrieve("test")
	//b3 := s3.Retrieve("test")
	//
	//log.Println(bytes.Equal(b1, b2))
	//log.Println(bytes.Equal(b2, b3))
	//
	//select {}
}
