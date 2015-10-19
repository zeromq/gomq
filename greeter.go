package gogozmq

import (
	"bytes"
	"io"
	"net"
)

type greeter struct {
	sockType byte
}

func (g *greeter) greet(conn net.Conn) error {
	_, err := g.send(conn)
	if err != nil {
		return err
	}

	buf := make([]byte, 64)
	err = g.verifyProtoSignature(conn, buf)
	if err != nil {
		return err
	}

	err = g.verifyProtoVersion(conn, buf)
	if err != nil {
		return err
	}

	err = g.verifyProtoSockType(conn, buf)
	if err != nil {
		return err
	}

	err = g.verifyProtoIdentity(conn, buf)
	return err
}

func (g *greeter) send(w io.Writer) (int, error) {
	finalShort := byte(0x00)
	identitySize := byte(0x00) // identity are not supported, so zero size identity

	greeting := []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 1, 0x7f, zmtpVersion, g.sockType, finalShort, identitySize}

	return w.Write(greeting)
}

func (g *greeter) verifyProtoSignature(conn net.Conn, buf []byte) error {
	n, err := conn.Read(buf[:1])
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrProtoBad
	}

	if buf[0] != signature[0] {
		return ErrProtoBad
	}

	_, err = conn.Read(buf[1:10])
	if err != nil {
		return err
	}

	if bytes.Compare(buf[0:10], signature) != 0 {
		return ErrProtoBad
	}

	return err
}

func (g *greeter) verifyProtoVersion(conn net.Conn, buf []byte) error {
	n, err := conn.Read(buf[10:11])
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

func (g *greeter) verifyProtoSockType(conn net.Conn, buf []byte) error {
	n, err := conn.Read(buf[11:12])
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

func (g *greeter) verifyProtoIdentity(conn net.Conn, buf []byte) error {
	n, err := conn.Read(buf[12:14])
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
