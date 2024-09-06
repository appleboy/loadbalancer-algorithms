package roundrobin

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/appleboy/loadbalancer-algorithms/proxy"
)

// ErrServersEmpty is returned when the server list is empty.
var ErrServersEmpty = errors.New("server list is empty")

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

// Ensure that roundrobin implements the RoundRobin interface.
var _ RoundRobin = (*roundrobin)(nil)

// roundrobin represents a round-robin load balancing algorithm.
type roundrobin struct {
	sync.Mutex
	servers []*proxy.Proxy
	next    uint32
	count   int
}

// NextServer returns the next server in the round-robin algorithm.
// If there are no servers available, it returns nil.
// The server selection is based on an atomic counter that increments with each call to NextServer.
// The selected server is determined by calculating the index using the modulo operation on the count of servers.
// The index is incremented atomically to ensure thread safety.
// This function acquires a lock before accessing the server list to prevent concurrent modifications.
// It is the responsibility of the caller to release the lock after using the returned server.
func (r *roundrobin) NextServer() *proxy.Proxy {
	index := atomic.AddUint32(&r.next, 1)
	r.Lock()
	server := r.servers[int(index-1)%r.count]
	r.Unlock()
	return server
}

// AddServers adds the given servers to the roundrobin load balancer.
// It takes a variadic parameter of type *proxy.Proxy, representing the servers to be added.
// If no servers are provided, it returns an error of type ErrServersEmpty.
// The function acquires a lock to ensure thread safety while modifying the server list.
// After adding the servers, it updates the count of servers in the load balancer.
// Finally, it releases the lock and returns nil.
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

// RemoveServers removes the specified servers from the roundrobin load balancer.
// It takes a variadic parameter 'names' which represents the names of the servers to be removed.
// If the 'names' parameter is empty, it returns an error of type 'ErrServersEmpty'.
// The function iterates over the 'names' and checks if each server name matches with any of the existing servers.
// If a match is found, the server is removed from the 'servers' slice by using the 'append' function.
// The 'count' field is updated to reflect the new number of servers after removal.
// The function is thread-safe and uses a lock to prevent concurrent access to the 'servers' slice.
// It returns nil if the removal operation is successful.
func (r *roundrobin) RemoveServers(names ...string) error {
	if len(names) == 0 {
		return ErrServersEmpty
	}
	r.Lock()
	defer r.Unlock()
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
	return nil
}

// Servers returns a list of all servers in the roundrobin load balancer.
func (r *roundrobin) Servers() []*proxy.Proxy {
	r.Lock()
	defer r.Unlock()
	return r.servers
}

// RemoveAll removes all servers from the roundrobin load balancer.
func (r *roundrobin) RemoveAll() {
	r.Lock()
	r.servers = r.servers[:0]
	r.count = 0
	r.Unlock()
	atomic.StoreUint32(&r.next, 0)
}

// New creates a new instance of the round-robin load balancer with the specified servers.
func New(servers ...*proxy.Proxy) (RoundRobin, error) {
	rb := &roundrobin{
		servers: servers,
		count:   len(servers),
	}

	return rb, nil
}
