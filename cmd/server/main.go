package main

import (
    "log"
    "net/http"
    "my-go-server/internal/server"
)

func main() {
    srv := server.NewServer()
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", srv); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}