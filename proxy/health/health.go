package health

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	// Default to 10 seconds. The minimum value is 1
	defaultPeriod = 10 * time.Second
	// If the value of period is greater than initialDelay then the initialDelay will be ignored
	// Defaults to 0 seconds. Minimum value is 0.
	defaultInitialDelay = 0 * time.Second
)

type Check func(addr *url.URL) bool

func New(origin *url.URL, opts ...Opts) *ProxyHealth {
	h := &ProxyHealth{
		origin:       origin,
		check:        defaultHttpCheck,
		period:       defaultPeriod,
		initialDelay: defaultInitialDelay,
		cancel:       make(chan struct{}),
	}

	for _, opt := range opts {
		opt(h)
	}

	h.run()
	return h
}

type ProxyHealth struct {
	origin *url.URL

	mu           sync.Mutex
	check        Check
	period       time.Duration
	initialDelay time.Duration
	cancel       chan struct{}
	isAvailable  bool
}

func (h *ProxyHealth) run() {
	checkHealth := func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		isAvailable := h.check(h.origin)
		h.isAvailable = isAvailable
	}

	if h.initialDelay > h.period {
		h.initialDelay = 0 * time.Second
	}

	// initial delay
	if h.initialDelay > 0 {
		select {
		case <-time.After(h.initialDelay):
		case <-h.cancel:
			return
		}
	}

	go func() {
		t := time.NewTicker(h.period)
		for {
			select {
			case <-t.C:
				checkHealth()
			case <-h.cancel:
				t.Stop()
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

func defaultHttpCheck(addr *url.URL) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(addr.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
