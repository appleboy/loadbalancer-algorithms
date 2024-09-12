package health

type Opts func(*ProxyHealth)

// WithCheck sets the health check function for the proxy.
func WithCheck(check Check) Opts {
	return func(h *ProxyHealth) {
		h.check = check
	}
}

// WithPeriodSeconds sets the period for the health check.
func WithPeriodSeconds(periodSeconds int) Opts {
	return func(h *ProxyHealth) {
		h.periodSeconds = periodSeconds
	}
}

// WithInitialDelaySeconds sets the initial delay for the health check.
func WithInitialDelaySeconds(initialDelaySeconds int) Opts {
	return func(h *ProxyHealth) {
		h.initialDelaySeconds = initialDelaySeconds
	}
}
