package domain

import "testing"

func TestQuestionsCreationRequest_isValid(t *testing.T) {
	tests := []struct {
		name string
		req  QuestionsCreationRequest
		want string
	}{
		{
			name: "empty request",
			req:  QuestionsCreationRequest{},
			want: "at least one question is required",
		},
		{
			name: "question text empty",
			req:  QuestionsCreationRequest{Questions: []QuestionCreationItemRequest{{Text: ""}}},
			want: "question text cannot be empty",
		},
		{
			name: "question type not single_choice",
			req:  QuestionsCreationRequest{Questions: []QuestionCreationItemRequest{{Text: "Some question?", Type: "multiple_choice"}}},
			want: "question type must be 'single_choice'",
		},
		{
			name: "question tags repeated",
			req:  QuestionsCreationRequest{Questions: []QuestionCreationItemRequest{{Text: "Some question?", Type: "single_choice", Tags: []string{"tag1", "tag1"}}}},
			want: "question tags must not be repeated",
		},
		{
			name: "question options less than 2",
			req:  QuestionsCreationRequest{Questions: []QuestionCreationItemRequest{{Text: "Some question?", Type: "single_choice", Options: []QuestionOptionCreationRequest{{Text: "option1", Correct: false}}}}},
			want: "each question must have between 2 and 8 options",
		},
		{
			name: "question options more than 8",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationItemRequest{
					{
						Text: "Some question?",
						Type: "single_choice",
						Options: []QuestionOptionCreationRequest{
							{Text: "option1", Correct: false},
							{Text: "option2", Correct: false},
							{Text: "option3", Correct: false},
							{Text: "option4", Correct: false},
							{Text: "option5", Correct: false},
							{Text: "option6", Correct: false},
							{Text: "option7", Correct: false},
							{Text: "option8", Correct: false},
							{Text: "option9", Correct: false},
						},
					},
				},
			},
			want: "each question must have between 2 and 8 options",
		},
		{
			name: "question options correct count not 1 (zero correct)",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationItemRequest{
					{
						Text: "Some question?",
						Type: "single_choice",
						Options: []QuestionOptionCreationRequest{
							{Text: "option1", Correct: false},
							{Text: "option2", Correct: false},
							{Text: "option3", Correct: false},
						},
					},
				},
			},
			want: "there must be exactly one correct option",
		},
		{
			name: "question options correct count not 1 (two correct)",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationItemRequest{
					{
						Text: "Some question?",
						Type: "single_choice",
						Options: []QuestionOptionCreationRequest{
							{Text: "option1", Correct: true},
							{Text: "option2", Correct: true},
							{Text: "option3", Correct: false},
						},
					},
				},
			},
			want: "there must be exactly one correct option",
		},
		{
			name: "option text empty",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationItemRequest{
					{
						Type: "single_choice",
						Text: "Some question?",
						Options: []QuestionOptionCreationRequest{
							{Text: "option1", Correct: false},
							{Text: "", Correct: true},
						},
					},
				},
			},
			want: "option text cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.req.IsValid()
			if got {
				t.Errorf("QuestionsCreationRequest.isValid() = true, want error: %v", tt.want)
			}
			if err != tt.want {
				t.Errorf("QuestionsCreationRequest.isValid() error = %v, want %v", err, tt.want)
			}
		})
	}
}
