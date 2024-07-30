package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/KDT2006/foreverstore-go/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey:            newEncryptionKey(),
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")
	// s3 := makeServer(":5000", ":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second)

	go s2.Start()
	time.Sleep(time.Second)

	// go s3.Start()
	time.Sleep(time.Second)

	// for i := 0; i < 20; i++ {
	key := "picture_1.jpg"
	data := bytes.NewReader([]byte("This is a big file!"))
	s2.Store(key, data)
	time.Sleep(time.Millisecond * 5)

	if err := s2.store.Delete(s2.ID, key); err != nil {
		log.Fatal(err)
	}

	r, err := s2.Get(key)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
	// }
}
