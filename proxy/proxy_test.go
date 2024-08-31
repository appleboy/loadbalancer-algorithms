package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestProxy_ServeHTTP(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new proxy with the test server's URL
	proxyURL, _ := url.Parse(ts.URL)
	proxy := NewProxy("foobar", proxyURL)

	// Create a mock request
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method
	proxy.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}
