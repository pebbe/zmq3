package zmq3

/*
#include <zmq.h>
#include "zmq2.h"
#include <stdint.h>
*/
import "C"

import (
	"time"
	"unsafe"
)

func (soc *Socket) getString(opt C.int, bufsize int) (string, error) {
	value := make([]byte, bufsize)
	size := C.size_t(bufsize)
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value[0]), &size); i != 0 {
		return "", errget(err)
	}
	return string(value[:int(size)]), nil
}

func (soc *Socket) getInt(opt C.int) (int, error) {
	value := C.int(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return int(value), nil
}

func (soc *Socket) getInt64(opt C.int) (int64, error) {
	value := C.int64_t(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return int64(value), nil
}

func (soc *Socket) getUInt64(opt C.int) (uint64, error) {
	value := C.uint64_t(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return uint64(value), nil
}

func (soc *Socket) getUInt32(opt C.int) (uint32, error) {
	value := C.uint32_t(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return uint32(value), nil
}

// ZMQ_TYPE: Retrieve socket type
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc3
func (soc *Socket) GetType() (Type, error) {
	v, err := soc.getInt(C.ZMQ_TYPE)
	return Type(v), err
}

// ZMQ_RCVMORE: More message data parts to follow
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc4
func (soc *Socket) GetRcvmore() (bool, error) {
	if major < 3 {
		v, err := soc.getInt64(C.ZMQ_RCVMORE)
		return v != 0, err
	}
	v, err := soc.getInt(C.ZMQ_RCVMORE)
	return v != 0, err
}

// ZMQ_SNDHWM: Retrieves high water mark for outbound messages
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc5
//
// In 0MQ version 2, returns common value for both inbound and outbound messages
func (soc *Socket) GetSndhwm() (int, error) {
	if major < 3 {
		i, e := soc.getUInt64(C.ZMQ_SNDHWM)
		return int(i), e
	}
	return soc.getInt(C.ZMQ_SNDHWM)
}

// ZMQ_RCVHWM: Retrieve high water mark for inbound messages
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc6
//
// In 0MQ version 2, returns common value for both inbound and outbound messages
func (soc *Socket) GetRcvhwm() (int, error) {
	if major < 3 {
		i, e := soc.getUInt64(C.ZMQ_RCVHWM)
		return int(i), e
	}
	return soc.getInt(C.ZMQ_RCVHWM)
}

// ZMQ_AFFINITY: Retrieve I/O thread affinity
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc7
func (soc *Socket) GetAffinity() (uint64, error) {
	return soc.getUInt64(C.ZMQ_AFFINITY)
}

// ZMQ_IDENTITY: Set socket identity
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc8
func (soc *Socket) GetIdentity() (string, error) {
	return soc.getString(C.ZMQ_IDENTITY, 256)
}

// ZMQ_RATE: Retrieve multicast data rate
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc9
func (soc *Socket) GetRate() (int, error) {
	if major < 3 {
		i, e := soc.getInt64(C.ZMQ_RATE)
		return int(i), e
	}
	return soc.getInt(C.ZMQ_RATE)
}

// ZMQ_RECOVERY_IVL: Get multicast recovery interval
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc10
func (soc *Socket) GetRecoveryIvl() (time.Duration, error) {
	if major < 3 {
		v, e := soc.getInt64(C.ZMQ_RECOVERY_IVL_MSEC)
		if v >= 0 {
			return time.Duration(v) * time.Millisecond, e
		}
		v, e = soc.getInt64(C.ZMQ_RECOVERY_IVL)
		return time.Duration(v) * time.Second, e
	}
	v, err := soc.getInt(C.ZMQ_RECOVERY_IVL)
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_SNDBUF: Retrieve kernel transmit buffer size
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc11
func (soc *Socket) GetSndbuf() (int, error) {
	if major < 3 {
		i, e := soc.getUInt64(C.ZMQ_SNDBUF)
		return int(i), e
	}
	return soc.getInt(C.ZMQ_SNDBUF)
}

// ZMQ_RCVBUF: Retrieve kernel receive buffer size
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc12
func (soc *Socket) GetRcvbuf() (int, error) {
	if major < 3 {
		i, e := soc.getUInt64(C.ZMQ_RCVBUF)
		return int(i), e
	}
	return soc.getInt(C.ZMQ_RCVBUF)
}

// ZMQ_LINGER: Retrieve linger period for socket shutdown
//
// Returns time.Duration(-1) for infinite
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc13
func (soc *Socket) GetLinger() (time.Duration, error) {
	v, err := soc.getInt(C.ZMQ_LINGER)
	if v < 0 {
		return time.Duration(-1), err
	}
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_RECONNECT_IVL: Retrieve reconnection interval
//
// Returns time.Duration(-1) for no reconnection
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc14
func (soc *Socket) GetReconnectIvl() (time.Duration, error) {
	v, err := soc.getInt(C.ZMQ_RECONNECT_IVL)
	if v < 0 {
		return time.Duration(-1), err
	}
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_RECONNECT_IVL_MAX: Retrieve maximum reconnection interval
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc15
func (soc *Socket) GetReconnectIvlMax() (time.Duration, error) {
	v, err := soc.getInt(C.ZMQ_RECONNECT_IVL_MAX)
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_BACKLOG: Retrieve maximum length of the queue of outstanding connections
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc16
func (soc *Socket) GetBacklog() (int, error) {
	return soc.getInt(C.ZMQ_BACKLOG)
}

// ZMQ_MAXMSGSIZE: Maximum acceptable inbound message size
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc17
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetMaxmsgsize() (int64, error) {
	if major < 3 {
		return -1, ErrorNotImplemented
	}
	return soc.getInt64(C.ZMQ_MAXMSGSIZE)
}

// ZMQ_MULTICAST_HOPS: Maximum network hops for multicast packets
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc18
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetMulticastHops() (int, error) {
	if major < 3 {
		return -1, ErrorNotImplemented
	}
	return soc.getInt(C.ZMQ_MULTICAST_HOPS)
}

// ZMQ_RCVTIMEO: Maximum time before a socket operation returns with EAGAIN
//
// Returns time.Duration(-1) for infinite
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc19
func (soc *Socket) GetRcvtimeo() (time.Duration, error) {
	v, err := soc.getInt(C.ZMQ_RCVTIMEO)
	if v < 0 {
		return time.Duration(-1), err
	}
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_SNDTIMEO: Maximum time before a socket operation returns with EAGAIN
//
// Returns time.Duration(-1) for infinite
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc20
func (soc *Socket) GetSndtimeo() (time.Duration, error) {
	v, err := soc.getInt(C.ZMQ_SNDTIMEO)
	if v < 0 {
		return time.Duration(-1), err
	}
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_IPV4ONLY: Retrieve IPv4-only socket override status
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc21
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetIpv4only() (bool, error) {
	if major < 3 {
		return false, ErrorNotImplemented
	}
	v, err := soc.getInt(C.ZMQ_IPV4ONLY)
	return v != 0, err
}

// ZMQ_DELAY_ATTACH_ON_CONNECT: Retrieve attach-on-connect value
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc22
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetDelayAttachOnConnect() (bool, error) {
	if major < 3 {
		return false, ErrorNotImplemented
	}
	v, err := soc.getInt(C.ZMQ_DELAY_ATTACH_ON_CONNECT)
	return v != 0, err
}

// ZMQ_FD: Retrieve file descriptor associated with the socket
// see socketget_unix.go and socketget_windows.go

// ZMQ_EVENTS: Retrieve socket event state
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc24
func (soc *Socket) GetEvents() (State, error) {
	if major < 3 {
		v, err := soc.getUInt32(C.ZMQ_EVENTS)
		return State(v), err
	}
	v, err := soc.getInt(C.ZMQ_EVENTS)
	return State(v), err
}

// ZMQ_LAST_ENDPOINT: Retrieve the last endpoint set
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc25
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetLastEndpoint() (string, error) {
	if major < 3 {
		return "", ErrorNotImplemented
	}
	return soc.getString(C.ZMQ_LAST_ENDPOINT, 1024)
}

// ZMQ_TCP_KEEPALIVE: Override SO_KEEPALIVE socket option
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc26
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetTcpKeepalive() (int, error) {
	if major < 3 {
		return -1, ErrorNotImplemented
	}
	return soc.getInt(C.ZMQ_TCP_KEEPALIVE)
}

// ZMQ_TCP_KEEPALIVE_IDLE: Override TCP_KEEPCNT(or TCP_KEEPALIVE on some OS)
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc27
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetTcpKeepaliveIdle() (int, error) {
	if major < 3 {
		return -1, ErrorNotImplemented
	}
	return soc.getInt(C.ZMQ_TCP_KEEPALIVE_IDLE)
}

// ZMQ_TCP_KEEPALIVE_CNT: Override TCP_KEEPCNT socket option
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc28
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetTcpKeepaliveCnt() (int, error) {
	if major < 3 {
		return -1, ErrorNotImplemented
	}
	return soc.getInt(C.ZMQ_TCP_KEEPALIVE_CNT)
}

// ZMQ_TCP_KEEPALIVE_INTVL: Override TCP_KEEPINTVL socket option
//
// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc29
//
// Returns ErrorNotImplemented in 0MQ version 2.
func (soc *Socket) GetTcpKeepaliveIntvl() (int, error) {
	if major < 3 {
		return -1, ErrorNotImplemented
	}
	return soc.getInt(C.ZMQ_TCP_KEEPALIVE_INTVL)
}
