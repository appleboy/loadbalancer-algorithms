package roundrobin

import (
	"fmt"
	"net/url"
	"testing"
)

var servers = []*url.URL{
	{Host: "192.168.1.10"},
	{Host: "192.168.1.11"},
	{Host: "192.168.1.12"},
	{Host: "192.168.1.13"},
}

var addServers = []*url.URL{
	{Host: "192.168.2.10"},
	{Host: "192.168.2.11"},
	{Host: "192.168.2.12"},
	{Host: "192.168.2.13"},
}

func TestNext(t *testing.T) {
	r, _ := New(servers...)

	for _, s := range servers {
		if r.NextServer().Host != s.Host {
			t.Fatalf("Expected %s, but got %s", s.Host, r.NextServer().Host)
		}
	}
}

func TestAdd(t *testing.T) {
	r, _ := New(servers...)
	_ = r.AddServers(addServers...)
	newServers := append(servers, addServers...)

	for _, s := range newServers {
		if r.NextServer().Host != s.Host {
			t.Fatalf("Expected %s, but got %s", s.Host, r.NextServer().Host)
		}
	}
}

func TestRemove(t *testing.T) {
	expectServers := servers[:len(servers)-1]

	r, _ := New(servers...)
	_ = r.RemoveServer(&url.URL{
		Host: "192.168.1.13",
	})

	for _, s := range expectServers {
		if r.NextServer().Host != s.Host {
			t.Fatalf("Expected %s, but got %s", s.Host, r.NextServer().Host)
		}
	}
}

func TestGetServers(t *testing.T) {
	r, _ := New(servers...)
	_ = r.RemoveServer(&url.URL{
		Host: "192.168.1.13",
	})

	for _, s := range r.Servers() {
		if r.NextServer().Host != s.Host {
			t.Fatalf("Expected %s, but got %s", s.Host, r.NextServer().Host)
		}
	}
}

func TestRemoveAll(t *testing.T) {
	r, _ := New(servers...)
	r.RemoveAll()

	if len(r.Servers()) != 0 {
		t.Fatalf("Expected zero, but got %d", len(r.Servers()))
	}
}

func ExampleRoundRobin() {
	r, _ := New(servers...)

	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println()

	_ = r.AddServers(addServers...)

	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)
	fmt.Println(r.NextServer().Host)

	// Output:
	// 192.168.1.10
	// 192.168.1.11
	// 192.168.1.12
	// 192.168.1.13
	//
	// 192.168.2.10
	// 192.168.2.11
	// 192.168.2.12
	// 192.168.2.13
	// 192.168.1.10
}

func BenchmarkNext(b *testing.B) {
	r, _ := New(servers...)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.NextServer()
	}
}
