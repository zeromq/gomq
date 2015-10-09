package gogozmq

import (
	"testing"
	"time"

	"github.com/zeromq/goczmq"
)

func TestPushChanneler(t *testing.T) {
	push := NewPushChanneler("inproc://TestPubChanneler")
	defer push.Destroy()

	pull, err := goczmq.NewPull("inproc://TestPubChanneler")
	if err != nil {
		t.Error(err)
	}
	defer pull.Destroy()

	select {
	case push.SendChan <- [][]byte{[]byte("hello")}:
	case <-time.After(1):
		t.Fatalf("send timed out")
	}

	resp, err := pull.RecvMessageNoWait()
	if err != nil {
		t.Error(err)
	}

	if want, got := 1, len(resp); want != got {
		t.Errorf("want '%#v', got '%#v'", want, got)
	}

	if want, got := "hello", string(resp[0]); want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}
