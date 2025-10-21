package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Hello(t *testing.T) {
	srv := New()

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var got struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if got.Message != "Hello, World!" {
		t.Fatalf("unexpected message: %q", got.Message)
	}
}
