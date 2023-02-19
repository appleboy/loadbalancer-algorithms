package weighted

import (
	"errors"
	"net/url"
	"sync"
)

var (
	ErrServersEmpty   = errors.New("server list is empty")
	ErrServerNotFound = errors.New("server not found")
)

type server struct {
	url    *url.URL
	weight int
}

type RoundRobin interface {
	NextServer() *url.URL
	AddServer(*url.URL, int) error
	RemoveServer(*url.URL) error
	Servers() []*url.URL
	RemoveAll()
	// Reset resets all current weights.
	Reset()
}

// Weighted Round Robin
// http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling
type roundrobin struct {
	sync.Mutex
	servers []*server

	// index indicates the server selected last time, and i is initialized with -1
	index int
	// cw is the current weight in scheduling, and cw is initialized with zero.
	cw int
	// maxWeigt is the maximum weight of all the servers.
	maxWeigt int
	// gcd is the greatest common divisor of all server weights.
	gcd int
	// current server list count
	count int
}

//	while (true) {
//	  i = (i + 1) mod n;
//	  if (i == 0) {
//	      cw = cw - gcd(S);
//	      if (cw <= 0) {
//	          cw = max(S);
//	          if (cw == 0)
//	          return NULL;
//	      }
//	  }
//	  if (W(Si) >= cw)
//	      return Si;
//	}
//
// reference: http://kb.linuxvirtualserver.org/wiki/Weighted_Round-Robin_Scheduling
func (r *roundrobin) NextServer() *url.URL {
	if r.count == 0 {
		return nil
	}

	if r.count == 1 {
		return r.servers[0].url
	}

	for {
		r.index = (r.index + 1) % r.count
		if r.index == 0 {
			r.cw = r.cw - r.gcd
			if r.cw <= 0 {
				r.cw = r.maxWeigt
				if r.cw == 0 {
					return nil
				}
			}
		}

		if r.servers[r.index].weight >= r.cw {
			return r.servers[r.index].url
		}
	}
}

func (r *roundrobin) AddServer(url *url.URL, weight int) error {
	if weight > 0 {
		if r.gcd == 0 {
			r.gcd = weight
			r.maxWeigt = weight
			r.index = -1
			r.cw = 0
		} else {
			r.gcd = gcd(r.gcd, weight)
			if r.maxWeigt < weight {
				r.maxWeigt = weight
			}
		}
	}

	r.servers = append(r.servers, &server{
		url:    url,
		weight: weight,
	})
	r.count += 1
	return nil
}

func (r *roundrobin) RemoveServer(url *url.URL) error {
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

func (r *roundrobin) Servers() []*url.URL {
	r.Lock()
	urls := make([]*url.URL, len(r.servers))
	for i, s := range r.servers {
		urls[i] = s.url
	}
	r.Unlock()

	return urls
}

func (r *roundrobin) RemoveAll() {
	r.servers = r.servers[:0]
	r.count = 0
	r.cw = 0
	r.index = -1
	r.gcd = 0
	r.maxWeigt = 0
}

// Reset resets all current weights.
func (r *roundrobin) Reset() {
	r.index = -1
	r.cw = 0
}

func New() (RoundRobin, error) {
	rb := &roundrobin{
		servers: []*server{},
		count:   0,
		index:   -1,
		cw:      0,
	}

	return rb, nil
}

func checkURL(a, b *url.URL) bool {
	return a.Path == b.Path && a.Host == b.Host && a.Scheme == b.Scheme
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
