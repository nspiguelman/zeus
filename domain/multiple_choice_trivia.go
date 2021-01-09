package domain

import "time"

var idIncremental int = 0

type MultipleChoiceTrivia struct {
	id   int
	question string
	date *time.Time
}

func NewMultipleChoiceTrivia(question string) *MultipleChoiceTrivia {
	date := time.Now()
	idIncremental++
	return &(MultipleChoiceTrivia{idIncremental, question, &date})
}

func (trivia *MultipleChoiceTrivia) GetQuestion() string {
	return trivia.question
}

func (trivia *MultipleChoiceTrivia) GetId() int {
	return trivia.id
}

func (trivia *MultipleChoiceTrivia) GetDate() *time.Time {
	return trivia.date
}
