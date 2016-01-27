package zeromq

import "testing"

func TestNewServer(t *testing.T) {
	server, err := NewServer(NewSecurityNull())
	if err != nil {
		t.Error(err)
	}

	_, err = server.Connect("tcp://127.0.0.1:9999")
	if err != ErrInvalidSockAction {
		t.Error(err)
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(NewSecurityNull())
	if err != nil {
		t.Error(err)
	}

	err = client.Bind("tcp://127.0.0.1:9999")
	if err != ErrInvalidSockAction {
		t.Error(err)
	}
}

func TestClientServer(t *testing.T) {
	server, err := NewServer(NewSecurityNull())
	if err != nil {
		t.Error(err)
	}

	err = server.Bind("tcp://127.0.0.1:9999")
	if err != nil {
		t.Error(err)
	}

	client, err := NewClient(NewSecurityNull())
	if err != nil {
		t.Error(err)
	}

	_, err = client.Connect("tcp://127.0.0.1:9998")
	if err != nil {
		t.Error(err)
	}
}
