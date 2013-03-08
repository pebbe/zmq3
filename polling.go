package zmq3

/*
#include <zmq.h>
*/
import "C"

import (
	"fmt"
	"time"
)

type Poller struct {
	items []C.zmq_pollitem_t
	socks []*Socket
}

// Create a new Poller
func NewPoller() *Poller {
	return &Poller{
		items: make([]C.zmq_pollitem_t, 0),
		socks: make([]*Socket, 0)}
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
	p.socks = append(p.socks, soc)
}

/*
Input/output multiplexing

If timeout < 0, wait forever until a matching event is detected

Only sockets with matching socket states are returned in the map.
The values in the map are the actual sockets states.

Example:

    poller := zmq.NewPoller()
    poller.Register(socket0, zmq.POLLIN)
    poller.Register(socket1, zmq.POLLIN)
    //  Process messages from both sockets
    for {
        sockets, _ := poller.Poll(-1)
        for socket := range sockets {
            switch socket {
            case socket0:
                msg, _ := socket0.Recv(0)
                //  Process msg
            case socket1:
                msg, _ := socket1.Recv(0)
                //  Process msg
            }
        }
    }
*/
func (p *Poller) Poll(timeout time.Duration) (map[*Socket]State, error) {
	mp := make(map[*Socket]State)
	t := timeout
	if t > 0 {
		t = t / time.Millisecond
	}
	if t < 0 {
		t = -1
	}
	rv, err := C.zmq_poll(&p.items[0], C.int(len(p.items)), C.long(t))
	if rv < 0 {
		return mp, errget(err)
	}
	for i, it := range p.items {
		if it.events&it.revents != 0 {
			mp[p.socks[i]] = State(it.revents)
		}
	}
	return mp, nil
}

func (p *Poller) String() string {
	str := make([]string, 0)
	for i, poll := range p.items {
		str = append(str, fmt.Sprintf("%v%v", p.socks[i], State(poll.events)))
	}
	return fmt.Sprint("Poller", str)
}
