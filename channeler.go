package gogozmq

import (
	"bytes"
	"net"
	"strings"
)

var (
	signature = []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 1, 0x7f}
)

type Channeler struct {
	conn         net.Conn
	sockType     byte
	endpoints    string
	subscribe    string
	SendChan     chan<- [][]byte
	RecvChan     <-chan [][]byte
	sendDoneChan chan bool
}

func newChanneler(sockType byte, endpoints, subscribe string) (*Channeler, error) {
	sendChan := make(chan [][]byte)
	recvChan := make(chan [][]byte)
	sendDoneChan := make(chan bool)

	parts := strings.Split(endpoints, "://")
	if len(parts) != 2 {
		panic("endpoint should have 2 parts")
	}

	c := &Channeler{
		sockType:     sockType,
		endpoints:    endpoints,
		subscribe:    subscribe,
		SendChan:     sendChan,
		RecvChan:     recvChan,
		sendDoneChan: sendDoneChan,
	}

	var err error
	c.conn, err = net.Dial(parts[0], parts[1])
	if err != nil {
		return c, err
	}

	zmtpGreetOutgoing := &zmtpGreeter{
		sockType: sockType,
	}

	_, err = zmtpGreetOutgoing.send(c.conn)
	if err != nil {
		return c, err
	}

	buf := make([]byte, 64)

	err = c.verifyProtoSignature(buf)
	if err != nil {
		return c, err
	}

	err = c.verifyProtoVersion(buf)
	if err != nil {
		return c, err
	}

	err = c.verifyProtoSockType(buf)
	if err != nil {
		return c, err
	}

	err = c.verifyProtoIdentity(buf)
	if err != nil {
		return c, err
	}

	go c.sendMessages(sendChan)

	return c, err
}

func NewPushChanneler(endpoints string) (*Channeler, error) {
	c, err := newChanneler(Push, endpoints, "")
	return c, err
}

func (c *Channeler) Destroy() {
	close(c.SendChan)
	<-c.sendDoneChan
	c.conn.Close()
}

func (c *Channeler) verifyProtoSignature(buf []byte) error {
	n, err := c.conn.Read(buf[:1])
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrProtoBad
	}

	if buf[0] != signature[0] {
		return ErrProtoBad
	}

	_, err = c.conn.Read(buf[1:10])
	if err != nil {
		return err
	}

	if bytes.Compare(buf[0:10], signature) != 0 {
		return ErrProtoBad
	}

	return err
}

func (c *Channeler) verifyProtoVersion(buf []byte) error {
	n, err := c.conn.Read(buf[10:11])
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrProtoBad
	}

	if buf[10] < zmtpVersion {
		return ErrProtoBad
	}

	return nil
}

func (c *Channeler) verifyProtoSockType(buf []byte) error {
	n, err := c.conn.Read(buf[11:12])
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrProtoBad
	}

	if buf[11] != Pull {
		return ErrProtoSockType
	}

	return err
}

func (c *Channeler) verifyProtoIdentity(buf []byte) error {
	n, err := c.conn.Read(buf[12:14])
	if err != nil {
		return err
	}

	if n != 2 {
		return ErrProtoBad
	}

	if buf[12] != 0 {
		return ErrProtoIdentity
	}

	if buf[13] != 0 {
		return ErrProtoIdentity
	}

	return nil
}

func (c *Channeler) sendMessages(sendChan <-chan [][]byte) {
	zmtpMessageOutgoing := &zmtpMessage{}
	more := true

	for more {
		zmtpMessageOutgoing.msg, more = <-sendChan

		if more {
			_, err := zmtpMessageOutgoing.send(c.conn)

			if err != nil {
				// reconnection is not handled at the moment
				break
			}
		}
	}

	c.sendDoneChan <- true
}
