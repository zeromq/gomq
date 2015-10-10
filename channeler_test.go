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
}
