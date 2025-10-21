# My Go Server

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

### Prerequisites

- Go 1.16 or later
- Make (optional, for using the Makefile)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/my-go-server.git
   cd my-go-server
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

### Running the Server

To run the server, execute the following command:

```
go run cmd/server/main.go
```

The server will start listening on the specified port (default is 8080).

### Running Tests

To run unit tests, use:

```
go test ./...
```

For integration tests, run:

```
go test ./test/integration
```

### API Documentation

The API endpoints are documented in the `api/openapi.yaml` file. You can use tools like Swagger UI to visualize and interact with the API.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

### License

This project is licensed under the MIT License. See the LICENSE file for details.