package main

import (
	"log"
	"my-go-server/internal/server"
)

func main() {
	srv := server.New()
	log.Println("Starting server on :8080")
	if err := srv.Start(); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
