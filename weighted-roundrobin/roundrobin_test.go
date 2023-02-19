package roundrobin

import (
	"net/url"
	"testing"
)

var servers = []*server{
	{
		url: &url.URL{
			Host: "192.168.1.10",
		},
		weight: 4,
	},
	{
		url: &url.URL{
			Host: "192.168.1.11",
		},
		weight: 3,
	},
	{
		url: &url.URL{
			Host: "192.168.1.12",
		},
		weight: 2,
	},
}

func BenchmarkNext(b *testing.B) {
	r, _ := New()
	for _, server := range servers {
		r.AddServer(server.url, server.weight)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.NextServer()
	}
}
