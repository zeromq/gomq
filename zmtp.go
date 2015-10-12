package gogozmq

import (
	"io"
	"time"
)

const (
	Push                 byte          = 8
	Pull                 byte          = 7
	zmtpVersion          byte          = 1
	zmtpHandshakeTimeout time.Duration = time.Millisecond * 100
)

type zmtpGreet struct {
	sockType byte
}

func (z *zmtpGreet) send(w io.Writer) (int, error) {
	finalShort := byte(0x00)
	identitySize := byte(0x00) // identity are not supported, so zero size identity

	greeting := []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 1, 0x7f, zmtpVersion, z.sockType, finalShort, identitySize}

	return w.Write(greeting)
}
