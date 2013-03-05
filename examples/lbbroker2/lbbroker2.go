//
//  Load-balancing broker.
//  Demonstrates use of higher level functions.
//
package main

import (
	zmq "github.com/pebbe/zmq3"

	"fmt"
	"strings"
	"time"
)

const (
	NBR_CLIENTS  = 10
	NBR_WORKERS  = 3
	WORKER_READY = "\001" //  Signals worker is ready
)

//  Basic request-reply client using REQ socket
//
func client_task() {
	client, _ := zmq.NewSocket(zmq.REQ)
	defer client.Close()
	client.Connect("ipc://frontend.ipc")

	//  Send request, get reply
	for {
		client.SendMessage("HELLO")
		reply, _ := client.RecvMessage(0)
		if len(reply) == 0 {
			break
		}
		fmt.Println("Client:", strings.Join(reply, "\n\t"))
		time.Sleep(time.Second)
	}
}

//  Worker using REQ socket to do load-balancing
//
func worker_task() {
	worker, _ := zmq.NewSocket(zmq.REQ)
	defer worker.Close()
	worker.Connect("ipc://backend.ipc")

	//  Tell broker we're ready for work
	worker.SendMessage(WORKER_READY)

	//  Process messages as they arrive
	for {
		msg, e := worker.RecvMessage(0)
		if e != nil {
			break //  Interrupted ??
		}
		msg[len(msg)-1] = "OK"
		worker.SendMessage(msg)
	}
}

//  Now we come to the main task. This has the identical functionality to
//  the previous lbbroker example but uses higher level functions to read
//  and send messages:

func main() {
	//  Prepare our sockets
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	backend, _ := zmq.NewSocket(zmq.ROUTER)
	defer frontend.Close()
	defer backend.Close()
	frontend.Bind("ipc://frontend.ipc")
	backend.Bind("ipc://backend.ipc")

	client_nbr := 0
	for ; client_nbr < NBR_CLIENTS; client_nbr++ {
		go client_task()
	}
	for worker_nbr := 0; worker_nbr < NBR_WORKERS; worker_nbr++ {
		go worker_task()
	}

	//  Queue of available workers
	workers := make([]string, 0, 10)

	pool := func(soc *zmq.Socket, ch chan []string) {
		for {
			msg, err := soc.RecvMessage(0)
			if err != nil {
				panic(err)
			}
			ch <- msg
		}
	}
	chBack := make(chan []string, 100)
	chFrontIf := make(chan []string, 100)
	go pool(backend, chBack)
	go pool(frontend, chFrontIf)

	chNil := make(chan []string) // a channel that is never used
	for client_nbr > 0 {
		chFront := chNil
		//  Poll frontend only if we have available workers
		if len(workers) > 0 {
			chFront = chFrontIf
		}
		select {
		case msg := <-chBack:
			//  Use worker identity for load-balancing
			worker_id := msg[0]
			msg = msg[2:]

			workers = append(workers, worker_id)

			//  If client reply, send rest back to frontend
			if msg[0] != WORKER_READY {
				frontend.SendMessage(msg)
				client_nbr--
			}
		case msg := <-chFront:
			//  Route client request to first available worker
			backend.SendMessage(workers[0], "", msg)

			//  Dequeue and drop the next worker identity
			workers = workers[1:]
		}
	}

	time.Sleep(100 * time.Millisecond)
}
