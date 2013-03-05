package zmq3

import (
	"fmt"
)

/*
Send multi-part message on socket.

Any `[]string' or `[][]byte' is split into separate `string's or `[]byte's

Any other part that isn't a `string' or `[]byte' is converted
to `string' with `fmt.Sprintf("%v", part)'.

Returns error code of sending the final part.
*/
func (soc *Socket) SendMessage(parts ...interface{}) (err error) {
	pp := make([]interface{}, 0)
	for _, p := range parts {
		switch t := p.(type) {
		case []string:
			for _, s := range t {
				pp = append(pp, s)
			}
		case [][]byte:
			for _, b := range t {
				pp = append(pp, b)
			}
		default:
			pp = append(pp, t)
		}
	}

	n := len(pp)
	if n == 0 {
		return
	}
	opt := SNDMORE
	for i, p := range pp {
		if i == n-1 {
			opt = 0
		}
		switch t := p.(type) {
		case string:
			_, err = soc.Send(t, opt)
		case []byte:
			_, err = soc.SendBytes(t, opt)
		default:
			_, err = soc.Send(fmt.Sprintf("%v", t), opt)
		}
	}
	return // error code of last call
}

/*
Receive parts as message from socket.

Returns error code of receiving the final part.
*/
func (soc *Socket) RecvMessage(flags Flag) (msg []string, err error) {
	msg = make([]string, 0)
	for {
		s, err := soc.Recv(flags)
		if err == nil {
			msg = append(msg, s)
		}
		more, _ := soc.GetRcvmore()
		if !more {
			break
		}
	}
	return // error code of last call
}

/*
Receive parts as message from socket.

Returns error code of receiving the final part.
*/
func (soc *Socket) RecvMessageBytes(flags Flag) (msg [][]byte, err error) {
	msg = make([][]byte, 0)
	for {
		b, err := soc.RecvBytes(flags)
		if err == nil {
			msg = append(msg, b)
		}
		more, _ := soc.GetRcvmore()
		if !more {
			break
		}
	}
	return // error code of last call
}
