//
//  Reading from multiple sockets.
//  This version uses zmq.Poll()
//
package main

import (
	zmq "github.com/pebbe/zmq3"

	"fmt"
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

	//  Initialize poll set
	poller := zmq.NewPoller()
	poller.Register(receiver, zmq.POLLIN)
	poller.Register(subscriber, zmq.POLLIN)
	//  Process messages from both sockets
	for {
		items, _ := poller.Poll(-1)
		if items[0]&zmq.POLLIN != 0 {
			task, _ := receiver.Recv(0)
			//  Process task
			fmt.Println("Got task:", task)
		}
		if items[1]&zmq.POLLIN != 0 {
			update, _ := subscriber.Recv(0)
			//  Process weather update
			fmt.Println("Got weather update:", update)
		}
	}
}
