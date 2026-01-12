// handlers_test.go
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var publicTests = []struct {
	name   string
	method string
	path   string
	body   []byte
	expect int
}{
	{
		name:   "Register with invalid JSON",
		method: "POST",
		path:   "/register",
		body:   []byte(`{"email": "test", "password": }`),
		expect: http.StatusBadRequest,
	},
	{
		name:   "Login with invalid JSON",
		method: "POST",
		path:   "/login",
		body:   []byte(`{"email": "test", "password": }`),
		expect: http.StatusBadRequest,
	},
	{
		name:   "Register with empty body",
		method: "POST",
		path:   "/register",
		body:   []byte(``),
		expect: http.StatusBadRequest,
	},
}

func TestPublicEndpoints(t *testing.T) {
	for _, tt := range publicTests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			switch tt.path {
			case "/register":
				Register(w, req)
			case "/login":
				Login(w, req)
			}

			if w.Code != tt.expect {
				t.Errorf("Expected %d, got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestProtectedEndpoints_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test skipped in short mode")
	}
}
