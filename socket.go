package zeromq

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"time"

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

type Connection struct {
	netconn  net.Conn
	zmtpconn *zmtp.Connection
}

type Socket interface {
	Recv() ([]byte, error)
	Send([]byte) error
	Connect(endpoint string) error
	Bind(endpoint string) (net.Addr, error)
}

type socket struct {
	sockType      zmtp.SocketType
	isServer      bool
	conns         []*Connection
	retryInterval time.Duration
	lock          sync.Mutex
}

func NewSocket(sockType zmtp.SocketType, isServer bool, mechanism zmtp.SecurityMechanism) Socket {
	return &socket{
		isServer: isServer,
		sockType: sockType,
		conns:    make([]*Connection, 0),
	}
}

func (s *socket) Connect(endpoint string) error {
	if s.isServer {
		return ErrInvalidSockAction
	}

	parts := strings.Split(endpoint, "://")

	log.Printf("connecting with %q on %q", parts[0], parts[1])

	netconn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		return err
	}

	conn := &Connection{
		netconn:  netconn,
		zmtpconn: zmtp.NewConnection(netconn),
	}

	s.conns = append(s.conns, conn)
	return nil
}

func (s *socket) Bind(endpoint string) (net.Addr, error) {
	var addr net.Addr

	if !s.isServer {
		return addr, ErrInvalidSockAction
	}

	parts := strings.Split(endpoint, "://")

	log.Printf("listening for %q on %q", parts[0], parts[1])

	ln, err := net.Listen(parts[0], parts[1])
	if err != nil {
		return addr, err
	}

	netconn, err := ln.Accept()
	if err != nil {
		return addr, err
	}

	conn := &Connection{
		netconn:  netconn,
		zmtpconn: zmtp.NewConnection(netconn),
	}

	s.conns = append(s.conns, conn)
	return netconn.LocalAddr(), nil
}

func NewSecurityNull() *zmtp.SecurityNull {
	return zmtp.NewSecurityNull()
}

func NewClient(mechanism zmtp.SecurityMechanism) Socket {
	return NewSocket(ClientSocketType, false, mechanism)
}

func NewServer(mechanism zmtp.SecurityMechanism) Socket {
	return NewSocket(ServerSocketType, true, mechanism)
}

func (s *socket) Recv() ([]byte, error) {
	var msg []byte
	return msg, ErrNotImplemented
}

func (s *socket) Send([]byte) error {
	return ErrNotImplemented
}
