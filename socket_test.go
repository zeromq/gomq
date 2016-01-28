package zeromq

import (
	"net"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Log("attempting to create server...")

	var addr net.Addr
	var err error

	go func() {
		client := NewClient(NewSecurityNull())
		err = client.Connect("tcp://127.0.0.1:9999")
		if err != nil {
			t.Error(err)
		}
	}()

	server := NewServer(NewSecurityNull())
	addr, err = server.Bind("tcp://127.0.0.1:9999")
	t.Logf("NETADDR: %q", addr.String())
	if err != nil {
		t.Error(err)
	}
}
