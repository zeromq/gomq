package zeromq

import (
	"errors"

	"github.com/zeromq/gozmtp"
)

type Socket interface {
	Recv() ([]byte, error)
	Send([]byte) error
}

type socket struct{}

var (
	ClientSocketType = zmtp.ClientSocketType
	ServerSocketType = zmtp.ServerSocketType

	ErrNotImplemented = errors.New("not implemented")
)

// NewSocket accepts a zmtp.SocketType and an endpoint, and returns a Socket
// interface that is ready to send and receive messages
func NewSocket(socktype zmtp.SocketType, endpoint string) (Socket, error) {
	return &socket{}, ErrNotImplemented
}

// NewClient accepts an endpoint and returns a new Client Socket that is
// ready to send and receive messages
func NewClient(endpoint string) (Socket, error) {
	return NewSocket(ClientSocketType, endpoint)
}

// NewServer accepts an endpoint and returns a new Server Socket that is
// ready to send and receive messages
func NewServer(endpoint string) (Socket, error) {
	return NewSocket(ServerSocketType, endpoint)
}

// Recv receives a ZeroMQ message
func (s *socket) Recv() ([]byte, error) {
	var msg []byte
	return msg, ErrNotImplemented
}

// Send sends a ZeroMQ message
func (s *socket) Send([]byte) error {
	return ErrNotImplemented
}
