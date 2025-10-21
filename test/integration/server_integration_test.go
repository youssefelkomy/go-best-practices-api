package integration_test

import (
	"my-go-server/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerIntegration(t *testing.T) {
	srv := server.New()
	ts := httptest.NewServer(srv)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
