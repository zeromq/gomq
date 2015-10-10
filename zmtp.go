package gogozmq

import (
	"io"
	"time"
)

const (
	Push                               = 7
	Pull                               = 8
	zmtpVersionMajor     byte          = 3
	zmtpVersionMinor     byte          = 0
	zmtpHandshakeTimeout time.Duration = time.Millisecond * 100
)

type zmtpGreet struct {
	signature [10]byte
	version   [2]byte
	mechanism [20]byte
	asServer  [1]byte
	filler    [31]byte
}

func (z *zmtpGreet) sendSignature(w io.Writer) (int, error) {
	return w.Write(z.signature[:])
}

func (z *zmtpGreet) sendVersion(w io.Writer) (int, error) {
	return w.Write(z.version[:])
}

func (z *zmtpGreet) sendMechanism(w io.Writer) (int, error) {
	return w.Write(z.mechanism[:])
}

func (z *zmtpGreet) sendAsServer(w io.Writer) (int, error) {
	return w.Write(z.asServer[:])
}

func (z *zmtpGreet) sendFiller(w io.Writer) (int, error) {
	return w.Write(z.filler[:])
}
