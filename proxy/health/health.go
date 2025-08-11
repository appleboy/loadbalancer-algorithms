package health

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	// Default to 10 seconds. The minimum value is 1
	defaultPeriod = 10
	// If the value of period is greater than initialDelay then the initialDelay will be ignored
	// Defaults to 0 seconds. Minimum value is 0.
	defaultInitialDelay = 0
	// Default success thresholds
	defaultSuccessThreshold = 1
	// Default failure threshold
	defaultFailureThreshold = 3
)

type Check func(addr *url.URL) error

func New(origin *url.URL, opts ...Opts) *ProxyHealth {
	h := &ProxyHealth{
		origin:              origin,
		check:               defaultHTTPCheck,
		periodSeconds:       defaultPeriod,
		initialDelaySeconds: defaultInitialDelay,
		successThreshold:    defaultSuccessThreshold,
		failureThreshold:    defaultFailureThreshold,
		successCount:        0,
		failureCount:        0,
		cancel:              make(chan struct{}),
	}

	for _, opt := range opts {
		opt(h)
	}

	h.run()
	return h
}

type ProxyHealth struct {
	origin *url.URL

	mu                  sync.Mutex
	check               Check
	periodSeconds       int
	initialDelaySeconds int
	successThreshold    int
	failureThreshold    int
	successCount        int
	failureCount        int
	cancel              chan struct{}
	isAvailable         bool
	errors              error
}

// checkHealth checks the health of the proxy origin.
func (h *ProxyHealth) checkHealth() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	err := h.check(h.origin)

	if err == nil {
		h.successCount++
		h.failureCount = 0
	} else {
		h.successCount = 0
		h.failureCount++
	}

	if h.successCount >= h.successThreshold {
		h.isAvailable = true
		h.successCount = 0
	}

	if h.failureCount >= h.failureThreshold {
		h.isAvailable = false
		h.failureCount = 0
	}

	return err
}

func (h *ProxyHealth) run() {
	// initial delay
	if h.initialDelaySeconds > 0 {
		select {
		case <-time.After(time.Duration(h.initialDelaySeconds) * time.Second):
		case <-h.cancel:
			return
		}
	}

	go func() {
		for {
			select {
			case <-h.cancel:
				return
			default:
			}

			h.errors = h.checkHealth()

			select {
			case <-time.After(time.Duration(h.periodSeconds) * time.Second):
			case <-h.cancel:
				return
			}
		}
	}()
}

// stop stops the currently rinning check func.
func (h *ProxyHealth) stop() {
	if h.cancel == nil {
		return
	}

	close(h.cancel)
	h.cancel = nil
}

// IsAvailable returns whether the proxy origin was successfully connected at the last check time.
func (h *ProxyHealth) IsAvailable() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.isAvailable
}

// defaultHTTPCheck is a default health check function that checks
// if the HTTP connection to the address is successful.
func defaultHTTPCheck(addr *url.URL) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
		// never follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(addr.String())
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log the error or handle it appropriately
			// For now, we'll silently ignore Close errors in defer
			_ = err
		}
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return errors.New("invalid status code")
}

// defaultTCPCheck is a default health check function that checks
// if the TCP connection to the address is successful.
func defaultTCPCheck(addr *url.URL) bool {
	conn, err := net.DialTimeout("tcp", addr.Host, 5*time.Second)
	if err != nil {
		return false
	}
	return conn.Close() == nil
}

// defaultDNSCheck is a default health check function that checks
// if the DNS resolution to the address is successful.
func defaultDNSCheck(addr *url.URL) bool {
	resolver := net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	addrs, err := resolver.LookupHost(ctx, addr.Host)
	if err != nil {
		return false
	}
	if len(addrs) < 1 {
		return false
	}
	return true
}
