package gogozmq

import (
	"errors"
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
