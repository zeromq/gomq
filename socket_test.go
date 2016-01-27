package zeromq

import "testing"

func TestClientServer(t *testing.T) {
	server, err := NewServer("tcp://127.0.0.1:9999")
	if want, have := ErrNotImplemented, err; want != have {
		t.Errorf("want %q, got %q")
	}

	client, err := NewClient("tcp://127.0.0.1:9999")
	if want, have := ErrNotImplemented, err; want != have {
		t.Errorf("want %q, got %q")
	}

	err = client.Send([]byte("hello"))
	if want, have := ErrNotImplemented, err; want != have {
		t.Errorf("want %q, got %q")
	}

	_, err = server.Recv()
	if want, have := ErrNotImplemented, err; want != have {
		t.Errorf("want %q, got %q")
	}

}
