// +build !windows

package zmq3

/*
#cgo pkg-config: libzmq
#include <zmq.h>
*/
import "C"


// Retrieve file descriptor associated with the socket
func (soc *Socket) GetFd() (int, error) {
	return soc.getInt(C.ZMQ_FD)
}

