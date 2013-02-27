//
//  Task worker.
//  Connects PULL socket to tcp://localhost:5557
//  Collects workloads from ventilator via that socket
//  Connects PUSH socket to tcp://localhost:5558
//  Sends results to sink via that socket
//
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
	"strconv"
	"time"
)

func main() {
	context, _ := zmq.NewContext()

	//  Socket to receive messages on
	receiver, _ := context.NewSocket(zmq.PULL)
	receiver.Connect("tcp://localhost:5557")

	//  Socket to send messages to
	sender, _ := context.NewSocket(zmq.PUSH)
	sender.Connect("tcp://localhost:5558")

	//  Process tasks forever
	for {
		s, _ := receiver.Recv(0)

		//  Simple progress indicator for the viewer
        fmt.Print(s + ".");

		//  Do the work
		msec, _ := strconv.Atoi(s)
		time.Sleep(time.Duration(msec) * time.Millisecond)

		//  Send results to sink
		sender.Send("", 0)
	}
}
