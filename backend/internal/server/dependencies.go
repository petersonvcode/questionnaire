package server

import "github.com/petersonvcode/questionnaire/backend/internal/domain"

type QuestionCreationService interface {
	CreateQuestions(req domain.QuestionsCreationRequest) ([]*domain.Question, error)
}

type QuestionRetrievalService interface {
	GetQuestionByID(id int64) (*domain.Question, error)
	GetRandomQuestion() (*domain.Question, error)
	GetTags() ([]domain.QuestionTagAPIResponse, error)
	CountQuestionsWithTags(tagIDs []int64) (int, error)
}

type QuestionDeletionService interface {
	DeleteQuestionByID(id int64) error
}
