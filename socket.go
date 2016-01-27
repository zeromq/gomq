package zeromq

import (
	"errors"
	"net"
	"strings"

	"github.com/zeromq/gozmtp"
)

var (
	ClientSocketType = zmtp.ClientSocketType
	ServerSocketType = zmtp.ServerSocketType

	NullSecurityMechanismType  = zmtp.NullSecurityMechanismType
	PlainSecurityMechanismType = zmtp.PlainSecurityMechanismType
	CurveSecurityMechanismTyp  = zmtp.CurveSecurityMechanismType

	ErrNotImplemented    = errors.New("not implemented")
	ErrInvalidSockAction = errors.New("action not valid on this socket")
)

type Socket interface {
	Recv() ([]byte, error)
	Send([]byte) error
	Connect(endpoint string) (net.Conn, error)
	Bind(endpoint string) error
}

type socket struct {
	sockType zmtp.SocketType
	isServer bool
	conns    []net.Conn
}

func (s *socket) Connect(endpoint string) (net.Conn, error) {
	if s.isServer {
		return nil, ErrInvalidSockAction
	}

	parts := strings.Split(endpoint, "://")
	conn, err := net.Dial(parts[0], parts[1])
	return conn, err
}

func (s *socket) Bind(endpoint string) error {
	if !s.isServer {
		return ErrInvalidSockAction
	}

	return ErrNotImplemented
}

func NewSecurityNull() *zmtp.SecurityNull {
	return zmtp.NewSecurityNull()
}

func NewSocket(sockType zmtp.SocketType, isServer bool, mechanism zmtp.SecurityMechanism) (Socket, error) {
	return &socket{
		isServer: isServer,
		sockType: sockType,
	}, nil
}

func NewClient(mechanism zmtp.SecurityMechanism) (Socket, error) {
	return NewSocket(ClientSocketType, false, mechanism)
}

func NewServer(mechanism zmtp.SecurityMechanism) (Socket, error) {
	return NewSocket(ServerSocketType, true, mechanism)
}

func (s *socket) Recv() ([]byte, error) {
	var msg []byte
	return msg, ErrNotImplemented
}

func (s *socket) Send([]byte) error {
	return ErrNotImplemented
}
