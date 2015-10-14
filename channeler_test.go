package gogozmq

import (
	"testing"

	"github.com/zeromq/goczmq"
)

func TestPushChanneler(t *testing.T) {
	endpoint := "tcp://127.0.0.1:9999"

	pull, err := goczmq.NewPull(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	defer pull.Destroy()

	push, err := NewPushChanneler(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	push.SendChan <- [][]byte{[]byte("Hello")}

	msg, more, err := pull.RecvFrame()
	if err != nil {
		t.Fatal(err)
	}

	if want, have := "Hello", string(msg); want != have {
		t.Error("want %#v, have %#v", want, have)
	}

	if want, have := 0, more; want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}
