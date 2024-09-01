package roundrobin

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/appleboy/loadbalancer-algorithms/proxy"
)

func TestNextServer(t *testing.T) {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
	}

	r, _ := New(servers...)

	for i := 0; i < len(servers); i++ {
		nextServer := r.NextServer()
		expectedServer := servers[i]
		if nextServer != expectedServer {
			t.Fatalf("Expected server %s, but got %s", expectedServer.GetName(), nextServer.GetName())
		}
	}
}

func TestAddServers(t *testing.T) {
	r, _ := New()

	server1 := proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"})
	server2 := proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"})

	err := r.AddServers(server1, server2)
	if err != nil {
		t.Fatalf("Failed to add servers: %v", err)
	}

	servers := r.Servers()
	if len(servers) != 2 {
		t.Fatalf("Expected 2 servers, but got %d", len(servers))
	}

	if servers[0] != server1 {
		t.Fatalf("Expected server1, but got %v", servers[0])
	}

	if servers[1] != server2 {
		t.Fatalf("Expected server2, but got %v", servers[1])
	}

	for i := 0; i < len(servers); i++ {
		nextServer := r.NextServer()
		expectedServer := servers[i]
		if nextServer != expectedServer {
			t.Fatalf("Expected server %s, but got %s", expectedServer.GetName(), nextServer.GetName())
		}
	}
}

func BenchmarkNext(b *testing.B) {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.13"}),
	}
	r, _ := New(servers...)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.NextServer()
	}
}

func ExampleRoundRobin() {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
	}
	r, _ := New(servers...)

	fmt.Println(r.NextServer().GetName())
	fmt.Println(r.NextServer().GetName())
	fmt.Println(r.NextServer().GetName())
	fmt.Println(r.NextServer().GetName())
	fmt.Println()

	addServers := []*proxy.Proxy{
		proxy.NewProxy("d1", &url.URL{Host: "192.168.2.10"}),
		proxy.NewProxy("d2", &url.URL{Host: "192.168.2.11"}),
		proxy.NewProxy("d3", &url.URL{Host: "192.168.2.12"}),
	}

	_ = r.AddServers(addServers...)

	fmt.Println(r.NextServer().GetName())
	fmt.Println(r.NextServer().GetName())
	fmt.Println(r.NextServer().GetName())
	fmt.Println(r.NextServer().GetName())
	fmt.Println(r.NextServer().GetName())
	fmt.Println(r.NextServer().GetName())

	// Output:
	// s1
	// s2
	// s3
	// s1
	//
	// d2
	// d3
	// s1
	// s2
	// s3
	// d1
}
