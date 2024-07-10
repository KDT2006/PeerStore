package main

import (
	"log"

	"github.com/KDT2006/foreverstore-go/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":4000")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
