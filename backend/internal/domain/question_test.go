package domain

import "testing"

func TestQuestionsCreationRequest_isValid(t *testing.T) {
	twoOptionsOneCorrect := []QuestionOptionCreationRequest{
		{Text: "option1", Correct: true},
		{Text: "option2", Correct: false},
	}
	threeOptionsTwoCorrect := []QuestionOptionCreationRequest{
		{Text: "option1", Correct: true},
		{Text: "option2", Correct: true},
		{Text: "option3", Correct: false},
	}

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
			req:  QuestionsCreationRequest{Questions: []QuestionCreationRequestItem{{Text: ""}}},
			want: "question text cannot be empty",
		},
		{
			name: "unknown question type",
			req: QuestionsCreationRequest{Questions: []QuestionCreationRequestItem{
				{Text: "Some question?", Type: "true_false", Options: twoOptionsOneCorrect},
			}},
			want: "question 1: question type must be 'single_choice' or 'multiple_choice'",
		},
		{
			name: "question tags repeated",
			req: QuestionsCreationRequest{Questions: []QuestionCreationRequestItem{
				{Text: "Some question?", Type: "single_choice", Tags: []string{"tag1", "tag1"}},
			}},
			want: "question tags must not be repeated",
		},
		{
			name: "question options less than 2",
			req: QuestionsCreationRequest{Questions: []QuestionCreationRequestItem{
				{Text: "Some question?", Type: "single_choice", Options: []QuestionOptionCreationRequest{{Text: "option1", Correct: false}}},
			}},
			want: "each question must have between 2 and 8 options",
		},
		{
			name: "question options more than 8",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationRequestItem{
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
			name: "single choice with zero correct options",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationRequestItem{
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
			want: "question 1: there must be exactly one correct option",
		},
		{
			name: "single choice with two correct options",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationRequestItem{
					{
						Text:    "Some question?",
						Type:    "single_choice",
						Options: threeOptionsTwoCorrect,
					},
				},
			},
			want: "question 1: there must be exactly one correct option",
		},
		{
			name: "option text empty",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationRequestItem{
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
			want: "question 1: option text cannot be empty",
		},
		{
			name: "multiple choice with one correct option",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationRequestItem{
					{
						Text:    "Some question?",
						Type:    "multiple_choice",
						Options: twoOptionsOneCorrect,
					},
				},
			},
			want: "question 1: there must be at least two correct options",
		},
		{
			name: "multiple choice with zero correct options",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationRequestItem{
					{
						Text: "Some question?",
						Type: "multiple_choice",
						Options: []QuestionOptionCreationRequest{
							{Text: "option1", Correct: false},
							{Text: "option2", Correct: false},
						},
					},
				},
			},
			want: "question 1: there must be at least two correct options",
		},
		{
			name: "valid single choice",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationRequestItem{
					{
						Text:    "Some question?",
						Type:    "single_choice",
						Tags:    []string{"tag1"},
						Options: twoOptionsOneCorrect,
					},
				},
			},
			want: "",
		},
		{
			name: "valid multiple choice",
			req: QuestionsCreationRequest{
				Questions: []QuestionCreationRequestItem{
					{
						Text:    "Some question?",
						Type:    "multiple_choice",
						Options: threeOptionsTwoCorrect,
					},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, msg := tt.req.IsValid()
			wantValid := tt.want == ""
			if got != wantValid {
				t.Errorf("IsValid() = %v, want %v (msg %q)", got, wantValid, msg)
			}
			if msg != tt.want {
				t.Errorf("IsValid() msg = %q, want %q", msg, tt.want)
			}
		})
	}
}
