package zmq3_test

import (
	zmq "github.com/pebbe/zmq3"

	"errors"
	"fmt"
	"testing"
)

var (
	err32 = errors.New("rc != 32")
)

func TestVersion(t *testing.T) {
	major, _, _ := zmq.Version()
	if major != 3 {
		t.Errorf("Expected major version 3, got %d", major)
	}
}

func TestConnectResolve(t *testing.T) {

	sock, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		t.Fatal("NewSocket:", err)
	}
	defer func() {
		if sock != nil {
			sock.SetLinger(0)
			sock.Close()
		}
	}()

	err = sock.Connect("tcp://localhost:1234")
	if err != nil {
		t.Error("sock.Connect:", err)
	}

	fails := []string{
		"tcp://localhost:invalid",
		"tcp://in val id:1234",
		"invalid://localhost:1234",
	}
	for _, fail := range fails {
		if err = sock.Connect(fail); err == nil {
			t.Errorf("Connect %s, expected fail, got success", fail)
		}
	}

	err = sock.Close()
	sock = nil
	if err != nil {
		t.Error("sock.Close:", err)
	}
}

func bounce(server, client *zmq.Socket) (msg string, err error) {

	content := "12345678ABCDEFGH12345678abcdefgh"

	//  Send message from client to server
	rc, err := client.Send(content, zmq.SNDMORE|zmq.DONTWAIT)
	if err != nil {
		return "client.Send SNDMORE|DONTWAIT:", err
	}
	if rc != 32 {
		return "client.Send SNDMORE|DONTWAIT:", err32
	}

	rc, err = client.Send(content, zmq.DONTWAIT)
	if err != nil {
		return "client.Send DONTWAIT:", err
	}
	if rc != 32 {
		return "client.Send DONTWAIT:", err32
	}

	//  Receive message at server side
	msg, err = server.Recv(0)
	if err != nil {
		return "server.Recv 1:", err
	}

	//  Check that message is still the same
	if msg != content {
		return "server.Recv 1:", errors.New(fmt.Sprintf("%q != %q", msg, content))
	}

	rcvmore, err := server.GetRcvmore()
	if err != nil {
		return "server.GetRcvmore 1:", err
	}
	if !rcvmore {
		return "server.GetRcvmore 1:", errors.New(fmt.Sprint("rcvmore ==", rcvmore))
	}

	//  Receive message at server side
	msg, err = server.Recv(0)
	if err != nil {
		return "server.Recv 2:", err
	}

	//  Check that message is still the same
	if msg != content {
		return "server.Recv 2:", errors.New(fmt.Sprintf("%q != %q", msg, content))
	}

	rcvmore, err = server.GetRcvmore()
	if err != nil {
		return "server.GetRcvmore 2:", err
	}
	if rcvmore {
		return "server.GetRcvmore 2:", errors.New(fmt.Sprint("rcvmore == ", rcvmore))
	}

	// The same, from server back to client

	//  Send message from server to client
	rc, err = server.Send(content, zmq.SNDMORE)
	if err != nil {
		return "server.Send SNDMORE:", err
	}
	if rc != 32 {
		return "server.Send SNDMORE:", err32
	}

	rc, err = server.Send(content, 0)
	if err != nil {
		return "server.Send 0:", err
	}
	if rc != 32 {
		return "server.Send 0:", err32
	}

	//  Receive message at client side
	msg, err = client.Recv(0)
	if err != nil {
		return "client.Recv 1:", err
	}

	//  Check that message is still the same
	if msg != content {
		return "client.Recv 1:", errors.New(fmt.Sprintf("%q != %q", msg, content))
	}

	rcvmore, err = client.GetRcvmore()
	if err != nil {
		return "client.GetRcvmore 1:", err
	}
	if !rcvmore {
		return "client.GetRcvmore 1:", errors.New(fmt.Sprint("rcvmore ==", rcvmore))
	}

	//  Receive message at client side
	msg, err = client.Recv(0)
	if err != nil {
		return "client.Recv 2:", err
	}

	//  Check that message is still the same
	if msg != content {
		return "client.Recv 2:", errors.New(fmt.Sprintf("%q != %q", msg, content))
	}

	rcvmore, err = client.GetRcvmore()
	if err != nil {
		return "client.GetRcvmore 2:", err
	}
	if rcvmore {
		return "client.GetRcvmore 2:", errors.New(fmt.Sprint("rcvmore == ", rcvmore))
	}
	return "OK", nil
}
