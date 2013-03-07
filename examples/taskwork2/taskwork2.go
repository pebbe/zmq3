//
//  Task worker - design 2.
//  Adds pub-sub flow to receive and respond to kill signal
//
package main

import (
	zmq "github.com/pebbe/zmq3"

	"fmt"
	"strconv"
	"time"
)

func main() {
	//  Socket to receive messages on
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Connect("tcp://localhost:5557")

	//  Socket to send messages to
	sender, _ := zmq.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Connect("tcp://localhost:5558")

	//  Socket for control input
	controller, _ := zmq.NewSocket(zmq.SUB)
	defer controller.Close()
	controller.Connect("tcp://localhost:5559")
	controller.SetSubscribe("")

	//  Process messages from receiver and controller
	items := zmq.NewPoller()
	items.Register(receiver, zmq.POLLIN)
	items.Register(controller, zmq.POLLIN)
	//  Process messages from both sockets
	for {
		events, _ := items.Poll(-1)
		if events[0]&zmq.POLLIN != 0 {
			msg, _ := receiver.Recv(0)

			//  Do the work
			t, _ := strconv.Atoi(msg)
			time.Sleep(time.Duration(t) * time.Millisecond)

			//  Send results to sink
			sender.Send(msg, 0)

			//  Simple progress indicator for the viewer
			fmt.Printf(".")
		}
		//  Any controller command acts as 'KILL'
		if events[1]&zmq.POLLIN != 0 {
			break //  Exit loop
		}
	}
	fmt.Println()
}
