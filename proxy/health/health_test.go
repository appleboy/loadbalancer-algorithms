package health

import (
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
		WithPeriod(50*time.Millisecond),
	)
	defer h.stop()

	// Wait for the initial check to complete
	time.Sleep(100 * time.Millisecond)

	if !h.IsAvailable() {
		t.Fatalf("Expected IsAvailable to be true, but got false")
	}

	// Change the origin to make the check fail
	h.origin, _ = url.Parse("http://invalid.com")

	// Wait for the next check to complete
	time.Sleep(100 * time.Millisecond)

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
