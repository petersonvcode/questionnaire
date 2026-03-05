package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/petersonvcode/questionnaire/backend/internal/domain"
)

func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.healthHandler)
	mux.HandleFunc("/questions", s.questionsHandler)
	mux.HandleFunc("/questions/{id}", s.questionHandler)
	return s.corsMiddleware(mux)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) questionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		question, err := s.questionRetrievalService.GetRandomQuestion()
		if err != nil {
			slog.Error("Failed to get random question", "error", err)
			http.Error(w, "Failed to get random question", http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(question)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	case http.MethodPost:
		var req domain.QuestionsCreationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}
		questions, err := s.questionCreationService.CreateQuestions(req)
		if err != nil {
			slog.Error("Failed to create questions", "error", err)
			http.Error(w, "Failed to create questions", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(questions)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s *Server) questionHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		question, err := s.questionRetrievalService.GetQuestionByID(idInt)
		if err != nil {
			slog.Error("Failed to get question", "error", err)
			http.Error(w, "Failed to get question", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(question)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	case http.MethodDelete:
		err := s.questionDeletionService.DeleteQuestionByID(idInt)
		if err != nil {
			slog.Error("Failed to delete question", "error", err)
			http.Error(w, "Failed to delete question", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
