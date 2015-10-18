package gogozmq

import (
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	// Push is a ZMQ_PUSH socket
	Push byte = 8

	// Pull is a ZMQ_PULL socket
	Pull byte = 7

	zmtpVersion              byte          = 1
	zmtpHandshakeTimeout     time.Duration = time.Millisecond * 100
	shortMessageEnvelopeSize               = 2
	flagIndex                              = 0
	shortSizeIndex                         = 1
	finalShortSizeFlag                     = 0
)

var (
	// ErrProtoBad is returned when the protocol version header is incorrect
	ErrProtoBad = errors.New("bad protocol signature")

	// ErrProtoVersion is returned for unsupported ZMTP versions
	ErrProtoVersion = errors.New("bad protocol version")

	// ErrProtoSockType is returned for socket types we do not support
	ErrProtoSockType = errors.New("unsupported socket type")

	// ErrProtoIdentity is returned when the identity portion of
	// a handshake is incorrect
	ErrProtoIdentity = errors.New("bad protocol identity")
)

type zmtpGreeter struct {
	sockType byte
}

type zmtpMessage struct {
	msg [][]byte
}

func (z *zmtpGreeter) send(w io.Writer) (int, error) {
	finalShort := byte(0x00)
	identitySize := byte(0x00) // identity are not supported, so zero size identity

	greeting := []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 1, 0x7f, zmtpVersion, z.sockType, finalShort, identitySize}

	return w.Write(greeting)
}

func (z *zmtpMessage) send(w io.Writer) (int, error) {
	payloadSize := len(z.msg[0])
	if payloadSize > 255 {
		return 0, fmt.Errorf("long messages not supported")
	}

	envelopeSize := shortMessageEnvelopeSize + payloadSize

	// only support single part and short message (under 255 bytes)
	envelope := make([]byte, envelopeSize)

	envelope[flagIndex] = finalShortSizeFlag
	envelope[shortSizeIndex] = byte(payloadSize)

	if payloadSize > 0 {
		copy(envelope[2:], z.msg[0])
	}
	return w.Write(envelope)
}
