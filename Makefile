# Makefile

.PHONY: build run test clean

build:
	go build -o my-go-server ./cmd/server

run: build
	./my-go-server

test:
	go test ./...

clean:
	go clean
	rm -f my-go-server