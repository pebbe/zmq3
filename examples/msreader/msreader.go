//
//  Reading from multiple sockets.
//  This version uses a simple recv loop
//
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
	"time"
)

func main() {

	//  Prepare our context and sockets
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

	//  Process messages from both sockets
	//  We prioritize traffic from the task ventilator
	for {

		//  Process any waiting tasks
		for {
			task, err := receiver.Recv(zmq.DONTWAIT)
			if err != nil {
				break
			}
			//  process task
			fmt.Println("Got task:", task)
		}

		//  Process any waiting weather updates
		for {
			udate, err := subscriber.Recv(zmq.DONTWAIT)
			if err != nil {
				break
			}
			//  process weather update
			fmt.Println("Got weather update:", udate)
		}

		//  No activity, so sleep for 1 msec
		time.Sleep(time.Millisecond)
	}
}
