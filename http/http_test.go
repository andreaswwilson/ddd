package http_test

import (
	"net/http"
	"net/http/httptest"
)

// Return a test server
func setupTestServer(jsonResponse string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(int(status))
		if _, err := w.Write([]byte(jsonResponse)); err != nil {
			panic(err)
		}
	}))
}
