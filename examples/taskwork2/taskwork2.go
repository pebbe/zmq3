//
//  Task worker - design 2.
//  Adds pub-sub flow to receive and respond to kill signal
//
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
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
	chReceive := make(chan string)
	chControle := make(chan string)
	go func() {
		for {
			msg, e := receiver.Recv(0)
			if e != nil {
				break
			}
			chReceive <- msg
		}
	}()
	go func() {
		for {
			msg, e := controller.Recv(0)
			if e != nil {
				break
			}
			chControle <- msg
		}
	}()

	//  Process messages from both sockets
	for run := true; run; {
		select {
		case msg := <-chReceive:
			//  Do the work
			t, _ := strconv.Atoi(msg)
			time.Sleep(time.Duration(t) * time.Millisecond)

			//  Send results to sink
			sender.Send(msg, 0)

			//  Simple progress indicator for the viewer
			fmt.Printf(".")
		case <-chControle:
			//  Any controller command acts as 'KILL'
			run = false
		}
	}
	fmt.Println()
}
