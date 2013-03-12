package zmq3

/*
#include <zmq.h>
#if ZMQ_VERSION_MAJOR < 3
int zmq_ctx_get (void *context, int option_name) { return -1; }
int zmq_ctx_set (void *context, int option_name, int option_value) { return -1; }
int zmq_disconnect (void *socket, const char *endpoint) { return -1; }
int zmq_msg_recv (zmq_msg_t *msg, void *socket, int flags) { return -1; }
int zmq_proxy (const void *frontend, const void *backend, const void *capture) { return -1; }
int zmq_socket_monitor (void *s, const char *addr, int events)  { return -1; }
int zmq_unbind (void *socket, const char *endpoint) { return -1; }
void *zmq_ctx_new () { return NULL; }

int zmq2_recv (void *socket, zmq_msg_t *msg, int flags) {
    return zmq_recv (socket, msg, flags);
}
int zmq2_send (void *socket, zmq_msg_t *msg, int flags) {
    return zmq_send(socket, msg, flags);
}
int zmq3_send (void *socket, void *buf, size_t len, int flags) { return -1; }

#else // ZMQ_VERSION_MAJOR >= 3

int zmq2_recv (void *socket, zmq_msg_t *msg, int flags) { return -1; }
int zmq2_send (void *socket, zmq_msg_t *msg, int flags) { return -1; }
int zmq3_send (void *socket, void *buf, size_t len, int flags) {
    return zmq_send(socket, buf, len, flags);
}


#endif // ZMQ_VERSION_MAJOR
*/
import "C"
