package roundrobin

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/appleboy/loadbalancer-algorithms/proxy"
)

var (
	ErrServersEmpty   = errors.New("server list is empty")
	ErrServerNotFound = errors.New("server not found")
)

// RoundRobin is an interface that defines the methods for a round-robin load balancer algorithm.
type RoundRobin interface {
	// NextServer returns the next server in the rotation.
	NextServer() *proxy.Proxy

	// AddServers adds one or more servers to the load balancer.
	AddServers(...*proxy.Proxy) error

	// RemoveServers removes one or more servers from the load balancer.
	RemoveServers(...string) error

	// Servers returns a list of all servers in the load balancer.
	Servers() []*proxy.Proxy

	// RemoveAll removes all servers from the load balancer.
	RemoveAll()
}

type roundrobin struct {
	sync.Mutex
	servers []*proxy.Proxy
	next    uint32
	count   int
}

func (r *roundrobin) NextServer() *proxy.Proxy {
	if r.count == 0 {
		return nil
	}
	index := atomic.AddUint32(&r.next, 1)
	server := r.servers[int(index-1)%r.count]
	return server
}

func (r *roundrobin) AddServers(servers ...*proxy.Proxy) error {
	if len(servers) == 0 {
		return ErrServersEmpty
	}
	r.Lock()
	r.servers = append(r.servers, servers...)
	r.count = len(r.servers)
	r.Unlock()
	return nil
}

func (r *roundrobin) RemoveServers(names ...string) error {
	if len(names) == 0 {
		return ErrServersEmpty
	}
	r.Lock()
	for _, name := range names {
		for i, server := range r.servers {
			if server.GetName() != name {
				continue
			}
			r.servers = append(r.servers[:i], r.servers[i+1:]...)
			r.count = len(r.servers)
			break
		}
	}
	r.Unlock()
	return nil
}

func (r *roundrobin) Servers() []*proxy.Proxy {
	return r.servers
}

func (r *roundrobin) RemoveAll() {
	r.servers = r.servers[:0]
	r.count = 0
	atomic.StoreUint32(&r.next, 0)
}

func New(servers ...*proxy.Proxy) (RoundRobin, error) {
	if len(servers) == 0 {
		return nil, ErrServersEmpty
	}

	rb := &roundrobin{
		servers: servers,
		count:   len(servers),
	}

	return rb, nil
}
