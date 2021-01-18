package domain

import (
	"strings"
)

type Question struct {
	ID   int
	Question string
	Description string //TODO: Maybe it's not needed
	// options map[string]bool //Set of options
	// correctAnswer string
	// date *time.Time
}

func NewMultipleChoiceTrivia(question string, options map[string]bool, correctAnswer string, description string) *Question {
	// date := time.Now()
	correctAnswerFormatted := strings.ToLower(correctAnswer)
	optionsFormatted := make(map[string]bool)

	for k, v := range options {
		option := strings.ToLower(k)
		optionsFormatted[option] = v
	}

	if _, ok := optionsFormatted[correctAnswerFormatted]; !ok {
		panic("Correct answer is not an possible option");
	}

	return &(Question{Question: question, Description: description})
}

/*
func (t *Question) GetOptions() []string {
	v := make([]string, len(t.options))
	id := 0
	for k, _ := range t.options {
		v[id] = k
		id++
	}
	return v
}

func (t *Question) IsCorrect(option string) bool {
	if _, ok := t.options[option]; !ok {
		panic("Option does not exist");
	}

	return t.correctAnswer == option
}
*/