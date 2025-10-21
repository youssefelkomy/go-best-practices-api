package integration_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "my-go-server/internal/server"
)

func TestServerIntegration(t *testing.T) {
    // Initialize the server
    srv := server.New() // Assuming New() initializes and returns a server instance
    ts := httptest.NewServer(srv)
    defer ts.Close()

    // Define test cases
    tests := []struct {
        name       string
        method     string
        url        string
        wantStatus int
    }{
        {"GET /", http.MethodGet, ts.URL + "/", http.StatusOK},
        // Add more test cases as needed
    }

    // Execute test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := http.NewRequest(tt.method, tt.url, nil)
            if err != nil {
                t.Fatalf("could not create request: %v", err)
            }

            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                t.Fatalf("could not send request: %v", err)
            }
            defer resp.Body.Close()

            if resp.StatusCode != tt.wantStatus {
                t.Errorf("got status %v, want %v", resp.StatusCode, tt.wantStatus)
            }
        })
    }
}