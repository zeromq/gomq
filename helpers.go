package gogozmq

import (
	"bytes"
	"net"
)

func verifyProtoSignature(conn net.Conn, buf []byte) error {
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

func verifyProtoVersion(conn net.Conn, buf []byte) error {
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

func verifyProtoSockType(conn net.Conn, buf []byte) error {
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

func verifyProtoIdentity(conn net.Conn, buf []byte) error {
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

func clientHandshake(conn net.Conn, buf []byte) error {
	err := verifyProtoSignature(conn, buf)
	if err != nil {
		return err
	}

	err = verifyProtoVersion(conn, buf)
	if err != nil {
		return err
	}

	err = verifyProtoSockType(conn, buf)
	if err != nil {
		return err
	}

	err = verifyProtoIdentity(conn, buf)
	return err
}
