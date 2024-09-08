package health

import "time"

type Opts func(*ProxyHealth)

// WithCheck sets the health check function for the proxy.
func WithCheck(check Check) Opts {
	return func(h *ProxyHealth) {
		h.check = check
	}
}

// WithPeriod sets the period for the health check.
func WithPeriod(period time.Duration) Opts {
	return func(h *ProxyHealth) {
		h.period = period
	}
}

// WithInitialDelay sets the initial delay for the health check.
func WithInitialDelay(initialDelay time.Duration) Opts {
	return func(h *ProxyHealth) {
		h.initialDelay = initialDelay
	}
}
