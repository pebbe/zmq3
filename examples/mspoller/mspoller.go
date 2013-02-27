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

type msg struct {
	msg string
	err error
}

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

	chTask := make(chan msg)
	chWup := make(chan msg)
	go func() {
		for {
			s, err := receiver.Recv(0)
			chTask <- msg{s, err}
		}
	}()
	go func() {
		for {
			s, err := subscriber.Recv(0)
			chWup <- msg{s, err}
		}
	}()

	//  Process messages from both sockets
	for {
		select {
		case task, _ := <-chTask:
			//  Process task
			fmt.Println("Got task:", task.msg)
		case update, _ := <-chWup:
			//  Process weather update
			fmt.Println("Got weather update:", update.msg)
		}
	}

}
