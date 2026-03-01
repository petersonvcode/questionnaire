package main

import (
	"log/slog"

	"github.com/petersonvcode/questionnaire/backend/internal/database"
	"github.com/petersonvcode/questionnaire/backend/internal/repositories"
	"github.com/petersonvcode/questionnaire/backend/internal/server"
	"github.com/petersonvcode/questionnaire/backend/internal/services"
)

func main() {
	server := setupServerWithDependencies()
	err := server.Start()
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}

func setupServerWithDependencies() *server.Server {
	db, err := database.New()
	if err != nil {
		slog.Error("Failed to create database", "error", err)
		panic(err)
	}

	// repositories
	questionRepository := repositories.NewQuestionRepository(db)

	// services
	questionCreationService := services.NewQuestionCreationService(questionRepository)
	questionRetrievalService := services.NewQuestionRetrievalService(questionRepository)
	questionDeletionService := services.NewQuestionDeletionService(questionRepository)

	server := server.NewServer(8080, questionCreationService, questionRetrievalService, questionDeletionService)
	return server
}
