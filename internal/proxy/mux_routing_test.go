package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMuxRouting(t *testing.T) {
	// Simulate the exact main.go mux registration order
	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("login"))
	})
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("portal health"))
	})
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("proxy -> " + r.URL.Path))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("dashboard"))
	})
	mux.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("portal services"))
	})
	mux.HandleFunc("/api/probe", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("portal probe"))
	})
	mux.HandleFunc("/proxy/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("proxy -> " + r.URL.Path))
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("logout"))
	})
	mux.HandleFunc("/api/me", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("portal me"))
	})

	tests := []struct {
		path     string
		wantBody string
	}{
		{"/api/health", "portal health"},
		{"/api/services", "portal services"},
		{"/api/probe", "portal probe"},
		{"/api/me", "portal me"},
		{"/api/funds", "proxy -> /api/funds"},
		{"/api/portfolio/summary", "proxy -> /api/portfolio/summary"},
		{"/api/", "proxy -> /api/"},
		{"/proxy/5000/dashboard", "proxy -> /proxy/5000/dashboard"},
		{"/proxy/5000/api/funds", "proxy -> /proxy/5000/api/funds"},
		{"/login", "login"},
		{"/", "dashboard"},
		{"/logout", "logout"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "http://localhost:8747"+tt.path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		body, _ := io.ReadAll(w.Body)
		if string(body) != tt.wantBody {
			t.Errorf("GET %s: got %q, want %q", tt.path, string(body), tt.wantBody)
		} else {
			fmt.Printf("  OK GET %-35s -> %s\n", tt.path, string(body))
		}
	}
}
