// A Go interface to ZeroMQ version 3.
//
// http://www.zeromq.org/
package zmq3

/*
#cgo pkg-config: libzmq
#include <zmq.h>
#include <stdlib.h>
#include <string.h>
#if ZMQ_VERSION_MAJOR < 3
#error ZeroMQ is too old. Version 3 required.
#endif
*/
import "C"

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"
)

//. Util

func errget(err error) error {
	errno, ok := err.(syscall.Errno)
	if ok && errno >= C.ZMQ_HAUSNUMERO {
		return errors.New(C.GoString(C.zmq_strerror(C.int(errno))))
	}
	return err
}

// Report 0MQ library version.
func Version() (int, int, int) {
	var major, minor, patch C.int
	C.zmq_version(&major, &minor, &patch)
	return int(major), int(minor), int(patch)
}

//. Context

var (
	errCtxClosed = errors.New("Context is closed")
)

type Context struct {
	ctx    unsafe.Pointer
	opened bool
}

// Create new 0MQ context.
func NewContext() (ctx *Context, err error) {
	ctx = &Context{}
	c, e := C.zmq_ctx_new()
	if c == nil {
		err = errget(e)
	} else {
		ctx.ctx = c
		ctx.opened = true
		runtime.SetFinalizer(ctx, (*Context).Close)
	}
	return
}

// If not called explicitly, the context will be closed on garbage collection
func (ctx *Context) Close() error {
	if ctx.opened {
		ctx.opened = false
		i, err := C.zmq_ctx_destroy(ctx.ctx)
		if int(i) != 0 {
			return errget(err)
		}
	}
	return nil
}

func (ctx *Context) getOption(o C.int) (int, error) {
	if !ctx.opened {
		return -1, errCtxClosed
	}
	n, err := C.zmq_ctx_get(ctx.ctx, o)
	return int(n), errget(err)
}

// Returns the size of the 0MQ thread pool for this context.
func (ctx *Context) GetIoThreads() (int, error) {
	return ctx.getOption(C.ZMQ_IO_THREADS)
}

// Returns the maximum number of sockets allowed for this context.
func (ctx *Context) GetMaxSockets() (int, error) {
	return ctx.getOption(C.ZMQ_MAX_SOCKETS)
}

