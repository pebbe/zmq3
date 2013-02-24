//
//  Hello World server
//  Binds REP socket to tcp://*:5555
//  Expects "Hello" from client, replies with "World"
//

package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
	"time"
)

func main() {
	context, _ := zmq.NewContext()

	//  Socket to talk to clients
	responder, _ := context.NewSocket(zmq.REP)
	responder.Bind("tcp://*:5555")

	for {
		//  Wait for next request from client
		msg, _ := responder.Recv(0)
		fmt.Println("Received ", string(msg))

		//  Do some 'work'
		time.Sleep(time.Second)

		//  Send reply back to client
		reply := "World"
		responder.Send([]byte(reply), 0)
	}
}
