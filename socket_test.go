package zeromq

import (
	"net"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Log("attempting to create server...")

	var addr net.Addr
	var err error

	go func() {
		server := NewServer(NewSecurityNull())
		addr, err = server.Bind("tcp://127.0.0.1:9999")
		t.Logf("NETADDR: %q", addr.String())
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(100)

	client := NewClient(NewSecurityNull())
	err = client.Connect("tcp://127.0.0.1:9999")
	if err != nil {
		t.Error(err)
	}
}
