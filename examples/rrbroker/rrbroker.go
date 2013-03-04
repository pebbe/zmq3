//
//  Simple request-reply broker.
//
package main

import (
	zmq "github.com/pebbe/zmq3"
)

type Msg struct {
	s    string
	more bool
}

func main() {
	//  Prepare our sockets
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	defer frontend.Close()
	backend, _ := zmq.NewSocket(zmq.DEALER)
	defer backend.Close()
	frontend.Bind("tcp://*:5559")
	backend.Bind("tcp://*:5560")

	chFront := make(chan *Msg)
	chBack := make(chan *Msg)
	go func() {
		for {
			msg, _ := frontend.Recv(0)
			more, _ := frontend.GetRcvmore()
			chFront <- &Msg{msg, more}
		}
	}()
	go func() {
		for {
			msg, _ := backend.Recv(0)
			more, _ := backend.GetRcvmore()
			chBack <- &Msg{msg, more}
		}
	}()
	for {
		select {
		case msg := <-chFront:
			for {
				if msg.more {
					backend.Send(msg.s, zmq.SNDMORE)
				} else {
					backend.Send(msg.s, 0)
					break
				}
				msg = <-chFront
			}
		case msg := <-chBack:
			for {
				if msg.more {
					frontend.Send(msg.s, zmq.SNDMORE)
				} else {
					frontend.Send(msg.s, 0)
					break
				}
				msg = <-chBack
			}
		}
	}
}
