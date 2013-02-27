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

	context, _ := zmq.NewContext()
	defer context.Close()

	//  Connect to task ventilator
	receiver, _ := context.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Connect("tcp://localhost:5557")

	//  Connect to weather server
	subscriber, _ := context.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5556")
	subscriber.SetSubscribe("10001 ")

	chTask := make(chan string)
	chWup := make(chan string)
	go func() {
		for {
			msg, _ := receiver.Recv(0)
			chTask <- msg
		}
	}()
	go func() {
		for {
			msg, _ := subscriber.Recv(0)
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
