package domain

import "fmt"

// Database model - Table questions
type Question struct {
	ID      int64            `json:"id"`
	Text    string           `json:"text"`
	Type    string           `json:"type"`
	Options []QuestionOption `json:"options"`
	Tags    []QuestionTag    `json:"tags"`
}

// Database model - Table question_options
type QuestionOption struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"-"`
	Text       string `json:"text"`
	Correct    bool   `json:"correct"`
}

// Database model - Table question_tags
type QuestionTag struct {
	ID    int64  `json:"id"`
	Tag   string `json:"tag"`
	Group string `json:"group"`
}

// API request model - Question creation
type QuestionsCreationRequest struct {
	Questions []QuestionCreationRequestItem `json:"questions"`
}

// API request model - Question creation
type QuestionCreationRequestItem struct {
	Text    string                          `json:"text"`
	Type    string                          `json:"type"`
	Tags    []string                        `json:"tags"`
	Options []QuestionOptionCreationRequest `json:"options"`
}

// API request model - Question creation
type QuestionOptionCreationRequest struct {
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

// API response model
type QuestionAPIResponse struct {
	ID      int64                       `json:"id"`
	Text    string                      `json:"text"`
	Type    string                      `json:"type"`
	Tags    []string                    `json:"tags"`
	Options []QuestionOptionAPIResponse `json:"options"`
}

// API response model - QuestionOption
type QuestionOptionAPIResponse struct {
	ID      int64  `json:"id"`
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

// API response model - QuestionTag
type QuestionTagAPIResponse struct {
	ID    int64  `json:"id"`
	Tag   string `json:"tag"`
	Group string `json:"-"`
	Count int    `json:"count"`
}

const (
	MAX_OPTIONS_PER_QUESTION = 8
	QUESTION_TEXT_MAX_LENGTH = 4096
	TAG_MAX_LENGTH           = 64
	TAG_GROUP_MAX_LENGTH     = 64
	OPTION_TEXT_MAX_LENGTH   = 1024
)

// Validates the questions creation request.
// Returns true if the request is valid, false otherwise.
// Returns the error message if the request is not valid.
func (q QuestionsCreationRequest) IsValid() (bool, string) {
	if len(q.Questions) == 0 {
		return false, "at least one question is required"
	}
	for qIndex, question := range q.Questions {
		if question.Text == "" {
			return false, "question text cannot be empty"
		}
		if len(question.Text) > QUESTION_TEXT_MAX_LENGTH {
			return false, fmt.Sprintf("question text cannot be longer than %d characters", QUESTION_TEXT_MAX_LENGTH)
		}

		tagSet := make(map[string]struct{})
		for _, t := range question.Tags {
			if len(t) > TAG_MAX_LENGTH {
				return false, fmt.Sprintf("tag cannot be longer than %d characters", TAG_MAX_LENGTH)
			}
			if _, exists := tagSet[t]; exists {
				return false, "question tags must not be repeated"
			}
			tagSet[t] = struct{}{}
		}
		if len(question.Tags) > TAG_GROUP_MAX_LENGTH {
			return false, fmt.Sprintf("question tags cannot be longer than %d characters", TAG_MAX_LENGTH)
		}

		if len(question.Options) < 2 || len(question.Options) > MAX_OPTIONS_PER_QUESTION {
			return false, fmt.Sprintf("each question must have between 2 and %d options", MAX_OPTIONS_PER_QUESTION)
		}

		switch question.Type {
		case "single_choice":
			isValid, msg := question.IsSingleChoiceQuestionValid()
			if !isValid {
				return false, fmt.Sprintf("question %d: %s", qIndex+1, msg)
			}
		case "multiple_choice":
			isValid, msg := question.IsMultipleChoiceQuestionValid()
			if !isValid {
				return false, fmt.Sprintf("question %d: %s", qIndex+1, msg)
			}
		default:
			return false, fmt.Sprintf("question %d: question type must be 'single_choice' or 'multiple_choice'", qIndex+1)
		}
	}

	return true, ""
}

func (q QuestionCreationRequestItem) IsSingleChoiceQuestionValid() (bool, string) {
	correctCount := 0
	for _, opt := range q.Options {
		if opt.Text == "" {
			return false, "option text cannot be empty"
		}
		if len(opt.Text) > OPTION_TEXT_MAX_LENGTH {
			return false, fmt.Sprintf("option text cannot be longer than %d characters", OPTION_TEXT_MAX_LENGTH)
		}
		if opt.Correct {
			correctCount++
		}
	}
	if correctCount != 1 {
		return false, "there must be exactly one correct option"
	}
	return true, ""
}

func (q QuestionCreationRequestItem) IsMultipleChoiceQuestionValid() (bool, string) {
	correctCount := 0
	for _, opt := range q.Options {
		if opt.Text == "" {
			return false, "option text cannot be empty"
		}
		if len(opt.Text) > OPTION_TEXT_MAX_LENGTH {
			return false, fmt.Sprintf("option text cannot be longer than %d characters", OPTION_TEXT_MAX_LENGTH)
		}
		if opt.Correct {
			correctCount++
		}
	}
	if correctCount < 2 {
		return false, "there must be at least two correct options"
	}
	return true, ""
}
