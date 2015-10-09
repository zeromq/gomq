package gogozmq

const (
	Push = 7
	Pull = 8
)

type Channeler struct {
	sockType  int
	endpoints string
	subscribe string
	SendChan  chan<- [][]byte
	RecvChan  <-chan [][]byte
}

func newChanneler(sockType int, endpoints, subscribe string) *Channeler {
	sendChan := make(chan [][]byte)
	recvChan := make(chan [][]byte)
	return &Channeler{
		sockType:  sockType,
		endpoints: endpoints,
		subscribe: subscribe,
		SendChan:  sendChan,
		RecvChan:  recvChan,
	}
}

func NewPushChanneler(endpoints string) *Channeler {
	return newChanneler(Push, endpoints, "")
}

func NewPullChanneler(endpoints string) *Channeler {
	return newChanneler(Pull, endpoints, "")
}

func (c *Channeler) Destroy() {
}
