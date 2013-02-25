package zmq3

/*
#cgo pkg-config: libzmq
#include <zmq.h>
#include <stdint.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

func (soc *Socket) setString(opt C.int, s string) error {
	if !soc.opened {
		return errSocClosed
	}
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(cs), C.size_t(len(s))); i != 0 {
		return errget(err)
	}
    return nil
}

func (soc *Socket) setInt(opt C.int, value int) error {
	if !soc.opened {
		return errSocClosed
	}
	val := C.int(value)
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&val), C.size_t(unsafe.Sizeof(val))); i != 0 {
		return errget(err)
	}
    return nil
}

func (soc *Socket) setInt64(opt C.int, value int64) error {
	if !soc.opened {
		return errSocClosed
	}
	val := C.int64_t(value)
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&val), C.size_t(unsafe.Sizeof(val))); i != 0 {
		return errget(err)
	}
    return nil
}

func (soc *Socket) setUInt64(opt C.int, value uint64) error {
	if !soc.opened {
		return errSocClosed
	}
	val := C.uint64_t(value)
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&val), C.size_t(unsafe.Sizeof(val))); i != 0 {
		return errget(err)
	}
    return nil
}

// ZMQ_SNDHWM: Set high water mark for outbound messages
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc3
func (soc *Socket) SetSndhwm(value int) error {
	return soc.setInt(C.ZMQ_SNDHWM, value)
}

// ZMQ_RCVHWM: Set high water mark for inbound messages
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc4
func (soc *Socket) SetRcvhwm(value int) error {
	return soc.setInt(C.ZMQ_RCVHWM, value)
}

// ZMQ_AFFINITY: Set I/O thread affinity
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc5
func (soc *Socket) SetAffinity(value uint64) error {
	return soc.setUInt64(C.ZMQ_AFFINITY, value)
}

// ZMQ_SUBSCRIBE: Establish message filter
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc6
func (soc *Socket) SetSubscribe(filter string) error {
	return soc.setString(C.ZMQ_SUBSCRIBE, filter)
}

// ZMQ_UNSUBSCRIBE: Remove message filter
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc7
func (soc *Socket) SetUnsubscribe(filter string) error {
	return soc.setString(C.ZMQ_UNSUBSCRIBE, filter)
}

// ZMQ_IDENTITY: Set socket identity
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc8
func (soc *Socket) SetIdentity(value string) error {
	return soc.setString(C.ZMQ_IDENTITY, value)
}

// ZMQ_RATE: Set multicast data rate
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc9
func (soc *Socket) SetRate(value int) error {
	return soc.setInt(C.ZMQ_RATE, value)
}

// ZMQ_RECOVERY_IVL: Set multicast recovery interval
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc10
func (soc *Socket) SetRecoveryIvl(value int) error {
	return soc.setInt(C.ZMQ_RECOVERY_IVL, value)
}

// ZMQ_SNDBUF: Set kernel transmit buffer size
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc11
func (soc *Socket) SetSndbuf(value int) error {
	return soc.setInt(C.ZMQ_SNDBUF, value)
}

// ZMQ_RCVBUF: Set kernel receive buffer size
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc12
func (soc *Socket) SetRcvbuf(value int) error {
	return soc.setInt(C.ZMQ_RCVBUF, value)
}

// ZMQ_LINGER: Set linger period for socket shutdown
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc13
func (soc *Socket) SetLinger(value int) error {
	return soc.setInt(C.ZMQ_LINGER, value)
}

// ZMQ_RECONNECT_IVL: Set reconnection interval
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc14
func (soc *Socket) SetReconnectIvl(value int) error {
	return soc.setInt(C.ZMQ_RECONNECT_IVL, value)
}

// ZMQ_RECONNECT_IVL_MAX: Set maximum reconnection interval
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc15
func (soc *Socket) SetReconnectIvlMax(value int) error {
	return soc.setInt(C.ZMQ_RECONNECT_IVL_MAX, value)
}

// ZMQ_BACKLOG: Set maximum length of the queue of outstanding connections
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc16
func (soc *Socket) SetBacklog(value int) error {
	return soc.setInt(C.ZMQ_BACKLOG, value)
}

// ZMQ_MAXMSGSIZE: Maximum acceptable inbound message size
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc17
func (soc *Socket) SetMaxmsgsize(value int64) error {
	return soc.setInt64(C.ZMQ_MAXMSGSIZE, value)
}

// ZMQ_MULTICAST_HOPS: Maximum network hops for multicast packets
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc18
func (soc *Socket) SetMulticastHops(value int) error {
	return soc.setInt(C.ZMQ_MULTICAST_HOPS, value)
}

// ZMQ_RCVTIMEO: Maximum time before a recv operation returns with EAGAIN
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc19
func (soc *Socket) SetRcvtimeo(value int) error {
	return soc.setInt(C.ZMQ_RCVTIMEO, value)
}

// ZMQ_SNDTIMEO: Maximum time before a send operation returns with EAGAIN
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc20
func (soc *Socket) SetSndtimeo(value int) error {
	return soc.setInt(C.ZMQ_SNDTIMEO, value)
}

// ZMQ_IPV4ONLY: Use IPv4-only sockets
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc21
func (soc *Socket) SetIpv4only(value int) error {
	return soc.setInt(C.ZMQ_IPV4ONLY, value)
}

// ZMQ_DELAY_ATTACH_ON_CONNECT: Accept messages only when connections are made
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc22
func (soc *Socket) SetDelayAttachOnConnect(value int) error {
	return soc.setInt(C.ZMQ_DELAY_ATTACH_ON_CONNECT, value)
}

// ZMQ_ROUTER_MANDATORY: accept only routable messages on ROUTER sockets
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc23
func (soc *Socket) SetRouterMandatory(value int) error {
	return soc.setInt(C.ZMQ_ROUTER_MANDATORY, value)
}

// ZMQ_XPUB_VERBOSE: provide all subscription messages on XPUB sockets
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc24
func (soc *Socket) SetXpubVerbose(value int) error {
	return soc.setInt(C.ZMQ_XPUB_VERBOSE, value)
}

// ZMQ_TCP_KEEPALIVE: Override SO_KEEPALIVE socket option
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc25
func (soc *Socket) SetTcpKeepalive(value int) error {
	return soc.setInt(C.ZMQ_TCP_KEEPALIVE, value)
}

// ZMQ_TCP_KEEPALIVE_IDLE: Override TCP_KEEPCNT(or TCP_KEEPALIVE on some OS)
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc26
func (soc *Socket) SetTcpKeepaliveIdle(value int) error {
	return soc.setInt(C.ZMQ_TCP_KEEPALIVE_IDLE, value)
}

// ZMQ_TCP_KEEPALIVE_CNT: ZMQ_TCP_KEEPALIVE_CNT: Override TCP_KEEPCNT socket option
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc27
func (soc *Socket) SetTcpKeepaliveCnt(value int) error {
	return soc.setInt(C.ZMQ_TCP_KEEPALIVE_CNT, value)
}

// ZMQ_TCP_KEEPALIVE_INTVL: Override TCP_KEEPINTVL socket option
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc28
func (soc *Socket) SetTcpKeepaliveIntvl(value int) error {
	return soc.setInt(C.ZMQ_TCP_KEEPALIVE_INTVL, value)
}

/* TO DO:

// ZMQ_TCP_ACCEPT_FILTER: Assign filters to allow new TCP connections
//
// See: http://api.zeromq.org/3-2:zmq-setsockopt#toc29

       Assign arbitrary number of filters that will be applied for each new TCP transport connection on a
       listening socket. If no filters applied, then TCP transport allows connections from any ip. If at
       least one filter is applied then new connection source ip should be matched. To clear all filters
       call zmq_setsockopt(socket, ZMQ_TCP_ACCEPT_FILTER, NULL, 0). Filter is a null-terminated string with
       ipv6 or ipv4 CIDR.


       Option value type         binary data

       Option value unit         N/A

       Default value             no filters (allow from all)

       Applicable socket types   all listening sockets, when using
                                 TCP transports.
*/
