package zmq3_test

import (
	zmq "github.com/pebbe/zmq3"

	"errors"
	"fmt"
	"runtime"
)

func Example_test_connect_delay() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_connect_resolve() {

	/*
		copied from zmq4
	*/

	sock, err := zmq.NewSocket(zmq.PUB)
	if checkErr(err) {
		return
	}

	err = sock.Connect("tcp://localhost:1234")
	checkErr(err)

	err = sock.Connect("tcp://localhost:invalid")
	fmt.Println(err)

	err = sock.Connect("tcp://in val id:1234")
	fmt.Println(err)

	err = sock.Connect("invalid://localhost:1234")
	fmt.Println(err)

	err = sock.Close()
	checkErr(err)

	fmt.Println("Done")
	// Output:
	// invalid argument
	// invalid argument
	// protocol not supported
	// Done
}

func Example_test_disconnect_inproc() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_hwm() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_invalid_rep() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_last_endpoint() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_monitor() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_msg_flags() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_pair_inproc() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_pair_ipc() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_pair_tcp() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_reqrep_device() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_reqrep_inproc() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_reqrep_ipc() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_reqrep_tcp() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_router_mandatory() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_shutdown_stress() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_sub_forward() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_term_endpoint() {

	fmt.Println("Done")
	// Output:
	// Done
}

func Example_test_timeo() {

	fmt.Println("Done")
	// Output:
	// Done
}

func bounce(server, client *zmq.Socket) {

	content := "12345678ABCDEFGH12345678abcdefgh"

	//  Send message from client to server
	rc, err := client.Send(content, zmq.SNDMORE)
	if checkErr(err) {
		return
	}
	if rc != 32 {
		checkErr(errors.New("rc != 32"))
	}

	rc, err = client.Send(content, 0)
	if checkErr(err) {
		return
	}
	if rc != 32 {
		checkErr(errors.New("rc != 32"))
	}

	//  Receive message at server side
	msg, err := server.Recv(0)
	if checkErr(err) {
		return
	}

	//  Check that message is still the same
	if msg != content {
		checkErr(errors.New(fmt.Sprintf("%q != %q", msg, content)))
	}

	rcvmore, err := server.GetRcvmore()
	if checkErr(err) {
		return
	}
	if !rcvmore {
		checkErr(errors.New(fmt.Sprint("rcvmore ==", rcvmore)))
		return
	}

	//  Receive message at server side
	msg, err = server.Recv(0)
	if checkErr(err) {
		return
	}

	//  Check that message is still the same
	if msg != content {
		checkErr(errors.New(fmt.Sprintf("%q != %q", msg, content)))
	}

	rcvmore, err = server.GetRcvmore()
	if checkErr(err) {
		return
	}
	if rcvmore {
		checkErr(errors.New(fmt.Sprint("rcvmore == ", rcvmore)))
		return
	}

	// The same, from server back to client

	//  Send message from server to client
	rc, err = server.Send(content, zmq.SNDMORE)
	if checkErr(err) {
		return
	}
	if rc != 32 {
		checkErr(errors.New("rc != 32"))
	}

	rc, err = server.Send(content, 0)
	if checkErr(err) {
		return
	}
	if rc != 32 {
		checkErr(errors.New("rc != 32"))
	}

	//  Receive message at client side
	msg, err = client.Recv(0)
	if checkErr(err) {
		return
	}

	//  Check that message is still the same
	if msg != content {
		checkErr(errors.New(fmt.Sprintf("%q != %q", msg, content)))
	}

	rcvmore, err = client.GetRcvmore()
	if checkErr(err) {
		return
	}
	if !rcvmore {
		checkErr(errors.New(fmt.Sprint("rcvmore ==", rcvmore)))
		return
	}

	//  Receive message at client side
	msg, err = client.Recv(0)
	if checkErr(err) {
		return
	}

	//  Check that message is still the same
	if msg != content {
		checkErr(errors.New(fmt.Sprintf("%q != %q", msg, content)))
	}

	rcvmore, err = client.GetRcvmore()
	if checkErr(err) {
		return
	}
	if rcvmore {
		checkErr(errors.New(fmt.Sprint("rcvmore == ", rcvmore)))
		return
	}

}

func checkErr(err error) bool {
	if err != nil {
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("%v:%v: %v\n", filename, lineno, err)
		} else {
			fmt.Println(err)
		}
		return true
	}
	return false
}
