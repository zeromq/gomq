package zmtp

import "errors"

type Socket interface {
	Type() SocketType
	IsSocketTypeCompatible(socketType SocketType) bool
	IsCommandTypeValid(name string) bool
}

func NewSocket(socketType SocketType) (Socket, error) {
	switch socketType {
	case ClientSocketType:
		return clientSocket{}, nil
	case ServerSocketType:
		return serverSocket{}, nil
	default:
		return nil, errors.New("Invalid socket type")
	}
}

type clientSocket struct{}
type serverSocket struct{}

func (clientSocket) Type() SocketType {
	return ClientSocketType
}

func (clientSocket) IsSocketTypeCompatible(socketType SocketType) bool {
	return socketType == ServerSocketType
}

func (clientSocket) IsCommandTypeValid(name string) bool {
	return false
}

func (serverSocket) IsSocketTypeCompatible(socketType SocketType) bool {
	return socketType == ClientSocketType
}

func (serverSocket) IsCommandTypeValid(name string) bool {
	return false
}

func (serverSocket) Type() SocketType {
	return ServerSocketType
}
