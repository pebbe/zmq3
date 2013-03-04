//
//  Request-reply client.
//  Connects REQ socket to tcp://localhost:5559
//  Sends "Hello" to server, expects "World" back
//
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
)

func main() {
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://localhost:5559")

	for request := 0; request < 10; request++ {
		requester.Send("Hello", 0)
		reply, _ := requester.Recv(0)
		fmt.Printf("Received reply %d [%s]\n", request, reply)
	}
}