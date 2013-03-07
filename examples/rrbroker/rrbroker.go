//
//  Simple request-reply broker.
//
package main

import (
	zmq "github.com/pebbe/zmq3"
)

func main() {
	//  Prepare our sockets
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	defer frontend.Close()
	backend, _ := zmq.NewSocket(zmq.DEALER)
	defer backend.Close()
	frontend.Bind("tcp://*:5559")
	backend.Bind("tcp://*:5560")

	//  Initialize poll set
	items := zmq.NewPoller()
	items.Register(frontend, zmq.POLLIN)
	items.Register(backend, zmq.POLLIN)

	//  Switch messages between sockets
	for {
		events, _ := items.Poll(-1)
		if events[0]&zmq.POLLIN != 0 {
			for {
				msg, _ := frontend.Recv(0)
				if more, _ := frontend.GetRcvmore(); more {
					backend.Send(msg, zmq.SNDMORE)
				} else {
					backend.Send(msg, 0)
					break
				}
			}
		}
		if events[1]&zmq.POLLIN != 0 {
			for {
				msg, _ := backend.Recv(0)
				if more, _ := backend.GetRcvmore(); more {
					frontend.Send(msg, zmq.SNDMORE)
				} else {
					frontend.Send(msg, 0)
					break
				}
			}
		}
	}
}
