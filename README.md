# Go Best Practices API

This project is a simple HTTP server built with Go, designed to demonstrate best practices in server implementation and testing.

## Project Structure

```
my-go-server
├── cmd
│   └── server
│       └── main.go          # Entry point of the application
├── internal
│   ├── server
│   │   ├── server.go        # Server implementation
│   │   └── server_test.go   # Unit tests for the server
│   └── config
│       └── config.go        # Configuration settings
├── pkg
│   └── handler
│       ├── handler.go       # HTTP handler functions
│       └── handler_test.go   # Unit tests for handlers
├── api
│   └── openapi.yaml         # API specification
├── test
│   └── integration
│       └── server_integration_test.go # Integration tests
├── .github
│   └── workflows
│       └── ci.yml           # CI workflow
├── go.mod                   # Module definition
├── Makefile                 # Build and test commands
└── README.md                # Project documentation
```

## Getting Started

# my-go-server

Lightweight example Go HTTP server (built with Gin) intended for local development,
experimentation and benchmarking. It includes a small set of endpoints, an embedded
OpenAPI spec and an interactive documentation UI (Redoc).

This README documents how to run, test, benchmark and profile the server.

---

Table of contents
- Overview and routes
- Quick start
- Build & run (release)
- Tests
- Docs (OpenAPI + Redoc)
- Benchmarking and profiling
- Notes & recommendations

## Overview and routes

Main implementation: `internal/server/server.go` (Gin engine).

Exposed endpoints (short):

- `GET /` — Service metadata (service, version, uptime, requests)
- `GET /info` — Start time and basic info
- `GET /time` — Current server time (RFC3339Nano)
- `GET /health` — Simple liveness (keeps compatibility with `pkg/handler`)
- `GET /metrics` — In-memory metrics (uptime, total requests)
- `GET /headers` — Echo request headers (diagnostic)
- `ANY /echo` — Echo method, query, headers and raw body
- `GET /hello` — Greeting JSON (keeps compatibility with `pkg/handler`)
- `GET /openapi.yaml` — Embedded OpenAPI specification (YAML)
- `GET /docs` — Interactive Redoc documentation (embedded spec)

Notes:
- `pkg/handler` contains small net/http-compatible handlers used for `/health` and `/hello`.
- The OpenAPI spec is embedded into the binary for reliable docs serving.

## Quick start (development)

Prerequisites:
- Go 1.20+ (this repo was tested with Go 1.25)

Bootstrap (one-time):

```bash
cd your-route/projects/go-best-practices-api
go mod tidy
```

Run (quick):

```bash
# development (prints Gin debug logs)
go run ./cmd/server
```

Open http://localhost:8080/docs in your browser to view interactive API docs.

## Build & run (release mode)

For benchmarking or running without Gin debug logs, use release mode:

```bash
cd your-route/projects/go-best-practices-api
export GIN_MODE=release
go build -o my-go-server ./cmd/server
./my-go-server &> server.log &
echo $!  # PID
```

Then visit:

- http://localhost:8080/ — service info
- http://localhost:8080/docs — interactive docs (Redoc)
- http://localhost:8080/openapi.yaml — raw OpenAPI spec

## Tests

Run all unit & integration tests:

```bash
go test ./...
```

Run only integration tests:

```bash
go test ./test/integration
```

## Docs (OpenAPI & Redoc)

- The OpenAPI spec file is maintained at `api/openapi.yaml` and a copy is embedded at
   `internal/server/openapi.yaml` so the binary serves the spec even when the working directory changes.
- Interactive docs are available at `/docs` (Redoc) and are served using an embedded copy of the spec
   so they work offline and without an extra network fetch.

## Benchmarking and profiling (recommended)

Start the server in release mode (see above). Use a modern benchmarking tool — `wrk`, `hey` or `vegeta`
— instead of `ab` when testing high concurrency.

Example `wrk` run (local loopback):

```bash
# moderate concurrency
wrk2 -t4 -c200 -d30s http://127.0.0.1:8080/info

# increase until throughput plateaus
wrk2 -t4 -c400 -d30s http://127.0.0.1:8080/info
```

Example `hey`:

```bash
hey -n 200000 -c 200 http://127.0.0.1:8080/info
```

If you push concurrency very high (thousands), raise OS limits first:

```bash
# temporary for current shell
ulimit -n 100000
sudo sysctl -w net.core.somaxconn=65535
sudo sysctl -w net.ipv4.tcp_max_syn_backlog=4096
sudo sysctl -w net.ipv4.ip_local_port_range="1024 65000"
```

Profiling with pprof (developer only)

1. Temporarily enable `net/http/pprof` in `cmd/server/main.go` (dev only):

```go
import _ "net/http/pprof"
go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()
```

2. Start server + benchmark, then collect a CPU profile:

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
# in pprof: (pprof) top  ; (pprof) web
```

3. Inspect heap profile similarly using `/debug/pprof/heap`.

## Practical tips & tuning

- Run Gin in release mode when benchmarking: `export GIN_MODE=release`.
- Use keep-alive in clients (wrk/hey/vegeta support this) to avoid connection churn.
- Route or disable the Gin logger when benchmarking; logging to stdout can hurt throughput.
- Pre-encode static JSON payloads if possible to avoid allocations in hot handlers (e.g. a prebuilt `[]byte` for `{"message":"Hello, World!"}`).
- Profile before optimizing: use pprof to find real hotspots (GC, allocation, syscalls).

## Contributing

Contributions are welcome — open an issue or a PR. If you add new endpoints, please update `api/openapi.yaml`.

## License

MIT
