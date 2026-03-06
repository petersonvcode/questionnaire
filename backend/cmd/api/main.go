package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/petersonvcode/questionnaire/backend/internal/database"
	"github.com/petersonvcode/questionnaire/backend/internal/repositories"
	"github.com/petersonvcode/questionnaire/backend/internal/server"
	"github.com/petersonvcode/questionnaire/backend/internal/services"
)

func main() {
	loadEnvFile()
	AdjustLogLevel()
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

	port := os.Getenv("QUESTIONNAIRE_PORT")
	if port == "" {
		port = "8080"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		slog.Error("Failed to convert port to int", "error", err)
		panic(err)
	}
	server := server.NewServer(portInt, questionCreationService, questionRetrievalService, questionDeletionService)
	return server
}

// Loads the environment variables from the .env and .env.<env> files
// .env.<env> takes precedence over .env
// Returns the environment variable ENV
func loadEnvFile() string {
	env := strings.ToLower(os.Getenv("ENV"))
	if env == "" {
		slog.Debug("ENV variable not set, defaulting to local")
		env = "local"
	}
	slog.Debug("ENV: " + env)

	// Loading environment variables from .env files
	environmentEnvFile := ".env." + env
	slog.Debug("loading environment variables from file", "environmentEnvFile", environmentEnvFile, "also", ".env")

	// Load .env file first
	godotenv.Load()
	// Then load the environment-specific file
	godotenv.Load(environmentEnvFile)
	slog.Debug("loaded environment variables from file", "environmentEnvFile", environmentEnvFile, "also", ".env")

	return env
}

func AdjustLogLevel() {
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelInfo)

	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		slog.Debug("LOG_LEVEL: " + levelStr)
		switch strings.ToUpper(levelStr) {
		case "DEBUG":
			logLevel.Set(slog.LevelDebug)
		case "INFO":
			logLevel.Set(slog.LevelInfo)
		case "WARN":
			logLevel.Set(slog.LevelWarn)
		case "ERROR":
			logLevel.Set(slog.LevelError)
		}
	}

	opts := &slog.HandlerOptions{Level: logLevel}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
}
