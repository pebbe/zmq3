// A Go interface to ZeroMQ version 3.
//
// http://www.zeromq.org/
//
// Although you can check for event states POLLIN and POLLOUT,
// the ZeroMQ function zmq_poll() is not available by this
// package. You can use Go's select statement with the same
// effect.
package zmq3

/*
#include <zmq.h>
#include <stdlib.h>
#include <string.h>
#if ZMQ_VERSION_MAJOR < 3
#error ZeroMQ is too old. Version 3 required.
#endif

char *get_event(zmq_msg_t *msg, int *ev, int *val) {
    zmq_event_t event;
    memcpy(&event, zmq_msg_data(msg), sizeof(event));
    *ev = event.event;
    *val = event.data.connected.fd;
    return strdup(event.data.connected.addr);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
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

// Get 0MQ error message string.
func Error(e int) string {
	return C.GoString(C.zmq_strerror(C.int(e)))
}

//. Context

type Context struct {
	ctx unsafe.Pointer
}

// Context as string.
func (c Context) String() string {
	it, _ := c.GetIoThreads()
	ms, _ := c.GetMaxSockets()
	return fmt.Sprintf("Context(IO_THREADS=%d,MAX_SOCKETS=%d)", it, ms)
}

// Create new 0MQ context.
func NewContext() (ctx *Context, err error) {
	ctx = &Context{}
	c, e := C.zmq_ctx_new()
	if c == nil {
		err = errget(e)
	} else {
		ctx.ctx = c
		runtime.SetFinalizer(ctx, (*Context).Close)
	}
	return
}

// If not called explicitly, the context will be closed on garbage collection
func (ctx *Context) Close() error {
	i, err := C.zmq_ctx_destroy(ctx.ctx)
	if int(i) != 0 {
		return errget(err)
	}
	ctx.ctx = unsafe.Pointer(nil)
	return nil
}

func (ctx *Context) getOption(o C.int) (int, error) {
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

/*
Socket type as string.
*/
func (t Type) String() string {
	switch t {
	case REQ:
		return "REQ"
	case REP:
		return "REP"
	case DEALER:
		return "DEALER"
	case ROUTER:
		return "ROUTER"
	case PUB:
		return "PUB"
	case SUB:
		return "SUB"
	case XPUB:
		return "XPUB"
	case XSUB:
		return "XSUB"
	case PUSH:
		return "PUSH"
	case PULL:
		return "PULL"
	case PAIR:
		return "PAIR"
	}
	return "INVALID"
}

// Used by  (*Socket)Send() and (*Socket)Recv()
type Flag int

const (
	// Flags for (*Socket)Send(), (*Socket)Recv()
	// For Send, see: http://api.zeromq.org/3-2:zmq-send#toc2
	// For Recv, see: http://api.zeromq.org/3-2:zmq-msg-recv#toc2
	DONTWAIT = Flag(C.ZMQ_DONTWAIT)
	SNDMORE  = Flag(C.ZMQ_SNDMORE)
)

/*
Socket flag as string.
*/
func (f Flag) String() string {
	ff := make([]string, 0)
	if f&DONTWAIT != 0 {
		ff = append(ff, "DONTWAIT")
	}
	if f&SNDMORE != 0 {
		ff = append(ff, "SNDMORE")
	}
	return strings.Join(ff, "|")
}

// Used by (*Socket)Monitor() and (*Socket)RecvEvent()
type Event int

const (
	// Flags for (*Socket)Monitor() and (*Socket)RecvEvent()
	// See: http://api.zeromq.org/3-2:zmq-socket-monitor#toc2
	EVENT_ALL             = Event(C.ZMQ_EVENT_ALL)
	EVENT_CONNECTED       = Event(C.ZMQ_EVENT_CONNECTED)
	EVENT_CONNECT_DELAYED = Event(C.ZMQ_EVENT_CONNECT_DELAYED)
	EVENT_CONNECT_RETRIED = Event(C.ZMQ_EVENT_CONNECT_RETRIED)
	EVENT_LISTENING       = Event(C.ZMQ_EVENT_LISTENING)
	EVENT_BIND_FAILED     = Event(C.ZMQ_EVENT_BIND_FAILED)
	EVENT_ACCEPTED        = Event(C.ZMQ_EVENT_ACCEPTED)
	EVENT_ACCEPT_FAILED   = Event(C.ZMQ_EVENT_ACCEPT_FAILED)
	EVENT_CLOSED          = Event(C.ZMQ_EVENT_CLOSED)
	EVENT_CLOSE_FAILED    = Event(C.ZMQ_EVENT_CLOSE_FAILED)
	EVENT_DISCONNECTED    = Event(C.ZMQ_EVENT_DISCONNECTED)
)

