# gogozmq
Pure Go Implementation of a Subset of ZeroMQ.
**Danger Will Robinson, Danger**: if this note is here, this code is highly experimental. There will be false starts, APIs will change and things will break. If you need to use ZeroMQ from Go right now, don't use this. Use [GoCZMQ](https://github.com/zeromq/goczmq) instead.

## Problems
* Go's performance is less than optimal when calling C Code
* Managing C dependencies when working with Go is cumbersome
* CZMQ has a large API surface and we are only interested in a subset of it

## Proposed Solution
GoGoZMQ will be a pure go implementation of a subset of ZMTP, wrapped with a Go friendly API. After a year of working on and with [GoCZMQ](https://github.com/zeromq/goczmq), we have come to some conclusions on our use of ZeroMQ from Go:
* The [GoCZMQ Channeler API](https://godoc.org/github.com/zeromq/goczmq#Channeler) covers the majority of the ways we wish to use ZeroMQ from Go.
* The TCP transport covers most of our use case for ZeroMQ in Go. 
* The ZMQ_PUB, ZMQ_SUB, ZMQ_CLIENT, and ZMQ_SERVER socket types cover most of our use cases.
* We do not care about ZeroMQ versions before version 4.

We feel if we get a working implementation of the Channeler API from GoCZMQ, that works with PUB/SUB and CLIENT/SERVER sockets over TCP, this covers a large set of cases we currently use GoCZMQ for.

For the initial implementation, we are going to punt on worrying about CURVE encryption and ZAP support. If we get as far as having a pure Go implementation of Channeler that successfully interacts with other ZeroMQ implementations, we will then worry about the problem of authentication and encryption.

Contribution solutions for problems outside of the scope of our initial problem set are welcome. See the [contributors guide](https://github.com/zeromq/goczmq/CONTRIBUTING.md) for our process.
