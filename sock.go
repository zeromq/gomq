package gogozmq

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

func (s *Sock) verifyProtoSignature(buf []byte) error {
	n, err := s.conn.Read(buf[:1])
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrProtoBad
	}

	if buf[0] != signature[0] {
		return ErrProtoBad
	}

	_, err = s.conn.Read(buf[1:10])
	if err != nil {
		return err
	}

	if bytes.Compare(buf[0:10], signature) != 0 {
		return ErrProtoBad
	}

	return err
}

func (s *Sock) verifyProtoVersion(buf []byte) error {
	n, err := s.conn.Read(buf[10:11])
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

func (s *Sock) verifyProtoSockType(buf []byte) error {
	n, err := s.conn.Read(buf[11:12])
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

func (s *Sock) verifyProtoIdentity(buf []byte) error {
	n, err := s.conn.Read(buf[12:14])
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

type Sock struct {
	conn     net.Conn
	sockType byte
	address  string
}

func NewPush(address string) (Sock, error) {
	s := Sock{
		sockType: Push,
		address:  address,
	}

	var err error
	addrParts := strings.Split(s.address, "://")
	s.conn, err = net.Dial(addrParts[0], addrParts[1])
	if err != nil {
		return s, err
	}

	zmtpGreetOutgoing := &zmtpGreeter{
		sockType: s.sockType,
	}

	_, err = zmtpGreetOutgoing.send(s.conn)
	if err != nil {
		return s, err
	}

	buf := make([]byte, 64)

	err = s.verifyProtoSignature(buf)
	if err != nil {
		return s, err
	}

	err = s.verifyProtoVersion(buf)
	if err != nil {
		return s, err
	}

	err = s.verifyProtoSockType(buf)
	if err != nil {
		return s, err
	}

	err = s.verifyProtoIdentity(buf)
	if err != nil {
		return s, err
	}

	return s, err
}

func (s *Sock) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("I can't read")
}

func (s *Sock) Write(b []byte) (n int, err error) {
	zmtpMessageOutgoing := &zmtpMessage{}
	zmtpMessageOutgoing.msg = [][]byte{b}
	return zmtpMessageOutgoing.send(s.conn)
}

func (s *Sock) Close() error {
	return s.conn.Close()
}

func (s *Sock) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *Sock) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Sock) SetDeadLine(t time.Time) error {
	return fmt.Errorf("I don't do deadlines")
}

func (s *Sock) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("I don't do deadlines")
}

func (s *Sock) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("I don't do deadlines")
}
