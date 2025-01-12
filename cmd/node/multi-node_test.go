package main

import (
	"bytes"
	"os"
	"testing"
	"time"
	"xdb/internal/p2p"
	"xdb/internal/shared"
)

const (
	HashKeyFixture    = "your-32-byte-secret-key-here!!!!"
	CollectionFixture = "test"
)

func TestMultiNodeDataReplication(t *testing.T) {
	n1 := NewNode(NodeOpts{
		peers:    []string{},
		hashKey:  HashKeyFixture,
		dataDir:  shared.DefaultXdbTestDataDirectory + "node1/",
		logDir:   shared.DefaultXdbTestLogDirectory + "node1/",
		nodeAddr: ":3001",
		apiAddr:  ":8081",
	})
	n1.run()
	n1.store.CreateCollection(CollectionFixture)

	n2Chan := make(chan *Node, 1)

	time.Sleep(time.Second * 1)

	go func() {
		n2 := NewNode(NodeOpts{
			peers:    []string{":3001"},
			hashKey:  HashKeyFixture,
			dataDir:  shared.DefaultXdbTestDataDirectory + "node2/",
			logDir:   shared.DefaultXdbTestLogDirectory + "node2/",
			nodeAddr: ":3002",
			apiAddr:  ":8082",
		})
		n2.run()
		n2.store.CreateCollection(CollectionFixture)
		n2Chan <- n2
	}()

	n2 := <-n2Chan

	time.Sleep(time.Second * 1)

	if len(n1.tcpServer.GetConnexions(false)) != 1 {
		t.Error("node1 should be connected to node2")
	}

	if len(n2.tcpServer.GetConnexions(false)) != 1 {
		t.Error("node2 should be connected to node1")
	}

	p2p.SendTestMessage(n1.tcpServer, CollectionFixture)

	time.Sleep(time.Second * 1)

	filesNode1, _ := os.ReadDir(shared.DefaultXdbTestDataDirectory + "node1/")

	if len(filesNode1) == 0 {
		t.Error("Expected 1 collection in node1, got ", len(filesNode1))
	}

	filesNode2, _ := os.ReadDir(shared.DefaultXdbTestDataDirectory + "node2/")

	if len(filesNode2) == 0 {
		t.Error("Expected 1 collection in node2, got ", len(filesNode2))
	}

	data1, _ := n1.store.Get(CollectionFixture)
	data2, _ := n2.store.Get(CollectionFixture)

	if !bytes.Equal(data1, data2) {
		t.Error("Data in all nodes should be equal")
	}

	n1.store.Clear()
	n2.store.Clear()
	deleteTestLogDir()
}

func deleteTestLogDir() {
	os.RemoveAll("./log")
}
