package round_robin

import (
	"errors"
	"net/url"
	"sync"
	"sync/atomic"
)

var ErrServersEmpty = errors.New("server list is empty")

type RoundRobin interface {
	Next() *url.URL
	Add(...*url.URL) error
	Remove(...*url.URL) error
}

type roundRobin struct {
	sync.Mutex
	urls  []*url.URL
	next  uint32
	count int
}

func (r *roundRobin) Next() *url.URL {
	index := atomic.AddUint32(&r.next, 1)
	return r.urls[int(index-1)%r.count]
}

func (r *roundRobin) Add(urls ...*url.URL) error {
	if len(urls) == 0 {
		return ErrServersEmpty
	}
	r.Lock()
	r.urls = append(r.urls, urls...)
	r.count = len(r.urls)
	r.Unlock()
	return nil
}

func (r *roundRobin) Remove(urls ...*url.URL) error {
	return nil
}

func New(urls ...*url.URL) (RoundRobin, error) {
	if len(urls) == 0 {
		return nil, ErrServersEmpty
	}

	return &roundRobin{
		urls:  urls,
		count: len(urls),
	}, nil
}