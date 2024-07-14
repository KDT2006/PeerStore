package main

import (
	"log"

	"github.com/KDT2006/foreverstore-go/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":4000",
		HandshakeFunc: func(any) error {
			return nil
		},
		Decoder: p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
