package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Return a test server
func ServeMux(t *testing.T) (*http.ServeMux, *Client) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClient("", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return mux, client
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %s, want %s", got, want)
	}
}
