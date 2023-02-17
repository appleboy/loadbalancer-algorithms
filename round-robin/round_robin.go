package round_robin

import (
	"errors"
	"net/url"
)

var ErrServersEmpty = errors.New("server list is empty")

type RoundRobin interface {
	Next() *url.URL
	Add(*url.URL) error
	Remove(*url.URL) error
}

type roundRobin struct {
	urls  []*url.URL
	next  uint32
	count int
}

func (r *roundRobin) Next() *url.URL {
	return nil
}

func (r *roundRobin) Add(u *url.URL) error {
	return nil
}

func (r *roundRobin) Remove(u *url.URL) error {
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
