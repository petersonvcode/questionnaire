package services

import "github.com/petersonvcode/questionnaire/backend/internal/domain"

type QuestionRetrievalService struct {
	questionRepository QuestionRepository
}

func NewQuestionRetrievalService(questionRepository QuestionRepository) *QuestionRetrievalService {
	return &QuestionRetrievalService{questionRepository: questionRepository}
}

func (s *QuestionRetrievalService) GetQuestionByID(id int64) (*domain.Question, error) {
	return s.questionRepository.GetQuestionByID(id)
}

func (s *QuestionRetrievalService) GetRandomQuestion() (*domain.Question, error) {
	return s.questionRepository.GetRandomQuestion()
}
