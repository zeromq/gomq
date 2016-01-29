package test

/*
#cgo !windows pkg-config: libczmq libzmq libsodium
#cgo windows LDFLAGS: -lws2_32 -liphlpapi -lrpcrt4 -lsodium -lzmq -lczmq
#cgo windows CFLAGS: -Wno-pedantic-ms-format -DLIBCZMQ_EXPORTS -DZMQ_DEFINED_STDINT -DLIBCZMQ_EXPORTS

extern startExternalServer();
*/
import "C"

func StartExternalServer() {
	C.startExternalServer()
}