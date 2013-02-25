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

func (soc *Socket) getString(opt C.int, bufsize int) (string, error) {
	if !soc.opened {
		return "", errSocClosed
	}
	value := make([]byte, bufsize)
	size := C.size_t(bufsize)
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value[0]), &size); i != 0 {
		return "", errget(err)
	}
	return string(value[:int(size)]), nil
}

func (soc *Socket) getInt(opt C.int) (int, error) {
	if !soc.opened {
		return 0, errSocClosed
	}
	value := C.int(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return int(value), nil
}

func (soc *Socket) getInt64(opt C.int) (int64, error) {
	if !soc.opened {
		return 0, errSocClosed
	}
	value := C.longlong(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return int64(value), nil
}

func (soc *Socket) getUInt64(opt C.int) (uint64, error) {
	if !soc.opened {
		return 0, errSocClosed
	}
	value := C.ulonglong(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return uint64(value), nil
}

// ZMQ_TYPE: Retrieve socket type
func (soc *Socket) GetType() (int, error) {
	return soc.getInt(C.ZMQ_TYPE)
}

// ZMQ_RCVMORE: More message data parts to follow
func (soc *Socket) GetRcvmore() (int, error) {
	return soc.getInt(C.ZMQ_RCVMORE)
}

// ZMQ_SNDHWM: Retrieves high water mark for outbound messages
func (soc *Socket) GetSndhwm() (int, error) {
	return soc.getInt(C.ZMQ_SNDHWM)
}

// ZMQ_RCVHWM: Retrieve high water mark for inbound messages
func (soc *Socket) GetRcvhwm() (int, error) {
	return soc.getInt(C.ZMQ_RCVHWM)
}

// ZMQ_AFFINITY: Retrieve I/O thread affinity
func (soc *Socket) GetAffinity() (uint64, error) {
	return soc.getUInt64(C.ZMQ_AFFINITY)
}

// ZMQ_IDENTITY: Set socket identity
func (soc *Socket) GetIdentity() (string, error) {
	return soc.getString(C.ZMQ_IDENTITY, 256)
}

// ZMQ_RATE: Retrieve multicast data rate
func (soc *Socket) GetRate() (int, error) {
	return soc.getInt(C.ZMQ_RATE)
}

// ZMQ_RECOVERY_IVL: Get multicast recovery interval
func (soc *Socket) GetRecoveryIvl() (int, error) {
	return soc.getInt(C.ZMQ_RECOVERY_IVL)
}

// ZMQ_SNDBUF: Retrieve kernel transmit buffer size
func (soc *Socket) GetSndbuf() (int, error) {
	return soc.getInt(C.ZMQ_SNDBUF)
}

// ZMQ_RCVBUF: Retrieve kernel receive buffer size
func (soc *Socket) GetRcvbuf() (int, error) {
	return soc.getInt(C.ZMQ_RCVBUF)
}

// ZMQ_LINGER: Retrieve linger period for socket shutdown
func (soc *Socket) GetLinger() (int, error) {
	return soc.getInt(C.ZMQ_LINGER)
}

// ZMQ_RECONNECT_IVL: Retrieve reconnection interval
func (soc *Socket) GetReconnectIvl() (int, error) {
	return soc.getInt(C.ZMQ_RECONNECT_IVL)
}

// ZMQ_RECONNECT_IVL_MAX: Retrieve maximum reconnection interval
func (soc *Socket) GetReconnectIvlMax() (int, error) {
	return soc.getInt(C.ZMQ_RECONNECT_IVL_MAX)
}

// ZMQ_BACKLOG: Retrieve maximum length of the queue of outstanding connections
func (soc *Socket) GetBacklog() (int, error) {
	return soc.getInt(C.ZMQ_BACKLOG)
}

// ZMQ_MAXMSGSIZE: Maximum acceptable inbound message size
func (soc *Socket) GetMaxmsgsize() (int64, error) {
	return soc.getInt64(C.ZMQ_MAXMSGSIZE)
}

// ZMQ_MULTICAST_HOPS: Maximum network hops for multicast packets
func (soc *Socket) GetMulticastHops() (int, error) {
	return soc.getInt(C.ZMQ_MULTICAST_HOPS)
}

// ZMQ_RCVTIMEO: Maximum time before a socket operation returns with EAGAIN
func (soc *Socket) GetRcvtimeo() (int, error) {
	return soc.getInt(C.ZMQ_RCVTIMEO)
}

// ZMQ_SNDTIMEO: Maximum time before a socket operation returns with EAGAIN
func (soc *Socket) GetSndtimeo() (int, error) {
	return soc.getInt(C.ZMQ_SNDTIMEO)
}

// ZMQ_IPV4ONLY: Retrieve IPv4-only socket override status
func (soc *Socket) GetIpv4only() (int, error) {
	return soc.getInt(C.ZMQ_IPV4ONLY)
}

// ZMQ_DELAY_ATTACH_ON_CONNECT: Retrieve attach-on-connect value
func (soc *Socket) GetDelayAttachOnConnect() (int, error) {
	return soc.getInt(C.ZMQ_DELAY_ATTACH_ON_CONNECT)
}

// ZMQ_FD: Retrieve file descriptor associated with the socket
// see socketget_unix.go and socketget_windows.go

// ZMQ_EVENTS: Retrieve socket event state
func (soc *Socket) GetEvents() (int, error) {
	return soc.getInt(C.ZMQ_EVENTS)
}

// ZMQ_LAST_ENDPOINT: Retrieve the last endpoint set
func (soc *Socket) GetLastEndpoint() (string, error) {
	return soc.getString(C.ZMQ_LAST_ENDPOINT, 1024)
}

// ZMQ_TCP_KEEPALIVE: Override SO_KEEPALIVE socket option
func (soc *Socket) GetTcpKeepalive() (int, error) {
	return soc.getInt(C.ZMQ_TCP_KEEPALIVE)
}

// ZMQ_TCP_KEEPALIVE_IDLE: Override TCP_KEEPCNT(or TCP_KEEPALIVE on some OS)
func (soc *Socket) GetTcpKeepaliveIdle() (int, error) {
	return soc.getInt(C.ZMQ_TCP_KEEPALIVE_IDLE)
}

// ZMQ_TCP_KEEPALIVE_CNT: Override TCP_KEEPCNT socket option
func (soc *Socket) GetTcpKeepaliveCnt() (int, error) {
	return soc.getInt(C.ZMQ_TCP_KEEPALIVE_CNT)
}

// ZMQ_TCP_KEEPALIVE_INTVL: Override TCP_KEEPINTVL socket option
func (soc *Socket) GetTcpKeepaliveIntvl() (int, error) {
	return soc.getInt(C.ZMQ_TCP_KEEPALIVE_INTVL)
}