/*
Socket event as string.
*/
func (e Event) String() string {
	if e == EVENT_ALL {
		return "EVENT_ALL"
	}
	ee := make([]string, 0)
	if e&EVENT_CONNECTED != 0 {
		ee = append(ee, "EVENT_CONNECTED")
	}
	if e&EVENT_CONNECT_DELAYED != 0 {
		ee = append(ee, "EVENT_CONNECT_DELAYED")
	}
	if e&EVENT_CONNECT_RETRIED != 0 {
		ee = append(ee, "EVENT_CONNECT_RETRIED")
	}
	if e&EVENT_LISTENING != 0 {
		ee = append(ee, "EVENT_LISTENING")
	}
	if e&EVENT_BIND_FAILED != 0 {
		ee = append(ee, "EVENT_BIND_FAILED")
	}
	if e&EVENT_ACCEPTED != 0 {
		ee = append(ee, "EVENT_ACCEPTED")
	}
	if e&EVENT_ACCEPT_FAILED != 0 {
		ee = append(ee, "EVENT_ACCEPT_FAILED")
	}
	if e&EVENT_CLOSED != 0 {
		ee = append(ee, "EVENT_CLOSED")
	}
	if e&EVENT_CLOSE_FAILED != 0 {
		ee = append(ee, "EVENT_CLOSE_FAILED")
	}
	if e&EVENT_DISCONNECTED != 0 {
		ee = append(ee, "EVENT_DISCONNECTED")
	}
	return strings.Join(ee, "|")
}

// Used by (soc *Socket)GetEvents()
type State int

const (
	// Flags for (*Socket)GetEvents()
	// See: http://api.zeromq.org/3-2:zmq-getsockopt#toc24
	POLLIN  = State(C.ZMQ_POLLIN)
	POLLOUT = State(C.ZMQ_POLLOUT)
)

/*
Socket state as string.
*/
func (s State) String() string {
	ss := make([]string, 0)
	if s&POLLIN != 0 {
		ss = append(ss, "POLLIN")
	}
	if s&POLLOUT != 0 {
		ss = append(ss, "POLLOUT")
	}
	return strings.Join(ss, "|")
}

/*
Socket functions starting with `Set` or `Get` are used for setting and
getting socket options.
*/
type Socket struct {
	soc unsafe.Pointer
}

/*
Socket as string.
*/
func (soc Socket) String() string {
	t, _ := soc.GetType()
	return "Socket(" + t.String() + ")"
}

/*
Create 0MQ socket.

For a description of socket types, see: http://api.zeromq.org/3-2:zmq-socket#toc3
*/
func (ctx *Context) NewSocket(t Type) (soc *Socket, err error) {
	soc = &Socket{}
	s, e := C.zmq_socket(ctx.ctx, C.int(t))
	if s == nil {
		err = errget(e)
	} else {
		soc.soc = s
		runtime.SetFinalizer(soc, (*Socket).Close)
	}
	return
}

// If not called explicitly, the socket will be closed on garbage collection
func (soc *Socket) Close() error {
	if i, err := C.zmq_close(soc.soc); int(i) != 0 {
		return errget(err)
	}
	soc.soc = unsafe.Pointer(nil)
	return nil
}

