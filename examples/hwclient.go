//
//  Hello World client.
//  Connects REQ socket to tcp://localhost:5555
//  Sends "Hello" to server, expects "World" back
//
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
)

func main() {
	context, _ := zmq.NewContext()

	//  Socket to talk to server
	fmt.Println("Connecting to hello world server…")
	requester, _ := context.NewSocket(zmq.REQ)
	requester.Connect("tcp://localhost:5555")

	for request_nbr := 0; request_nbr != 10; request_nbr++ {
		// send hello
		msg := fmt.Sprintf("Hello %d", request_nbr)
		fmt.Println("Sending ", msg)
		requester.Send([]byte(msg), 0)

		// Wait for reply:
		reply, _ := requester.Recv(0)
		fmt.Println("Received ", string(reply))
	}
}
