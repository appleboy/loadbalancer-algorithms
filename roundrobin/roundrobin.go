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
	sync.RWMutex
	servers []*proxy.Proxy
	next    uint32
}

// NextServer returns the next server in the round-robin algorithm.
// If there are no servers available, it returns nil.
// The server selection is based on an atomic counter that increments with each call.
// The selected server is determined by calculating the index using the modulo operation.
// This method is thread-safe using atomic operations and read locks.
func (r *roundrobin) NextServer() *proxy.Proxy {
	index := atomic.AddUint32(&r.next, 1)

	r.RLock()
	count := uint32(len(r.servers))
	if count == 0 {
		r.RUnlock()
		return nil
	}
	server := r.servers[(index-1)%count]
	r.RUnlock()
	return server
}

// AddServers adds the given servers to the roundrobin load balancer.
// It takes a variadic parameter of type *proxy.Proxy, representing the servers to be added.
// If no servers are provided, it returns an error of type ErrServersEmpty.
// The function uses a write lock to ensure thread safety while modifying the server list.
func (r *roundrobin) AddServers(servers ...*proxy.Proxy) error {
	if len(servers) == 0 {
		return ErrServersEmpty
	}

	r.Lock()
	r.servers = append(r.servers, servers...)
	r.Unlock()
	return nil
}

// RemoveServers removes the specified servers from the roundrobin load balancer.
// It takes a variadic parameter 'names' which represents the names of the servers to be removed.
// If the 'names' parameter is empty, it returns an error of type 'ErrServersEmpty'.
// The function is thread-safe and uses a write lock to prevent concurrent access.
// It uses an optimized in-place removal algorithm to minimize memory allocations.
func (r *roundrobin) RemoveServers(names ...string) error {
	if len(names) == 0 {
		return ErrServersEmpty
	}

	r.Lock()
	defer r.Unlock()

	// For small number of names, use linear search to avoid map allocation
	if len(names) <= 2 {
		for _, name := range names {
			for i := 0; i < len(r.servers); i++ {
				if r.servers[i].GetName() == name {
					// Remove by swapping with last element and truncating
					r.servers[i] = r.servers[len(r.servers)-1]
					r.servers = r.servers[:len(r.servers)-1]
					i-- // Adjust index since we moved an element into current position
					break
				}
			}
		}
		return nil
	}

	// For larger number of names, use map for better performance
	nameMap := make(map[string]struct{}, len(names))
	for _, name := range names {
		nameMap[name] = struct{}{}
	}

	// In-place filtering
	writeIndex := 0
	for readIndex := 0; readIndex < len(r.servers); readIndex++ {
		if _, exists := nameMap[r.servers[readIndex].GetName()]; !exists {
			r.servers[writeIndex] = r.servers[readIndex]
			writeIndex++
		}
	}

	// Truncate the slice and clear references to help GC
	for i := writeIndex; i < len(r.servers); i++ {
		r.servers[i] = nil
	}
	r.servers = r.servers[:writeIndex]
	return nil
}

// Servers returns a copy of all servers in the roundrobin load balancer.
// It returns a new slice to prevent external modifications to the internal state.
func (r *roundrobin) Servers() []*proxy.Proxy {
	r.RLock()
	defer r.RUnlock()

	// Return a copy of the slice to prevent external modifications
	servers := make([]*proxy.Proxy, len(r.servers))
	copy(servers, r.servers)
	return servers
}

// RemoveAll removes all servers from the roundrobin load balancer.
// It resets both the server list and the rotation counter atomically.
func (r *roundrobin) RemoveAll() {
	r.Lock()
	r.servers = r.servers[:0]
	r.Unlock()
	atomic.StoreUint32(&r.next, 0)
}

// New creates a new instance of the round-robin load balancer with the specified servers.
// If no servers are provided, it creates an empty load balancer that can have servers added later.
func New(servers ...*proxy.Proxy) (RoundRobin, error) {
	rb := &roundrobin{
		servers: make([]*proxy.Proxy, len(servers)),
	}

	// Copy servers to prevent external modifications
	copy(rb.servers, servers)

	return rb, nil
}
