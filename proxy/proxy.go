package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/appleboy/loadbalancer-algorithms/proxy/health"
)

// Example of a configuration file:
// server s1 example01.com
// server s2 example02.com
// server s3 example03.com

var value int32 = -1

// NewProxy creates a new instance of Proxy with the specified address.
func NewProxy(name string, addr *url.URL) *Proxy {
	return &Proxy{
		name:   name,
		proxy:  httputil.NewSingleHostReverseProxy(addr),
		health: health.NewProxyHealth(addr),
	}
}

// Proxy represents a reverse proxy for load balancing algorithms.
type Proxy struct {
	name    string
	proxy   *httputil.ReverseProxy
	loading uint32
	health  *health.ProxyHealth
}

// ServeHTTP handles the incoming HTTP request and forwards it to the underlying proxy server.
// It increments the load counter by 1 before forwarding the request and decrements it by the given value after the request is processed.
// This method is part of the Proxy struct and implements the http.Handler interface.
//
// Parameters:
// - w: The http.ResponseWriter used to write the response back to the client.
// - r: The http.Request representing the incoming request.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint32(&p.loading, 1)
	defer atomic.AddUint32(&p.loading, uint32(value))
	p.proxy.ServeHTTP(w, r)
}

// GetLoading returns the current loading of the proxy.
func (p *Proxy) GetLoading() uint32 {
	return atomic.LoadUint32(&p.loading)
}

// GetName returns the name of the proxy.
func (p *Proxy) GetName() string {
	return p.name
}

// IsAvailable returns whether the proxy origin was successfully connected at the last check time.
func (p *Proxy) IsAvailable() bool {
	return p.health.IsAvailable()
}
