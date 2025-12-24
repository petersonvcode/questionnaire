package domain

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	Port       int
	httpServer *http.Server
}

func NewServer(port int) *Server {
	return &Server{
		Port: port,
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
	}
}

func (s *Server) Start() error {
	slog.Info("Starting server", "port", s.Port)
	return s.httpServer.ListenAndServe()
}
