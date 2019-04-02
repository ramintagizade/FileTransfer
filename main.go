package main

import (
	"flag"
	"fmt"
	//"github.com/ramintagizade/FileTransfer/server"
	"./server"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Error: , File location is required ")
		return
	}
	filePath := flag.Args()[0]
	server.Run(filePath)
}
