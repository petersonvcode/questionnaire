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

func (s *QuestionRetrievalService) GetTags() ([]domain.QuestionTagAPIResponse, error) {
	return s.questionRepository.GetTags()
}

func (s *QuestionRetrievalService) CountQuestionsWithTags(tagIDs []int64) (int, error) {
	return s.questionRepository.CountQuestionsWithTags(tagIDs)
}

func (s *QuestionRetrievalService) GetQuestionsWithTags(tagIDs []int64, count int) ([]domain.Question, error) {
	return s.questionRepository.GetQuestionsWithTags(tagIDs, count)
}
