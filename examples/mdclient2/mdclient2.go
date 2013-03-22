//
//  Majordomo Protocol client example - asynchronous.
//  Uses the mdcli API to hide all MDP aspects
//
//  Lets us build this source without creating a library
package main

import (
	"github.com/pebbe/zmq3/examples/mdapi"

	"fmt"
	"os"
)

func main() {
	var verbose bool
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		verbose = true
	}
	session, _ := mdapi.NewMdcli2("tcp://localhost:5555", verbose)

    var count int
    for count = 0; count < 100000; count++ {
		session.Send("echo", "Hello world")
    }
    for count = 0; count < 100000; count++ {
		_, err := session.Recv()
		if err != nil {
			break	      //  Interrupted by Ctrl-C
		}
    }
    fmt.Printf("%d replies received\n", count)
}
