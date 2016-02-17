package zmq3_test

import (
	zmq "github.com/pebbe/zmq3"

	"errors"
	"fmt"
	"testing"
	"time"
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

func TestPoller(t *testing.T) {

	var sb, sc *zmq.Socket

	defer func() {
		for _, s := range []*zmq.Socket{sb, sc} {
			if s != nil {
				s.SetLinger(0)
				s.Close()
			}
		}
	}()

	sb, err := zmq.NewSocket(zmq.PAIR)
	if err != nil {
		t.Fatal("NewSocket:", err)
	}

	err = sb.Bind("tcp://127.0.0.1:9737")
	if err != nil {
		t.Fatal("sb.Bind:", err)
	}

	sc, err = zmq.NewSocket(zmq.PAIR)
	if err != nil {
		t.Fatal("NewSocket:", err)
	}

	err = sc.Connect("tcp://127.0.0.1:9737")
	if err != nil {
		t.Fatal("sc.Connect:", err)
	}

	poller := zmq.NewPoller()
	idxb := poller.Add(sb, 0)
	idxc := poller.Add(sc, 0)
	if idxb != 0 || idxc != 1 {
		t.Errorf("idxb=%d idxc=%d", idxb, idxc)
	}

	if pa, err := poller.PollAll(100 * time.Millisecond); err != nil {
		t.Error("PollAll 1:", err)
	} else if len(pa) != 2 {
		t.Errorf("PollAll 1 len = %d", len(pa))
	} else if pa[0].Events != 0 || pa[1].Events != 0 {
		t.Errorf("PollAll 1 events = %v, %v", pa[0], pa[1])
	}

	poller.Update(idxb, zmq.POLLOUT)
	poller.UpdateBySocket(sc, zmq.POLLIN)

	if pa, err := poller.PollAll(100 * time.Millisecond); err != nil {
		t.Error("PollAll 2:", err)
	} else if len(pa) != 2 {
		t.Errorf("PollAll 2 len = %d", len(pa))
	} else if pa[0].Events != zmq.POLLOUT || pa[1].Events != 0 {
		t.Errorf("PollAll 2 events = %v, %v", pa[0], pa[1])
	}

	poller.UpdateBySocket(sb, 0)

	content := "12345678ABCDEFGH12345678ABCDEFGH"

	//  Send message from client to server
	if rc, err := sb.Send(content, zmq.DONTWAIT); err != nil {
		t.Error("sb.Send DONTWAIT:", err)
	} else if rc != 32 {
		t.Error("sb.Send DONTWAIT:", err32)
	}

	if pa, err := poller.PollAll(100 * time.Millisecond); err != nil {
		t.Error("PollAll 3:", err)
	} else if len(pa) != 2 {
		t.Errorf("PollAll 3 len = %d", len(pa))
	} else if pa[0].Events != 0 || pa[1].Events != zmq.POLLIN {
		t.Errorf("PollAll 3 events = %v, %v", pa[0], pa[1])
	}

	//  Receive message
	if msg, err := sc.Recv(zmq.DONTWAIT); err != nil {
		t.Error("sb.Recv DONTWAIT:", err)
	} else if msg != content {
		t.Error("sb.Recv msg != content")
	}

	poller.UpdateBySocket(sb, zmq.POLLOUT)
	poller.Update(idxc, zmq.POLLIN)

	if pa, err := poller.PollAll(100 * time.Millisecond); err != nil {
		t.Error("PollAll 4:", err)
	} else if len(pa) != 2 {
		t.Errorf("PollAll 4 len = %d", len(pa))
	} else if pa[0].Events != zmq.POLLOUT || pa[1].Events != 0 {
		t.Errorf("PollAll 4 events = %v, %v", pa[0], pa[1])
	}

	err = sc.Close()
	sc = nil
	if err != nil {
		t.Error("sc.Close:", err)
	}

	err = sb.Close()
	sb = nil
	if err != nil {
		t.Error("sb.Close:", err)
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
