package server

import (
    "net/http"
)

type Server struct {
    addr string
}

func NewServer(addr string) *Server {
    return &Server{addr: addr}
}

func (s *Server) Start() error {
    http.HandleFunc("/", s.handleRequest)
    return http.ListenAndServe(s.addr, nil)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello, World!"))
}