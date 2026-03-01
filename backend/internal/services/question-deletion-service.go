package services

type QuestionDeletionService struct {
	questionRepository QuestionRepository
}

func NewQuestionDeletionService(questionRepository QuestionRepository) *QuestionDeletionService {
	return &QuestionDeletionService{questionRepository: questionRepository}
}

func (s *QuestionDeletionService) DeleteQuestionByID(id int64) error {
	return s.questionRepository.DeleteQuestionByID(id)
}
