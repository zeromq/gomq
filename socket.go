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

	NullSecurityMechanism  = zmtp.NullSecurityMechanismType
	PlainSecurityMechanism = zmtp.PlainSecurityMechanismType
	CurveSecurityMechanism = zmtp.CurveSecurityMechanismType

	ErrNotImplemented = errors.New("not implemented")
)

func NewSecurityNull() *zmtp.SecurityNull {
	return &zmtp.SecurityNull{}
}

func NewSocket(socktype zmtp.SocketType, endpoint string, mechanism zmtp.SecurityMechanism) (Socket, error) {
	return &socket{}, ErrNotImplemented
}

func NewClient(endpoint string, mechanism zmtp.SecurityMechanism) (Socket, error) {
	return NewSocket(ClientSocketType, endpoint, mechanism)
}

func NewServer(endpoint string, mechanism zmtp.SecurityMechanism) (Socket, error) {
	return NewSocket(ServerSocketType, endpoint, mechanism)
}

func (s *socket) Recv() ([]byte, error) {
	var msg []byte
	return msg, ErrNotImplemented
}

func (s *socket) Send([]byte) error {
	return ErrNotImplemented
}
