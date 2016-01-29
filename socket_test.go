package zeromq

import (
	"bytes"
	"net"
	"testing"
)

func TestNewClient(t *testing.T) {
	var addr net.Addr
	var err error

	go func() {
		client := NewClient(NewSecurityNull())
		err = client.Connect("tcp://127.0.0.1:9999")
		if err != nil {
			t.Error(err)
		}

		err := client.Send([]byte("HELLO"))
		if err != nil {
			t.Error(err)
		}

		msg, _ := client.Recv()
		if want, got := 0, bytes.Compare([]byte("WORLD"), msg); want != got {
			t.Errorf("want %v, got %v", want, got)
		}

		err = client.Send([]byte("GOODBYE"))
		if err != nil {
			t.Error(err)
		}
	}()

	server := NewServer(NewSecurityNull())

	addr, err = server.Bind("tcp://127.0.0.1:9999")

	if want, got := "127.0.0.1:9999", addr.String(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}

	if err != nil {
		t.Error(err)
	}

	msg, _ := server.Recv()

	if want, got := 0, bytes.Compare([]byte("HELLO"), msg); want != got {
		t.Errorf("want %v, got %v", want, got)
	}

	server.Send([]byte("WORLD"))

	msg, _ = server.Recv()

	if want, got := 0, bytes.Compare([]byte("GOODBYE"), msg); want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
