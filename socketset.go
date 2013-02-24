package zmq3

/*
#cgo pkg-config: libzmq
#include <zmq.h>
#include <stdlib.h>
#include <string.h>
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
		return err
	}
    return nil
}

func (soc *Socket) setInt(opt C.int, value int) error {
	if !soc.opened {
		return errSocClosed
	}
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&value), C.size_t(unsafe.Sizeof(value))); i != 0 {
		return err
	}
    return nil
}

func (soc *Socket) setInt64(opt C.int, value int64) error {
	if !soc.opened {
		return errSocClosed
	}
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&value), C.size_t(unsafe.Sizeof(value))); i != 0 {
		return err
	}
    return nil
}

func (soc *Socket) setUInt64(opt C.int, value uint64) error {
	if !soc.opened {
		return errSocClosed
	}
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&value), C.size_t(unsafe.Sizeof(value))); i != 0 {
		return err
	}
    return nil
}

// Set high water mark for outbound messages
func (soc *Socket) SetSndhwm(value int) error {
	return soc.setInt(C.ZMQ_SNDHWM, value)
}

// Set high water mark for inbound messages
func (soc *Socket) SetRcvhwm(value int) error {
	return soc.setInt(C.ZMQ_RCVHWM, value)
}

// Set I/O thread affinity
func (soc *Socket) SetAffinity(value uint64) error {
	return soc.setUInt64(C.ZMQ_AFFINITY, value)
}

// Establish message filter
func (soc *Socket) SetSubscribe(filter string) error {
	return soc.setString(C.ZMQ_SUBSCRIBE, filter)
}

// Remove message filter
func (soc *Socket) SetUnsubscribe(filter string) error {
	return soc.setString(C.ZMQ_UNSUBSCRIBE, filter)
}

// Set socket identity
func (soc *Socket) SetIdentity(value string) error {
	return soc.setString(C.ZMQ_IDENTITY, value)
}

// Set multicast data rate
func (soc *Socket) SetRate(value int) error {
	return soc.setInt(C.ZMQ_RATE, value)
}

// Set multicast recovery interval
func (soc *Socket) SetRecoveryIvl(value int) error {
	return soc.setInt(C.ZMQ_RECOVERY_IVL, value)
}

// Set kernel transmit buffer size
func (soc *Socket) SetSndbuf(value int) error {
	return soc.setInt(C.ZMQ_SNDBUF, value)
}

// Set kernel receive buffer size
func (soc *Socket) SetRcvbuf(value int) error {
	return soc.setInt(C.ZMQ_RCVBUF, value)
}

// Set linger period for socket shutdown
func (soc *Socket) SetLinger(value int) error {
	return soc.setInt(C.ZMQ_LINGER, value)
}

// Set reconnection interval
func (soc *Socket) SetReconnectIvl(value int) error {
	return soc.setInt(C.ZMQ_RECONNECT_IVL, value)
}

// Set maximum reconnection interval
func (soc *Socket) SetReconnectIvlMax(value int) error {
	return soc.setInt(C.ZMQ_RECONNECT_IVL_MAX, value)
}

// Set maximum length of the queue of outstanding connections
func (soc *Socket) SetBacklog(value int) error {
	return soc.setInt(C.ZMQ_BACKLOG, value)
}

// Maximum acceptable inbound message size
func (soc *Socket) SetMaxmsgsize(value int64) error {
	return soc.setInt64(C.ZMQ_MAXMSGSIZE, value)
}

// Maximum network hops for multicast packets
func (soc *Socket) SetMulticastHops(value int) error {
	return soc.setInt(C.ZMQ_MULTICAST_HOPS, value)
}

// Maximum time before a recv operation returns with EAGAIN
func (soc *Socket) SetRcvtimeo(value int) error {
	return soc.setInt(C.ZMQ_RCVTIMEO, value)
}

// Maximum time before a send operation returns with EAGAIN
func (soc *Socket) SetSndtimeo(value int) error {
	return soc.setInt(C.ZMQ_SNDTIMEO, value)
}

// Use IPv4-only sockets
func (soc *Socket) SetIPv4only(value int) error {
	return soc.setInt(C.ZMQ_IPV4ONLY, value)
}

// Accept messages only when connections are made
func (soc *Socket) SetDelayAttachOnConnect(value int) error {
	return soc.setInt(C.ZMQ_DELAY_ATTACH_ON_CONNECT, value)
}

// accept only routable messages on ROUTER sockets
func (soc *Socket) SetRouterMandatory(value int) error {
	return soc.setInt(C.ZMQ_ROUTER_MANDATORY, value)
}

// provide all subscription messages on XPUB sockets
func (soc *Socket) SetXpubVerbose(value int) error {
	return soc.setInt(C.ZMQ_XPUB_VERBOSE, value)
}

// Override SO_KEEPALIVE socket option
func (soc *Socket) SetTcpKeepalive(value int) error {
	return soc.setInt(C.ZMQ_TCP_KEEPALIVE, value)
}

// Override TCP_KEEPCNT(or TCP_KEEPALIVE on some OS)
func (soc *Socket) SetTcpKeepaliveIdle(value int) error {
	return soc.setInt(C.ZMQ_TCP_KEEPALIVE_IDLE, value)
}

// ZMQ_TCP_KEEPALIVE_CNT: Override TCP_KEEPCNT socket option
func (soc *Socket) SetTcpKeepaliveCnt(value int) error {
	return soc.setInt(C.ZMQ_TCP_KEEPALIVE_CNT, value)
}

// ZMQ_TCP_KEEPALIVE_INTVL: Override TCP_KEEPINTVL socket option
func (soc *Socket) SetTcpKeepaliveIntvl(value int) error {
	return soc.setInt(C.ZMQ_TCP_KEEPALIVE_INTVL, value)
}

/* TO DO:

   ZMQ_TCP_ACCEPT_FILTER: Assign filters to allow new TCP connections
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
