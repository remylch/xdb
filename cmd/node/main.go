package main

import (
	"flag"
	"log"
	"os"
	"xdb/config"
)

func main() {
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	configFilePath := flag.String("f", "", "f represent the default toml config file")
	config, err := config.LoadConfig(*configFilePath)

	if err != nil {
		panic(err)
	}

	node := NewNode(NodeOpts{
		peers:    config.Server.BootstrapNodes,
		hashKey:  config.Secret.HashKey,
		dataDir:  config.Storage.DataDir,
		nodeAddr: config.Server.NodeAddr,
		apiAddr:  config.Server.ApiAddr,
	})

	node.run()

	select {}
}
