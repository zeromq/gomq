package gogozmq

import (
	"testing"

	"github.com/zeromq/goczmq"
)

func TestPushSockShortMessage(t *testing.T) {
	endpoint := "tcp://127.0.0.1:9999"

	pull, err := goczmq.NewPull(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	defer pull.Destroy()

	push, err := NewPush(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	push.Write([]byte("Hello"))

	msg, more, err := pull.RecvFrame()
	if err != nil {
		t.Fatal(err)
	}

	if want, have := "Hello", string(msg); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	if want, have := 5, len(msg); want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}

	if want, have := 0, more; want != have {
		t.Errorf("want %#v, have %#v", want, have)
	}
}
