package gogozmq

import (
	"testing"

	"github.com/zeromq/goczmq"
)

func TestPushChanneler(t *testing.T) {
	endpoint := "tcp://127.0.0.1:9999"

	pull, err := goczmq.NewPull(endpoint)
	if err != nil {
		t.Error(err)
	}
	defer pull.Destroy()

	push, err := NewPushChanneler(endpoint)
	defer push.Destroy()
	if err != nil {
		t.Fatal(err)
	}

	push.SendChan <- [][]byte{[]byte("Hello")}

	msg, more, err := pull.RecvFrame()

	if err != nil {
		t.Fatal(err)
	}

	text := string(msg)

	if text != "Hello" {
		t.Fatal("Wrong message")
	}

	if more != 0 {
		t.Fatal("More flag is wrong")
	}
}
