package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer struct represents a remote node over a TCP established connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn
	// if we dial a conn => outbound = true
	// if we accept and retrieve a conn => outbound = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	t.listener = ln

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("new incoming connection %+v\n", conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, false)
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP Error: %s\n", err)
		return
	}

	// read loop
	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP Error: %s\n", err)
			continue
		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("message: %+v\n", msg)
	}
}
