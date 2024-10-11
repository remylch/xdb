package main

import (
	"xdb/api"
	"xdb/p2p"
	"xdb/store"
)

// A Node represents an instance of the db running on the network
type Node struct {
	httpServer *api.NodeHttpServer
	tcpServer  *p2p.Server
}

type NodeOpts struct {
	peers    []string
	hashKey  string
	dataDir  string
	nodeAddr string
	apiAddr  string
}

func newNode(opts NodeOpts) *Node {
	s := store.NewXDBStore(opts.dataDir, opts.hashKey)
	s1 := p2p.MakeServer(opts.nodeAddr, s, opts.peers...)
	httpServer := api.NewHttpServer(s, opts.apiAddr)

	return &Node{
		httpServer: httpServer,
		tcpServer:  s1,
	}
}

func (n *Node) run() {
	go n.httpServer.Start()
	n.tcpServer.Start()
	select {}
}
