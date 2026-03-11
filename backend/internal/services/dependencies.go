package services

import "github.com/petersonvcode/questionnaire/backend/internal/domain"

type QuestionRepository interface {
	InsertQuestions(questions []domain.QuestionCreationItemRequest) ([]*domain.Question, error)

	GetQuestionByID(id int64) (*domain.Question, error)
	GetRandomQuestion() (*domain.Question, error)

	GetTags() ([]domain.QuestionTagAPIResponse, error)
	CountQuestionsWithTags(tagIDs []int64) (int, error)

	DeleteQuestionByID(id int64) error
}
