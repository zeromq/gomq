package gogozmq

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

var (
	zmtpGreetOutgoing = &zmtpGreet{
		signature: [10]byte{0xff, 0, 0, 0, 0, 0, 0, 0, 1, 0x7f},
		version:   [2]byte{3, 0},
		mechanism: [20]byte{'N', 'U', 'L', 'L', 0},
		asServer:  [1]byte{0},
	}
)

type Channeler struct {
	conn      net.Conn
	sockType  int
	endpoints string
	subscribe string
	SendChan  chan<- [][]byte
	RecvChan  <-chan [][]byte
}

func newChanneler(sockType int, endpoints, subscribe string) (*Channeler, error) {
	sendChan := make(chan [][]byte)
	recvChan := make(chan [][]byte)

	parts := strings.Split(endpoints, "://")
	if len(parts) != 2 {
		panic("endpoint should have 2 parts")
	}

	c := &Channeler{
		sockType:  sockType,
		endpoints: endpoints,
		subscribe: subscribe,
		SendChan:  sendChan,
		RecvChan:  recvChan,
	}

	var err error
	c.conn, err = net.Dial(parts[0], parts[1])
	if err != nil {
		return c, err
	}

	_, err = zmtpGreetOutgoing.sendSignature(c.conn)
	if err != nil {
		return c, err
	}

	buf := make([]byte, 255)
	_, err = c.conn.Read(buf[:1])

	if buf[0] != 0xff {
		return c, fmt.Errorf("bad protocol signature")
	}

	_, err = c.conn.Read(buf[1:10])

	if bytes.Compare(buf[0:10], zmtpGreetOutgoing.signature[:]) != 0 {
		return c, fmt.Errorf("bad protocol signature")
	}

	_, err = zmtpGreetOutgoing.sendVersion(c.conn)
	if err != nil {
		return c, err
	}

	_, err = c.conn.Read(buf[10:11])
	if err != nil {
		return c, err
	}

	if buf[10] != 0x03 {
		return c, fmt.Errorf("bad protocol version")
	}

	return c, err
}

func NewPushChanneler(endpoints string) (*Channeler, error) {
	c, err := newChanneler(Push, endpoints, "")
	return c, err
}

func (c *Channeler) Destroy() {
	c.conn.Close()
}