func (ctx *Context) setOption(o C.int, n int) error {
	if !ctx.opened {
		return errCtxClosed
	}
	i, err := C.zmq_ctx_set(ctx.ctx, o, C.int(n))
	if int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Specifies the size of the 0MQ thread pool to handle I/O operations. If
your application is using only the inproc transport for messaging you may set this to zero,
otherwise set it to at least one. This option only applies before creating any sockets on the
context.

Default value   1
*/
func (ctx *Context) SetIoThreads(n int) error {
	return ctx.setOption(C.ZMQ_IO_THREADS, n)
}

/*
Sets the maximum number of sockets allowed on the context.

Default value   1024
*/
func (ctx *Context) SetMaxSockets(n int) error {
	return ctx.setOption(C.ZMQ_MAX_SOCKETS, n)
}

//. Sockets

// Specifies the type of a socket, used by (*Context)NewSocket()
type Type int

const (
	// Constants for (*Context)NewSocket()
	// See: http://api.zeromq.org/3-2:zmq-socket#toc3
	REQ    = Type(C.ZMQ_REQ)
	REP    = Type(C.ZMQ_REP)
	DEALER = Type(C.ZMQ_DEALER)
	ROUTER = Type(C.ZMQ_ROUTER)
	PUB    = Type(C.ZMQ_PUB)
	SUB    = Type(C.ZMQ_SUB)
	XPUB   = Type(C.ZMQ_XPUB)
	XSUB   = Type(C.ZMQ_XSUB)
	PUSH   = Type(C.ZMQ_PUSH)
	PULL   = Type(C.ZMQ_PULL)
	PAIR   = Type(C.ZMQ_PAIR)
)

// Used by  (*Socket)Send() and (*Socket)Recv()
type Flag int

const (
	// Flags for (*Socket)Send(), (*Socket)Recv()
	// See zmq_send: http://api.zeromq.org/3-2:zmq-send#toc2
	// See zmq_msg_recv: http://api.zeromq.org/3-2:zmq-msg-recv#toc2
	DONTWAIT = Flag(C.ZMQ_DONTWAIT)
	SNDMORE  = Flag(C.ZMQ_SNDMORE)
)

// Used by (soc *Socket)GetEvents()
type EventState int

const (
	POLLIN  = EventState(C.ZMQ_POLLIN)
	POLLOUT = EventState(C.ZMQ_POLLOUT)
)

var (
	errSocClosed = errors.New("Socket is closed")
)

/*
Socket functions starting with `Set` or `Get` are used for setting and
getting socket options.
*/
type Socket struct {
	ctx    *Context
	soc    unsafe.Pointer
	opened bool
}

/*
Create 0MQ socket.

For a description of socket types, see: http://api.zeromq.org/3-2:zmq-socket#toc3
*/
func (ctx *Context) NewSocket(t Type) (soc *Socket, err error) {
	soc = &Socket{}
	if !ctx.opened {
		err = errCtxClosed
		return
	}
	s, e := C.zmq_socket(ctx.ctx, C.int(t))
	if s == nil {
		err = errget(e)
	} else {
		soc.ctx = ctx
		soc.soc = s
		soc.opened = true
		runtime.SetFinalizer(soc, (*Socket).Close)
	}
	return
}

// If not called explicitly, the socket will be closed on garbage collection
func (soc *Socket) Close() error {
	if soc.opened {
		soc.opened = false
		if i, err := C.zmq_close(soc.soc); int(i) != 0 {
			return errget(err)
		}
	}
	return nil
}

/*
Accept incoming connections on a socket.

For a description of endpoint, see: http://api.zeromq.org/3-2:zmq-bind#toc2
*/
func (soc *Socket) Bind(endpoint string) error {
	if !soc.opened {
		return errSocClosed
	}
	s := C.CString(endpoint)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_bind(soc.soc, s); int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Create outgoing connection from socket.

For a description of endpoint, see: http://api.zeromq.org/3-2:zmq-connect#toc2
*/
func (soc *Socket) Connect(endpoint string) error {
	if !soc.opened {
		return errSocClosed
	}
	s := C.CString(endpoint)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_connect(soc.soc, s); int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Receive a message part from a socket.

For a description of flags, see: http://api.zeromq.org/3-2:zmq-msg-recv#toc2
*/
func (soc *Socket) Recv(flags Flag) ([]byte, error) {
	if !soc.opened {
		return []byte{}, errSocClosed
	}

	var msg C.zmq_msg_t
	if i, err := C.zmq_msg_init(&msg); i != 0 {
		return []byte{}, errget(err)
	}
	defer C.zmq_msg_close(&msg)

	size, err := C.zmq_msg_recv(&msg, soc.soc, C.int(flags))
	if size < 0 {
		return []byte{}, errget(err)
	}
	if size == 0 {
		return []byte{}, nil
	}
	data := make([]byte, int(size))
	C.memcpy(unsafe.Pointer(&data[0]), C.zmq_msg_data(&msg), C.size_t(size))
	return data, nil
}

/*
Send a message part on a socket.

For a description of flags, see: http://api.zeromq.org/3-2:zmq-send#toc2
*/
func (soc *Socket) Send(data []byte, flags Flag) (int, error) {
	if !soc.opened {
		return -1, errSocClosed
	}
	size, err := C.zmq_send(soc.soc, unsafe.Pointer(&data[0]), C.size_t(len(data)), C.int(flags))
	if size < 0 {
		return int(size), errget(err)
	}
	return int(size), nil
}
