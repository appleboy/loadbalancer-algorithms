package roundrobin

import (
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
