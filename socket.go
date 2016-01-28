package zeromq

import (
	"errors"
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

	defaultRetry = 250 * time.Millisecond
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
	SetRetry(retry time.Duration)
	GetRetry() time.Duration
}

type socket struct {
	sockType      zmtp.SocketType
	asServer      bool
	conns         []*Connection
	retryInterval time.Duration
	lock          sync.Mutex
	mechanism     zmtp.SecurityMechanism
}

func NewSocket(sockType zmtp.SocketType, asServer bool, mechanism zmtp.SecurityMechanism) Socket {
	return &socket{
		asServer:      asServer,
		sockType:      sockType,
		retryInterval: defaultRetry,
		mechanism:     mechanism,
		conns:         make([]*Connection, 0),
	}
}

func (s *socket) Connect(endpoint string) error {
	if s.asServer {
		return ErrInvalidSockAction
	}

	parts := strings.Split(endpoint, "://")

Connect:
	netconn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		time.Sleep(s.GetRetry())
		goto Connect
	}

	zmtpconn := zmtp.NewConnection(netconn)
	_, err = zmtpconn.Prepare(s.mechanism, s.sockType, s.asServer, nil)
	if err != nil {
		return err
	}

	conn := &Connection{
		netconn:  netconn,
		zmtpconn: zmtpconn,
	}

	s.conns = append(s.conns, conn)
	return nil
}

func (s *socket) Bind(endpoint string) (net.Addr, error) {
	var addr net.Addr

	if !s.asServer {
		return addr, ErrInvalidSockAction
	}

	parts := strings.Split(endpoint, "://")

	ln, err := net.Listen(parts[0], parts[1])
	if err != nil {
		return addr, err
	}

	netconn, err := ln.Accept()
	if err != nil {
		return addr, err
	}

	zmtpconn := zmtp.NewConnection(netconn)
	_, err = zmtpconn.Prepare(s.mechanism, s.sockType, s.asServer, nil)
	if err != nil {
		return netconn.LocalAddr(), err
	}

	conn := &Connection{
		netconn:  netconn,
		zmtpconn: zmtpconn,
	}

	s.conns = append(s.conns, conn)
	return netconn.LocalAddr(), nil
}

func (s *socket) GetRetry() time.Duration {
	return s.retryInterval
}

func (s *socket) SetRetry(r time.Duration) {
	s.retryInterval = r
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
