package zeromq

import "testing"

func TestClientServer(t *testing.T) {

	server, err := NewServer("tcp://127.0.0.1:9999", NewSecurityNull())
	if err != nil {
		t.Error(err)
	}

	client, err := NewClient("tcp://127.0.0.1:9999", NewSecurityNull())
	if err != nil {
		t.Error(err)
	}

	err = client.Send([]byte("hello"))
	if err != nil {
		t.Error(err)
	}

	msg, err := server.Recv()
	if err != nil {
		t.Error(err)
	}

	if want, have := "hello", string(msg); want != have {
		t.Errorf("want %q, have %q", want, have)
	}

}
