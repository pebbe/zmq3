// +build windows

package zmq3

/*
#include <zmq.h>
*/
import "C"

/*
ZMQ_FD: Retrieve file descriptor associated with the socket

See: http://api.zeromq.org/3-2:zmq-getsockopt#toc23
*/


/*
func (soc *Socket) GetFd() (uintptr, error) {
	if !soc.opened {
		return uintptr(0), errSocClosed
	}
	value := C.SOCKET(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, C.ZMQ_FD, &value, &size); i != 0 {
		return uintptr(0), errget(err)
	}
	return uintptr(value), nil
}
*/
