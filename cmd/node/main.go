package main

import (
	"flag"
	"log"
	"os"
	"xdb/config"
)

func main() {
	configFilePath := flag.String("f", "", "f represent the default toml config file")
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	conf, err := config.LoadConfig(*configFilePath)

	log.Println("config", conf)

	if err != nil {
		panic(err)
	}

	node := NewNode(NodeOpts{
		peers:    conf.Server.BootstrapNodes,
		hashKey:  conf.Secret.HashKey,
		dataDir:  conf.Storage.DataDir,
		nodeAddr: conf.Server.NodeAddr,
		apiAddr:  conf.Server.ApiAddr,
	})

	node.run()

	select {}
}
