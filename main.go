package main

import (
	"fmt"
	"os"

	"github.com/ecgbeald/marston/services/server"
)

func main() {
	nodeType := os.Args[1]
	fmt.Println(nodeType)
	switch nodeType {
	case "master":
		server.GetMasterNode().Start()
	case "worker":
		server.GetWorker().Start()
	default:
		panic("invalid")
	}
}
