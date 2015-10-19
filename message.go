package gogozmq

import (
	"fmt"
	"io"
)

type message struct {
	msg [][]byte
}

func (m *message) send(w io.Writer) (int, error) {
	payloadSize := len(m.msg[0])
	if payloadSize > 255 {
		return 0, fmt.Errorf("long messages not supported")
	}

	envelopeSize := shortMessageEnvelopeSize + payloadSize

	// only support single part and short message (under 255 bytes)
	envelope := make([]byte, envelopeSize)

	envelope[flagIndex] = finalShortSizeFlag
	envelope[shortSizeIndex] = byte(payloadSize)

	if payloadSize > 0 {
		copy(envelope[2:], m.msg[0])
	}
	return w.Write(envelope)
}
