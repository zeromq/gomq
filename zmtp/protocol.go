package zmtp

import "encoding/binary"

const (
	majorVersion uint8 = 3
	minorVersion uint8 = 0
)

const (
	signaturePrefix = 0xFF
	signatureSuffix = 0x7F
)

const (
	hasMoreBitFlag   = 0x1
	isLongBitFlag    = 0x2
	isCommandBitFlag = 0x4
)

var (
	version = [2]uint8{majorVersion, minorVersion}
)

var byteOrder = binary.BigEndian

const maxUint = ^uint(0)
const minUint = 0
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1
const maxUint64 = ^uint64(0)
const minUint64 = 0
const maxInt64 = int64(maxUint64 >> 1)
const minInt64 = -maxInt64 - 1

type greeting struct {
	SignaturePrefix byte
	_               [8]byte
	SignatureSuffix byte
	Version         [2]uint8
	Mechanism       [20]byte
	ServerFlag      byte
	_               [31]byte
}

type Command struct {
	Index int
	Name  string
	Body  []byte
}

type Message struct {
	Index int
	Body  []byte
}

type Error struct {
	Index int
	Error error
}
