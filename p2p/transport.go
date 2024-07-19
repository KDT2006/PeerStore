package p2p

import "net"

// Perr is an interface that represents the remote node.
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is anything that handles the communication
// between the nodes in the network. This can be of the form
// (TCP, UDP, websockets, ...)
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
