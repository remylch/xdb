package main

import (
	"xdb/api"
	"xdb/internal/log"
	"xdb/internal/p2p"
	"xdb/internal/store"

	"github.com/google/uuid"
)

// A Node represents an instance of the db running on the network
type Node struct {
	id         uuid.UUID
	version    string
	httpServer *api.NodeHttpServer
	tcpServer  *p2p.Server
	store      *store.XDBStore
}

type NodeOpts struct {
	peers    []string
	version  string
	hashKey  string
	nodeAddr string
	apiAddr  string

	dataDir string
	logDir  string
}

func NewNode(opts NodeOpts) *Node {
	logger := log.NewDefaultLogger(opts.logDir)

	s := store.NewXDBStore(opts.dataDir, opts.hashKey)
	s1 := p2p.MakeServer(opts.nodeAddr, s, opts.peers...)
	httpServer := api.NewHttpServer(s, opts.apiAddr, s1, logger)

	return &Node{
		id:         uuid.New(),
		version:    "0.0.1",
		store:      s,
		httpServer: httpServer,
		tcpServer:  s1,
	}
}

func (n *Node) run() {
	go func() {
		err := n.httpServer.Start()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := n.tcpServer.Start()
		if err != nil {
			panic(err)
		}
	}()
}
