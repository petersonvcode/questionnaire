package main

import (
	"log/slog"

	"github.com/petersonvcode/questionnaire/backend/internal/domain"
)

func main() {
	server := domain.NewServer(8080)
	err := server.Start()
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
