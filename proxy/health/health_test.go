package health

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestProxyHealth_IsAvailable(t *testing.T) {
	origin, _ := url.Parse("http://example.com")

	// Mock check function
	mockCheck := func(addr *url.URL) bool {
		return addr.String() == "http://example.com"
	}

	h := New(
		origin,
		WithCheck(mockCheck),
		WithPeriodSeconds(1),
	)
	defer h.stop()

	// Wait for the initial check to complete
	time.Sleep(1200 * time.Millisecond)

	if !h.IsAvailable() {
		t.Fatalf("Expected IsAvailable to be true, but got false")
	}

	// Change the origin to make the check fail
	h.origin, _ = url.Parse("http://invalid.com")

	// Wait for the next check to complete
	time.Sleep(1200 * time.Millisecond)

	if h.IsAvailable() {
		t.Fatalf("Expected IsAvailable to be false, but got true")
	}
}

func TestDefaultTCPCheck(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Valid TCP connection",
			url:  "tcp://example.com:80",
			want: true,
		},
		{
			name: "Invalid TCP connection",
			url:  "tcp://invalid:80",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := url.Parse(tt.url)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}

			got := defaultTCPCheck(addr)
			if got != tt.want {
				t.Errorf("defaultTCPCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultHTTPCheck(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		want       bool
		statusCode int
	}{
		{
			name:       "Valid HTTP connection",
			url:        "http://example.com",
			want:       true,
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid HTTP connection",
			url:        "http://invalid",
			want:       false,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Redirect HTTP connection",
			url:        "http://example.com/redirect",
			want:       false,
			statusCode: http.StatusMovedPermanently,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			// Override URL for valid connections
			if tt.want {
				tt.url = server.URL
			}

			addr, err := url.Parse(tt.url)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}

			got := defaultHTTPCheck(addr)
			if got != tt.want {
				t.Errorf("defaultHTTPCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultDNSCheck(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "Valid DNS resolution",
			url:  "http://example.com",
			want: true,
		},
		{
			name: "Invalid DNS resolution",
			url:  "http://invalid.invalid",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := url.Parse(tt.url)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}

			got := defaultDNSCheck(addr)
			if got != tt.want {
				t.Errorf("defaultDNSCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxyHealth_checkHealth(t *testing.T) {
	tests := []struct {
		name             string
		check            Check
		successThreshold int
		failureThreshold int
		failureCount     int
		wantAvailable    bool
		isAvailable      bool
	}{
		{
			name: "Health check passes",
			check: func(addr *url.URL) bool {
				return true
			},
			successThreshold: defaultSuccessThreshold,
			failureThreshold: defaultFailureThreshold,
			wantAvailable:    true,
			failureCount:     0,
		},
		{
			name: "Health check fails",
			check: func(addr *url.URL) bool {
				return false
			},
			successThreshold: defaultSuccessThreshold,
			failureThreshold: 1,
			wantAvailable:    false,
			failureCount:     0,
		},
		{
			name: "Health check passes after failures",
			check: func(addr *url.URL) bool {
				if addr.String() == "http://example.com" {
					return false
				}
				return true
			},
			successThreshold: defaultSuccessThreshold,
			failureThreshold: defaultFailureThreshold,
			wantAvailable:    true,
			isAvailable:      true,
			failureCount:     1,
		},
		{
			name: "failing health check after reaching failure threshold",
			check: func(addr *url.URL) bool {
				if addr.String() == "http://example.com" {
					return false
				}
				return true
			},
			successThreshold: defaultSuccessThreshold,
			failureThreshold: defaultFailureThreshold,
			wantAvailable:    false,
			isAvailable:      true,
			failureCount:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin, _ := url.Parse("http://example.com")
			h := &ProxyHealth{
				origin:           origin,
				check:            tt.check,
				successThreshold: tt.successThreshold,
				failureThreshold: tt.failureThreshold,
				cancel:           make(chan struct{}),
				failureCount:     tt.failureCount,
				isAvailable:      tt.isAvailable,
			}

			h.checkHealth()

			if h.isAvailable != tt.wantAvailable {
				t.Errorf("checkHealth() = %v, want %v", h.isAvailable, tt.wantAvailable)
			}
		})
	}
}
