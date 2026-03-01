package domain

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
	Questions []QuestionCreationItemRequest `json:"questions"`
}

// API request model - Question creation
type QuestionCreationItemRequest struct {
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
	Group string `json:"group"`
}

// Validates the questions creation request.
// Returns true if the request is valid, false otherwise.
// Returns the error message if the request is not valid.
func (q QuestionsCreationRequest) IsValid() (bool, string) {
	if len(q.Questions) == 0 {
		return false, "at least one question is required"
	}
	for _, item := range q.Questions {
		if item.Text == "" {
			return false, "question text cannot be empty"
		}
		if item.Type != "single_choice" {
			return false, "question type must be 'single_choice'"
		}
		tagSet := make(map[string]struct{})
		for _, t := range item.Tags {
			if _, exists := tagSet[t]; exists {
				return false, "question tags must not be repeated"
			}
			tagSet[t] = struct{}{}
		}
		if len(item.Options) < 2 || len(item.Options) > 8 {
			return false, "each question must have between 2 and 8 options"
		}
		correctCount := 0
		for _, opt := range item.Options {
			if opt.Text == "" {
				return false, "option text cannot be empty"
			}
			if opt.Correct {
				correctCount++
			}
		}
		if correctCount != 1 {
			return false, "there must be exactly one correct option"
		}
	}
	return true, ""
}
