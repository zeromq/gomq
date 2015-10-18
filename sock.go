package gogozmq

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Conn interface {
	net.Conn
	Connect(endpoint string) error
	CurrentIdx() int
}

type Sock struct {
	conns      []net.Conn
	currentIdx int
	sockType   byte
	address    string
}

func NewPushConn(address string) (Conn, error) {
	s := Sock{
		currentIdx: 0,
		sockType:   Push,
		address:    address,
	}

	addrParts := strings.Split(s.address, "://")
	conn, err := net.Dial(addrParts[0], addrParts[1])
	if err != nil {
		return s, err
	}
	s.conns = append(s.conns, conn)

	zmtpGreetOutgoing := &zmtpGreeter{
		sockType: s.sockType,
	}

	_, err = zmtpGreetOutgoing.send(s.conns[s.currentIdx])
	if err != nil {
		return s, err
	}

	buf := make([]byte, 64)

	err = clientHandshake(conn, buf)
	return s, err
}

func (s Sock) Connect(endpoint string) error {
	addrParts := strings.Split(s.address, "://")
	_, err := net.Dial(addrParts[0], addrParts[1])
	if err != nil {
		return err
	}
	return err
}

func (s Sock) CurrentIdx() int {
	return s.currentIdx
}

func (s Sock) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("I can't read")
}

func (s Sock) Write(b []byte) (n int, err error) {
	zmtpMessageOutgoing := &zmtpMessage{}
	zmtpMessageOutgoing.msg = [][]byte{b}
	n, err = zmtpMessageOutgoing.send(s.conns[s.currentIdx])
	s.currentIdx++
	return n, err
}

func (s Sock) Close() error {
	var err error
	for i := 0; i < len(s.conns); i++ {
		err = s.conns[i].Close()
		if err != nil {
			return err
		}
	}
	return err
}

func (s Sock) LocalAddr() net.Addr {
	return s.conns[s.currentIdx].LocalAddr()
}

func (s Sock) RemoteAddr() net.Addr {
	return s.conns[s.currentIdx].RemoteAddr()
}

func (s Sock) SetDeadline(t time.Time) error {
	return fmt.Errorf("I don't do deadlines")
}

func (s Sock) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("I don't do deadlines")
}

func (s Sock) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("I don't do deadlines")
}
