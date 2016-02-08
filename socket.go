package gomq

import (
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/zeromq/gomq/zmtp"
)

var (
	// ErrInvalidSockAction is returned when an action is performed
	// on a socket type that does not support the action
	ErrInvalidSockAction = errors.New("action not valid on this socket")

	defaultRetry = 250 * time.Millisecond
)

// Connection holds a connection to a ZeroMQ socket.
type Connection struct {
	netconn  net.Conn
	zmtpconn *zmtp.Connection
}

// Socket represents a ZeroMQ socket. Sockets may have multiple connections.
type Socket interface {
	Recv() ([]byte, error)
	Send([]byte) error
	Connect(endpoint string) error
	Bind(endpoint string) (net.Addr, error)
	Close()
}

type socket struct {
	sockType      zmtp.SocketType
	asServer      bool
	conns         map[string]*Connection
	clients       []string
	retryInterval time.Duration
	lock          *sync.RWMutex
	mechanism     zmtp.SecurityMechanism
	messageChan   chan *zmtp.Message
}

func newSocket(sockType zmtp.SocketType, asServer bool, mechanism zmtp.SecurityMechanism) Socket {
	return &socket{
		lock:          &sync.RWMutex{},
		asServer:      asServer,
		sockType:      sockType,
		retryInterval: defaultRetry,
		mechanism:     mechanism,
		conns:         make(map[string]*Connection),
		clients:       make([]string, 0),
		messageChan:   make(chan *zmtp.Message),
	}
}

// Connect connects to an endpoint.
// TODO: this call should be non blocking
func (s *socket) Connect(endpoint string) error {
	if s.asServer {
		return ErrInvalidSockAction
	}

	parts := strings.Split(endpoint, "://")

Connect:
	netconn, err := net.Dial(parts[0], parts[1])
	if err != nil {
		time.Sleep(s.retryInterval)
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

	s.addConn(conn)
	zmtpconn.Recv(s.messageChan)
	return nil
}

// Bind binds to an endpoint.
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

	s.addConn(conn)
	zmtpconn.Recv(s.messageChan)
	return netconn.LocalAddr(), nil
}

func (s *socket) addConn(conn *Connection) {
	s.lock.Lock()
	uuid, _ := newUUID()
	s.conns[uuid] = conn
	s.clients = append(s.clients, uuid)
	s.lock.Unlock()
}

func (s *socket) removeConn(uuid string) {
	s.lock.Lock()
	for k, v := range s.clients {
		if v == uuid {
			s.clients = append(s.clients[:k], s.clients[k+1:]...)
		}
	}
	delete(s.conns, uuid)
	s.lock.Unlock()
}

// Close closes all underlying connections in a socket.
func (s *socket) Close() {
	s.lock.Lock()
	for _, v := range s.clients {
		s.conns[v].netconn.Close()
		s.removeConn(v)
	}
	s.lock.Unlock()
}

// NewClient creates a new ZMQ_CLIENT socket.
func NewClient(mechanism zmtp.SecurityMechanism) Socket {
	return newSocket(zmtp.ClientSocketType, false, mechanism)
}

// NewServer creates a new ZMQ_SERVER socket.
func NewServer(mechanism zmtp.SecurityMechanism) Socket {
	return newSocket(zmtp.ServerSocketType, true, mechanism)
}

// Recv receives the next message from the socket.
func (s *socket) Recv() ([]byte, error) {
	msg := <-s.messageChan
	if msg.MessageType == zmtp.CommandMessage {
	}
	return msg.Body, msg.Err
}

// Send sends a message.
func (s *socket) Send(b []byte) error {
	return s.conns[s.clients[0]].zmtpconn.SendFrame(b)
}
