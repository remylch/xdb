package main

import (
	"bytes"
	env "github.com/joho/godotenv"
	"log"
	"os"
	"time"
	"xdb/p2p"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	s1 := p2p.MakeServer("./data/s1ddir", ":3000")
	s2 := p2p.MakeServer("./data/s2ddir", ":4000", ":3000")
	s3 := p2p.MakeServer("./data/s3ddir", ":6000", ":4000", ":3000")

	servers := []*p2p.Server{s1, s2, s3}

	for _, s := range servers {
		go func(s *p2p.Server) {
			if err := s.Start(); err != nil {
				log.Fatalln(err)
			}
		}(s)
	}

	time.Sleep(2 * time.Second)

	p2p.SendTestMessage(s3)

	time.Sleep(2 * time.Second)

	b1, _ := s1.Retrieve("test")
	b2, _ := s2.Retrieve("test")
	b3, _ := s3.Retrieve("test")

	log.Println(bytes.Equal(b1, b2))
	log.Println(bytes.Equal(b2, b3))

	select {}
}
