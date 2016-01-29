# gomq
Pure Go Implementation of a Subset of ZeroMQ.
**Danger Will Robinson, Danger**: if this note is here, this code is highly experimental. There will be false starts, APIs will change and things will break. If you need to use ZeroMQ from Go right now, don't use this. Use [GoCZMQ](https://github.com/zeromq/goczmq) instead.

## Problems
* Go's performance is less than optimal when calling C Code
* Managing C dependencies when working with Go is cumbersome
* CZMQ has a large API surface and we are only interested in a subset of it

## Proposed Solution
gomq will be a pure go implementation of a subset of ZMTP, wrapped with a Go friendly API. GoMQ will only implement ZMTP version 3.x and will not be backwards compatible with previous versions of ZMTP. The initial implementation will support ZMQ_CLIENT and ZMQ_SERVER with the NULL security mechanism.
