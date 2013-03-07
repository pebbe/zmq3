package zmq3

/*
#include <zmq.h>
*/
import "C"

import (
	"time"
)

type Poller struct {
	items []C.zmq_pollitem_t
}

// Create a new Poller
func NewPoller() *Poller {
	return &Poller{items: make([]C.zmq_pollitem_t, 0)}
}

// Add items to the poller
//
// Events is a bitwise OR of zmq.POLLIN and zmq.POLLOUT
func (p *Poller) Register(soc *Socket, events State) {
	var item C.zmq_pollitem_t
	item.socket = soc.soc
	item.fd = 0
	item.events = C.short(events)
	p.items = append(p.items, item)
}

/*
Input/output multiplexing

If timeout < 0, wait forever until a matching event is detected

Example:

    poller := zmq.NewPoller()
    poller.Register(socket0, zmq.POLLIN)
    poller.Register(socket1, zmq.POLLIN)
    //  Process messages from both sockets
    for {
        items, _ := poller.Poll(-1)
        if items[0]&zmq.POLLIN != 0 {
            msg, _ := socket0.Recv(0)
            //  Process msg
        }
        if items[1]&zmq.POLLIN != 0 {
            msg, _ := socket1.Recv(0)
            //  Process msg
        }
    }
*/
func (p *Poller) Poll(timeout time.Duration) ([]State, error) {
	t := timeout
	if t > 0 {
		t = t / time.Millisecond
	}
	if t < 0 {
		t = -1
	}
	rv, err := C.zmq_poll(&p.items[0], C.int(len(p.items)), C.long(t))
	if rv < 0 {
		return []State{}, errget(err)
	}
	states := make([]State, len(p.items))
	for i, it := range p.items {
		states[i] = State(it.revents)
	}
	return states, nil
}
