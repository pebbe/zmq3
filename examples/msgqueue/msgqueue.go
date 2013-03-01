//
//  Simple message queuing broker.
//  Same as request-reply broker but using QUEUE device
//
package main

import (
	zmq "github.com/pebbe/zmq3"
	"log"
)

func main() {
	var err error

	context, _ := zmq.NewContext()
	defer context.Close()

	//  Socket facing clients
	frontend, _ := context.NewSocket(zmq.ROUTER)
	defer frontend.Close()
	err = frontend.Bind("tcp://*:5559")
	if err != nil {
		log.Fatalln("Binding frontend:", err)
	}

    //  Socket facing services
	backend, _ := context.NewSocket(zmq.DEALER)
	defer backend.Close()
	err = backend.Bind("tcp://*:5560")
	if err != nil {
		log.Fatalln("Binding backend:", err)
	}

    //  Start the proxy
    err = zmq.Proxy(frontend, backend, nil)
	log.Fatalln("Proxy interrupted:", err)
}
