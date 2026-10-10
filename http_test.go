package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"
)

func TestAddress(t *testing.T) {
	for _, tc := range []struct {
		port, want string
		invalid    bool
	}{{"", ":8991", false}, {"8080", ":8080", false}, {"1", ":1", false}, {"65535", ":65535", false}, {"0", "", true}, {"65536", "", true}, {"abc", "", true}, {"-1", "", true}} {
		got, err := address(tc.port)
		if (err != nil) != tc.invalid || got != tc.want {
			t.Errorf("address(%q) = %q, %v", tc.port, got, err)
		}
	}
}

func TestHTTPContract(t *testing.T) {
	h := handler(slog.New(slog.NewJSONHandler(io.Discard, nil)))
	for _, tc := range []struct {
		method, path string
		status       int
	}{{"GET", "/", 200}, {"GET", "/healthz", 200}, {"GET", "/readyz", 200}, {"HEAD", "/", 200}, {"GET", "/missing", 404}, {"POST", "/", 405}, {"POST", "/healthz", 405}} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status {
				t.Fatalf("status %d, want %d", w.Code, tc.status)
			}
			if w.Header().Get("X-Content-Type-Options") != "nosniff" || w.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("missing response safeguards")
			}
			if tc.method == "GET" && tc.status == 200 {
				var body map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if tc.path == "/" && body["service"] != "Portlight" {
					t.Fatal("wrong service identity")
				}
				if tc.path != "/" && body["status"] != "ok" {
					t.Fatal("probe not healthy")
				}
			}
		})
	}
}
