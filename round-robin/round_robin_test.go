package round_robin

import (
	"net/url"
	"testing"
)

var servers = []*url.URL{
	{Host: "192.168.1.10"},
	{Host: "192.168.1.11"},
	{Host: "192.168.1.12"},
	{Host: "192.168.1.13"},
}

func TestNext(t *testing.T) {
	r, _ := New(servers...)

	for _, s := range servers {
		if r.Next().Host != s.Host {
			t.Fatalf("Expected %s, but got %s", s.Host, r.Next().Host)
		}
	}
}
