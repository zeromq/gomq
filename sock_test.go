package gogozmq

import (
	"testing"

	"github.com/zeromq/goczmq"
)

func TestPushSockShortMessage(t *testing.T) {
	endpoint1 := "tcp://127.0.0.1:9998"

	pull1, err := goczmq.NewPull(endpoint1)
	if err != nil {
		t.Fatal(err)
	}
	defer pull1.Destroy()

	push, err := NewPushConn(endpoint1)
	if err != nil {
		t.Fatal(err)
	}

	push.Write([]byte("Hello"))

	msg, more, err := pull1.RecvFrame()
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
