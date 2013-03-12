#ifndef _ZMQ2_H_
#define _ZMQ2_H_ 1

int zmq2_send (void *socket, zmq_msg_t *msg, int flags);
int zmq3_send (void *socket, void *buf, size_t len, int flags);
int zmq2_recv (void *socket, zmq_msg_t *msg, int flags);

#if ZMQ_VERSION_MAJOR < 3

#define ZMQ_DELAY_ATTACH_ON_CONNECT 0
#define ZMQ_DONTWAIT ZMQ_NOBLOCK
#define ZMQ_EVENT_ACCEPTED 0
#define ZMQ_EVENT_ACCEPT_FAILED 0
#define ZMQ_EVENT_ALL 0
#define ZMQ_EVENT_BIND_FAILED 0
#define ZMQ_EVENT_CLOSED 0
#define ZMQ_EVENT_CLOSE_FAILED 0
#define ZMQ_EVENT_CONNECTED 0
#define ZMQ_EVENT_CONNECT_DELAYED 0
#define ZMQ_EVENT_CONNECT_RETRIED 0
#define ZMQ_EVENT_DISCONNECTED 0
#define ZMQ_EVENT_LISTENING 0
#define ZMQ_IO_THREADS 0
#define ZMQ_IPV4ONLY 0
#define ZMQ_LAST_ENDPOINT 0
#define ZMQ_MAXMSGSIZE 0
#define ZMQ_MAX_SOCKETS 0
#define ZMQ_MULTICAST_HOPS 0
#define ZMQ_RCVHWM ZMQ_HWM
#define ZMQ_ROUTER_MANDATORY 0
#define ZMQ_SNDHWM ZMQ_HWM
#define ZMQ_TCP_ACCEPT_FILTER 0
#define ZMQ_TCP_KEEPALIVE 0
#define ZMQ_TCP_KEEPALIVE_CNT 0
#define ZMQ_TCP_KEEPALIVE_IDLE 0
#define ZMQ_TCP_KEEPALIVE_INTVL 0
#define ZMQ_XPUB_VERBOSE 0
int zmq_ctx_get (void *context, int option_name);
int zmq_ctx_set (void *context, int option_name, int option_value);
int zmq_disconnect (void *socket, const char *endpoint);
int zmq_msg_recv (zmq_msg_t *msg, void *socket, int flags);
int zmq_proxy (const void *frontend, const void *backend, const void *capture);
int zmq_socket_monitor (void *s, const char *addr, int events);
int zmq_unbind (void *socket, const char *endpoint);
void *zmq_ctx_new ();
typedef struct {
    int event;
    union {
    struct {
        char *addr;
        int fd;
    } connected;
    } data;
} zmq_event_t;

#else // ZMQ_VERSION_MAJOR >= 3

#define ZMQ_RECOVERY_IVL_MSEC 0

#endif // ZMQ_VERSION_MAJOR

#endif // _ZMQ2_H_
