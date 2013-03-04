//
//  Reading from multiple sockets.
//  This version does NOT zmq_poll()
//  It uses Go's select instead.
//
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
)

func main() {

	//  Connect to task ventilator
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Connect("tcp://localhost:5557")

	//  Connect to weather server
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5556")
	subscriber.SetSubscribe("10001 ")

	chTask := make(chan string)
	chWup := make(chan string)
	go func() {
		for {
			msg, e := receiver.Recv(0)
			if e != nil {
				break
			}
			chTask <- msg
		}
	}()
	go func() {
		for {
			msg, e := subscriber.Recv(0)
			if e != nil {
				break
			}
			chWup <- msg
		}
	}()

	//  Process messages from both sockets
	for {
		select {
		case task := <-chTask:
			//  Process task
			fmt.Println("Got task:", task)
		case update := <-chWup:
			//  Process weather update
			fmt.Println("Got weather update:", update)
		}
	}
}
