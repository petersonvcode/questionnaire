package services

import "github.com/petersonvcode/questionnaire/backend/internal/domain"

type QuestionRepository interface {
	InsertQuestions(questions []domain.QuestionCreationItemRequest) ([]*domain.Question, error)

	GetQuestionByID(id int64) (*domain.Question, error)
	GetRandomQuestion() (*domain.Question, error)

	DeleteQuestionByID(id int64) error
}
