package api

import (
	"encoding/json"
	"log"
	"net/http"
	"xdb/store"
)

const (
	DefaultAPIAddr = ":8080"
)

type NodeHttpServer struct {
	store *store.XDBStore
	addr  string
}

func NewHttpServer(store *store.XDBStore, addr string) *NodeHttpServer {
	if addr == "" {
		addr = DefaultAPIAddr
	}

	return &NodeHttpServer{
		store: store,
		addr:  addr,
	}
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func getCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"collections": {"users", "products"}})
}

func (s *NodeHttpServer) Start() error {
	http.HandleFunc("/health", healthcheckHandler)
	http.HandleFunc("/collections", getCollectionsHandler)
	log.Println("API listening on ", s.addr)
	return http.ListenAndServe(s.addr, nil)
}
