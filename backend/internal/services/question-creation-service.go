package services

import (
	"errors"

	"github.com/petersonvcode/questionnaire/backend/internal/domain"
)

type QuestionCreationService struct {
	questionRepository QuestionRepository
}

func NewQuestionCreationService(questionRepository QuestionRepository) *QuestionCreationService {
	return &QuestionCreationService{questionRepository: questionRepository}
}

func (s *QuestionCreationService) CreateQuestions(req domain.QuestionsCreationRequest) ([]*domain.Question, error) {
	valid, msg := req.IsValid()
	if !valid {
		return nil, errors.New(msg)
	}
	return s.questionRepository.InsertQuestions(req.Questions)
}
