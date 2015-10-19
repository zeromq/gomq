package gogozmq

var (
	signature = []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 1, 0x7f}
)

type Channeler struct {
	conn         Conn
	sockType     byte
	endpoints    string
	subscribe    string
	SendChan     chan<- []byte
	RecvChan     <-chan []byte
	sendDoneChan chan bool
}

func newChanneler(sockType byte, endpoints, subscribe string) (*Channeler, error) {
	sendChan := make(chan []byte)
	recvChan := make(chan []byte)
	sendDoneChan := make(chan bool)

	c := &Channeler{
		sockType:     sockType,
		endpoints:    endpoints,
		subscribe:    subscribe,
		SendChan:     sendChan,
		RecvChan:     recvChan,
		sendDoneChan: sendDoneChan,
	}

	var err error
	c.conn, err = NewPushConn(endpoints)
	if err != nil {
		return c, err
	}

	go c.sendMessages(sendChan)
	return c, err
}

func NewPushChanneler(endpoints string) (*Channeler, error) {
	c, err := newChanneler(Push, endpoints, "")
	return c, err
}

func (c *Channeler) Destroy() {
	close(c.SendChan)
	<-c.sendDoneChan
	c.conn.Close()
}

func (c *Channeler) sendMessages(sendChan <-chan []byte) {
	more := true
	var msg []byte
	for more {
		msg, more = <-sendChan

		if more {
			_, err := c.conn.Write(msg)

			if err != nil {
				// reconnection is not handled at the moment
				break
			}
		}
	}

	c.sendDoneChan <- true
}
