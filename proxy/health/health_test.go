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

	h := NewProxyHealth(
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
