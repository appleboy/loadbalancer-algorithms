package round_robin

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
		if r.Next().Host != s.Host {
			t.Fatalf("Expected %s, but got %s", s.Host, r.Next().Host)
		}
	}
}

func TestAdd(t *testing.T) {
	r, _ := New(servers...)
	r.Add(addServers...)
	newServers := append(servers, addServers...)

	for _, s := range newServers {
		if r.Next().Host != s.Host {
			t.Fatalf("Expected %s, but got %s", s.Host, r.Next().Host)
		}
	}
}

func ExampleRoundRobin() {
	r, _ := New(servers...)

	fmt.Println(r.Next().Host)
	fmt.Println(r.Next().Host)
	fmt.Println(r.Next().Host)
	fmt.Println(r.Next().Host)
	fmt.Println()

	r.Add(addServers...)

	fmt.Println(r.Next().Host)
	fmt.Println(r.Next().Host)
	fmt.Println(r.Next().Host)
	fmt.Println(r.Next().Host)
	fmt.Println(r.Next().Host)

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
