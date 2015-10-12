package gogozmq

import (
	"io"
	"time"
)

const (
	Push                     byte          = 8
	Pull                     byte          = 7
	zmtpVersion              byte          = 1
	zmtpHandshakeTimeout     time.Duration = time.Millisecond * 100
	shortMessageEnvelopeSize               = 2
	flagIndex                              = 0
	shortSizeIndex                         = 1
	finalShortSizeFlag                     = 0
)

type zmtpGreet struct {
	sockType byte
}

type zmtpMessage struct {
	msg [][]byte
}

func (z *zmtpGreet) send(w io.Writer) (int, error) {
	finalShort := byte(0x00)
	identitySize := byte(0x00) // identity are not supported, so zero size identity

	greeting := []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 1, 0x7f, zmtpVersion, z.sockType, finalShort, identitySize}

	return w.Write(greeting)
}

func (z *zmtpMessage) send(w io.Writer) (int, error) {
	payloadSize := len(z.msg[0])
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
