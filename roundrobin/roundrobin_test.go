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

func TestRemoveAll(t *testing.T) {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
	}

	r, _ := New(servers...)

	r.RemoveAll()

	if len(r.Servers()) != 0 {
		t.Fatalf("Expected 0 servers after RemoveAll, but got %d", len(r.Servers()))
	}
}

func TestServers(t *testing.T) {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
	}

	r, _ := New(servers...)

	result := r.Servers()

	if len(result) != len(servers) {
		t.Fatalf("Expected %d servers, but got %d", len(servers), len(result))
	}

	for i := 0; i < len(servers); i++ {
		if result[i] != servers[i] {
			t.Fatalf("Expected server %s, but got %s", servers[i].GetName(), result[i].GetName())
		}
	}
}

func TestRemoveServers(t *testing.T) {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
	}

	r, _ := New(servers...)

	err := r.RemoveServers("s1", "s3")
	if err != nil {
		t.Fatalf("Failed to remove servers: %v", err)
	}

	remainingServers := r.Servers()
	if len(remainingServers) != 1 {
		t.Fatalf("Expected 1 server after removal, but got %d", len(remainingServers))
	}

	if remainingServers[0].GetName() != "s2" {
		t.Fatalf("Expected server s2, but got %v", remainingServers[0])
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

func BenchmarkNextParallel(b *testing.B) {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
		proxy.NewProxy("s4", &url.URL{Host: "192.168.1.13"}),
	}
	r, _ := New(servers...)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.NextServer()
		}
	})
}

func BenchmarkAddServers(b *testing.B) {
	r, _ := New()
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.RemoveAll() // Reset for each iteration
		r.AddServers(servers...)
	}
}

func BenchmarkRemoveServers(b *testing.B) {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
		proxy.NewProxy("s4", &url.URL{Host: "192.168.1.13"}),
		proxy.NewProxy("s5", &url.URL{Host: "192.168.1.14"}),
	}

	r, _ := New(servers...)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset servers for each iteration
		r.RemoveAll()
		r.AddServers(servers...)
		r.RemoveServers("s1", "s3")
	}
}

func BenchmarkServers(b *testing.B) {
	servers := []*proxy.Proxy{
		proxy.NewProxy("s1", &url.URL{Host: "192.168.1.10"}),
		proxy.NewProxy("s2", &url.URL{Host: "192.168.1.11"}),
		proxy.NewProxy("s3", &url.URL{Host: "192.168.1.12"}),
		proxy.NewProxy("s4", &url.URL{Host: "192.168.1.13"}),
	}
	r, _ := New(servers...)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Servers()
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
