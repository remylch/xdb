package main

import (
	"encoding/json"
	env "github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"xdb/p2p"
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func startHttpServer() {
	http.HandleFunc("/health", healthcheckHandler)
	_ = http.ListenAndServe(os.Getenv("HTTP_ADDR"), nil)
}

func main() {
	err := env.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	peers := strings.Split(os.Getenv("BOOTSTRAP_NODES"), ",")
	s1 := p2p.MakeServer("./data/s1ddir", os.Getenv("NODE_ADDR"), peers...)

	go startHttpServer()

	s1.Start()

	select {}

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
