package proxy

import (
	"net/url"
	"sync"
	"time"
)

type Check func(addr *url.URL) bool

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

	// initial delay
	select {
	case <-time.After(h.initialDelay):
	case <-h.cancel:
		return
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
