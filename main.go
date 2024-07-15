package main

import (
	"fmt"
	"log"

	"github.com/KDT2006/foreverstore-go/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("Doing some logic with the peer outside of TCPTransport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":4000",
		HandshakeFunc: func(any) error {
			return nil
		},
		Decoder: p2p.DefaultDecoder{},
		OnPeer:  OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
