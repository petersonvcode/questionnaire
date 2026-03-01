package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	Port       int
	httpServer *http.Server

	// Dependencies
	questionCreationService  QuestionCreationService
	questionRetrievalService QuestionRetrievalService
	questionDeletionService  QuestionDeletionService
}

func NewServer(
	port int,
	questionCreationService QuestionCreationService,
	questionRetrievalService QuestionRetrievalService,
	questionDeletionService QuestionDeletionService,
) *Server {
	server := &Server{
		Port:                     port,
		questionCreationService:  questionCreationService,
		questionRetrievalService: questionRetrievalService,
		questionDeletionService:  questionDeletionService,
	}
	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: server.SetupRoutes(),
	}

	return server
}

func (s *Server) Start() error {
	slog.Info("Starting server", "port", s.Port)
	return s.httpServer.ListenAndServe()
}
