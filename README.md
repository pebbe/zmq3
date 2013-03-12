A Go interface to [ZeroMQ](http://www.zeromq.org/) version 3.

With limited support for ZeroMQ version 2.

## Install

    go get github.com/pebbe/zmq3

## Docs

 * [package help](http://godoc.org/github.com/pebbe/zmq3)

## To do

 * Re-implementing the remaining examples for [ØMQ - The Guide](http://zguide.zeromq.org/page:all).
   Currently, all examples from chapters 1 to 3 are finished.
 * Signal handling? (includes example `interrupt' in The Guide

## Support for ZeroMQ version 2

This package installs also with ZeroMQ version 2, but will miss some of
its functionality for version 3.

 * Things specific for version 2 are not implemented.
 * The following functions will return an error with ZeroMQ version 2:
  * GetMaxSockets
  * SetMaxSockets
  * (*Socket) Disconnect
  * (*Socket) Monitor
  * (*Socket) RecvEvent
  * (*Socket) GetDelayAttachOnConnect
  * (*Socket) GetIpv4only
  * (*Socket) GetLastEndpoint
  * (*Socket) GetMaxmsgsize
  * (*Socket) GetMulticastHops)
  * (*Socket) GetTcpKeepalive
  * (*Socket) GetTcpKeepaliveCnt
  * (*Socket) GetTcpKeepaliveIdle
  * (*Socket) GetTcpKeepaliveIntvl
  * (*Socket) SetDelayAttachOnConnect
  * (*Socket) SetIpv4only
  * (*Socket) SetMaxmsgsize
  * (*Socket) SetMulticastHops
  * (*Socket) SetRouterMandatory
  * (*Socket) SetTcpAcceptFilter
  * (*Socket) SetTcpKeepalive
  * (*Socket) SetTcpKeepaliveCnt
  * (*Socket) SetTcpKeepaliveIdle
  * (*Socket) SetTcpKeepaliveIntvl
  * (*Socket) SetXpubVerbose
