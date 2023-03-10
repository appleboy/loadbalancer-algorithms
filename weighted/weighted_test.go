package weighted

import (
	"errors"
	"fmt"
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

func TestRemoveServer(t *testing.T) {
	r, _ := New()
	for _, server := range servers {
		_ = r.AddServer(server.url, server.weight)
	}

	err := r.RemoveServer(&url.URL{
		Host: "192.168.1.100",
	})

	if !errors.Is(ErrServerNotFound, err) {
		t.Fatal(err)
	}

	_ = r.RemoveServer(&url.URL{
		Host: "192.168.1.10",
	})

	servs := r.Servers()
	if len(servs) != len(servers)-1 {
		t.Fatal("can't remove server")
	}
}

func ExampleRoundRobin() {
	r, _ := New()
	for _, server := range servers {
		_ = r.AddServer(server.url, server.weight)
	}

	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)

	// Output:
	// 192.168.1.10
	// 192.168.1.10
	// 192.168.1.11
	// 192.168.1.10
	// 192.168.1.11
	// 192.168.1.12
	// 192.168.1.10
	// 192.168.1.11
	// 192.168.1.12
}

func BenchmarkNext(b *testing.B) {
	r, _ := New()
	for _, server := range servers {
		_ = r.AddServer(server.url, server.weight)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.NextServer()
	}
}
