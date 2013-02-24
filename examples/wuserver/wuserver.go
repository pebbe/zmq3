//
//  Weather update server.
//  Binds PUB socket to tcp://*:5556
//  Publishes random weather updates
//
package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq3"
	"math/rand"
	"time"
)

func main() {

	//  Prepare our context and publisher
	context, _ := zmq.NewContext()
	publisher, _ := context.NewSocket(zmq.PUB)
	publisher.Bind("tcp://*:5556")
	publisher.Bind("ipc://weather.ipc")

	//  Initialize random number generator
	rand.Seed(time.Now().Unix())

	// loop for a while aparently
	for {

		//  Get values that will fool the boss
		zipcode := rand.Intn(100000)
		temperature := rand.Intn(215) - 80
		relhumidity := rand.Intn(50) + 10

		//  Send message to all subscribers
		msg := fmt.Sprintf("%05d %d %d", zipcode, temperature, relhumidity)
		publisher.Send([]byte(msg), 0)
	}
}
