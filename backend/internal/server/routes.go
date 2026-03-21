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
	mux.HandleFunc("/questions/tags", s.tagsHandler)
	mux.HandleFunc("/questions/tags/count", s.tagsCountHandler)
	return s.corsMiddleware(mux)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func isGetTagsQuestionsRequest(r *http.Request) (valid bool, tagIDs []int64, count int) {
	if r.Method != http.MethodGet {
		return false, nil, 0
	}
	tagIDsStr := r.URL.Query().Get("tagIDs")
	tagIDs, err := domain.ParseInt64QueryString(tagIDsStr)
	if err != nil {
		return false, nil, 0
	}
	if len(tagIDs) <= 0 {
		return false, nil, 0
	}

	countStr := r.URL.Query().Get("count")
	count, err = strconv.Atoi(countStr)
	if err != nil {
		return false, nil, 0
	}
	if count <= 0 {
		return false, nil, 0
	}
	if count > 100 {
		return true, tagIDs, 100
	}
	return true, tagIDs, count
}

func (s *Server) questionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		isTagsRequest, tagIDs, count := isGetTagsQuestionsRequest(r)
		if isTagsRequest {
			questions, err := s.questionRetrievalService.GetQuestionsWithTags(tagIDs, count)
			if err != nil {
				slog.Error("Failed to get questions with tags", "error", err)
				http.Error(w, "Failed to get questions with tags", http.StatusInternalServerError)
				return
			}
			resp, err := json.Marshal(questions)
			if err != nil {
				http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		} else {
			// Get random question
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
		}

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

func (s *Server) tagsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tags, err := s.questionRetrievalService.GetTags()
		if err != nil {
			slog.Error("Failed to get tags", "error", err)
			http.Error(w, "Failed to get tags", http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(tags)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s *Server) tagsCountHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tagIDsStr := r.URL.Query().Get("ids")
		tagIDs, err := domain.ParseInt64QueryString(tagIDsStr)
		if err != nil {
			http.Error(w, "Invalid tag IDs", http.StatusBadRequest)
			return
		}
		count, err := s.questionRetrievalService.CountQuestionsWithTags(tagIDs)
		if err != nil {
			slog.Error("Failed to count questions with tags", "error", err)
			http.Error(w, "Failed to count questions with tags", http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(count)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