/*
Accept incoming connections on a socket.

For a description of endpoint, see: http://api.zeromq.org/3-2:zmq-bind#toc2
*/
func (soc *Socket) Bind(endpoint string) error {
	s := C.CString(endpoint)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_bind(soc.soc, s); int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Stop accepting connections on a socket.

For a description of endpoint, see: http://api.zeromq.org/3-2:zmq-bind#toc2
*/
func (soc *Socket) Unbind(endpoint string) error {
	s := C.CString(endpoint)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_unbind(soc.soc, s); int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Create outgoing connection from socket.

For a description of endpoint, see: http://api.zeromq.org/3-2:zmq-connect#toc2
*/
func (soc *Socket) Connect(endpoint string) error {
	s := C.CString(endpoint)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_connect(soc.soc, s); int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Disconnect a socket.

For a description of endpoint, see: http://api.zeromq.org/3-2:zmq-connect#toc2
*/
func (soc *Socket) Disconnect(endpoint string) error {
	s := C.CString(endpoint)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_disconnect(soc.soc, s); int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Receive a message part from a socket.

For a description of flags, see: http://api.zeromq.org/3-2:zmq-msg-recv#toc2
*/
func (soc *Socket) Recv(flags Flag) (string, error) {
	b, err := soc.RecvBytes(flags)
	return string(b), err
}

/*
Receive a message part from a socket.

For a description of flags, see: http://api.zeromq.org/3-2:zmq-msg-recv#toc2
*/
func (soc *Socket) RecvBytes(flags Flag) ([]byte, error) {
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
func (soc *Socket) Send(data string, flags Flag) (int, error) {
	return soc.SendBytes([]byte(data), flags)
}

/*
Send a message part on a socket.

For a description of flags, see: http://api.zeromq.org/3-2:zmq-send#toc2
*/
func (soc *Socket) SendBytes(data []byte, flags Flag) (int, error) {
	d := data
	if len(data) == 0 {
		d = []byte{0}
	}
	size, err := C.zmq_send(soc.soc, unsafe.Pointer(&d[0]), C.size_t(len(data)), C.int(flags))
	if size < 0 {
		return int(size), errget(err)
	}
	return int(size), nil
}

/*
Register a monitoring callback.

See: http://api.zeromq.org/3-2:zmq-socket-monitor#toc2

Example:

    package main

    import (
        zmq "github.com/pebbe/zmq3"
        "log"
        "time"
    )

    func rep_socket_monitor(ctx *zmq.Context, addr string) {
        s, err := ctx.NewSocket(zmq.PAIR)
        if err != nil {
            log.Fatalln(err)
        }
        err = s.Connect(addr)
        if err != nil {
            log.Fatalln(err)
        }
        for {
            a, b, c, err := s.RecvEvent(0)
            if err != nil {
                log.Println(err)
                break
            }
            log.Println(a, b, c)
        }
    }

    func main() {

        ctx, err := zmq.NewContext()
        if err != nil {
            log.Fatalln(err)
        }

        // this would hang because the socket created by Monitor() remains open
        //defer ctx.Close()

        // REP socket
        rep, err := ctx.NewSocket(zmq.REP)
        if err != nil {
            log.Fatalln(err)
        }
        defer rep.Close()

        // REP socket monitor, all events
        err = rep.Monitor("inproc://monitor.rep", zmq.EVENT_ALL)
        if err != nil {
            log.Fatalln(err)
        }
        go rep_socket_monitor(ctx, "inproc://monitor.rep")

        // Generate an event
        rep.Bind("tcp://*:5555")
        if err != nil {
            log.Fatalln(err)
        }

        // Allow some time for event detection
        time.Sleep(time.Second)
    }
*/
func (soc *Socket) Monitor(addr string, events Event) error {
	s := C.CString(addr)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_socket_monitor(soc.soc, s, C.int(events)); i != 0 {
		return errget(err)
	}
	return nil
}

/*
Receive a message part from a socket interpreted as an event.

For a description of flags, see: http://api.zeromq.org/3-2:zmq-msg-recv#toc2

For a description of event_type, see: http://api.zeromq.org/3-2:zmq-socket-monitor#toc2

For an example, see: func (*Socket) Monitor
*/
func (soc *Socket) RecvEvent(flags Flag) (event_type Event, addr string, value int, err error) {
	var msg C.zmq_msg_t
	if i, e := C.zmq_msg_init(&msg); i != 0 {
		err = errget(e)
		return
	}
	defer C.zmq_msg_close(&msg)

	size, e := C.zmq_msg_recv(&msg, soc.soc, C.int(flags))
	if size < 0 {
		err = errget(e)
		return
	}

	var t C.zmq_event_t
	if size < C.int(unsafe.Sizeof(t)) {
		err = errors.New("Not an event")
		return
	}

	et := C.int(0)
	val := C.int(0)
	addrs := C.get_event(&msg, &et, &val)
	defer C.free(unsafe.Pointer(addrs))

	event_type = Event(et)
	addr = C.GoString(addrs)
	value = int(val)

	return
}

/*
Start start built-in ØMQ proxy

See: http://api.zeromq.org/3-2:zmq-proxy
*/
func Proxy(frontend, backend, capture *Socket) error {
	var capt unsafe.Pointer
	if capture != nil {
		capt = capture.soc
	}
	_, err := C.zmq_proxy(frontend.soc, backend.soc, capt)
	return errget(err)
}
