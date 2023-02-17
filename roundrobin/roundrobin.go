package roundrobin

import (
	"errors"
	"net/url"
	"sync"
	"sync/atomic"
)

var (
	ErrServersEmpty   = errors.New("server list is empty")
	ErrServerNotFound = errors.New("server not found")
)

type server struct {
	url *url.URL
}

type RoundRobin interface {
	NextServer() *url.URL
	AddServers(...*url.URL) error
	RemoveServer(*url.URL) error
	Servers() []*url.URL
}

type roundRobin struct {
	sync.Mutex
	servers []*server
	next    uint32
	count   int
}

func (r *roundRobin) NextServer() *url.URL {
	index := atomic.AddUint32(&r.next, 1)
	server := r.servers[int(index-1)%r.count]
	return server.url
}

func (r *roundRobin) AddServers(urls ...*url.URL) error {
	if len(urls) == 0 {
		return ErrServersEmpty
	}
	r.Lock()
	for _, url := range urls {
		r.servers = append(r.servers, &server{url: url})
	}
	r.count = len(r.servers)
	r.Unlock()
	return nil
}

func (r *roundRobin) RemoveServer(url *url.URL) error {
	r.Lock()
	defer r.Unlock()
	for i, s := range r.servers {
		if checkURL(url, s.url) {
			r.servers = append(r.servers[:i], r.servers[i+1:]...)
			return nil
		}
	}
	return ErrServerNotFound
}

func (r *roundRobin) Servers() []*url.URL {
	r.Lock()
	urls := make([]*url.URL, len(r.servers))
	for i, s := range r.servers {
		urls[i] = s.url
	}
	r.Unlock()

	return urls
}

func New(urls ...*url.URL) (RoundRobin, error) {
	if len(urls) == 0 {
		return nil, ErrServersEmpty
	}

	rb := &roundRobin{
		servers: []*server{},
	}

	for _, url := range urls {
		rb.servers = append(rb.servers, &server{url: url})
	}
	rb.count = len(rb.servers)

	return rb, nil
}

func checkURL(a, b *url.URL) bool {
	return a.Path == b.Path && a.Host == b.Host && a.Scheme == b.Scheme
}
