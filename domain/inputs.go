package domain

type KahootInput struct {
	Name string `json:"name" validate:"min=1,max=50"`
}

type UserInput struct {
	Username string `json:"username" validate:"min=4, max=15"`
	//PIN
}

type QuestionInput struct {
	Question    string        `json:"question" validate:"min=1,max=255"`
	Description string        `json:"description" validate:"min=1,max=255"`
	Answers     []AnswerInput `json:"answers" validate:"min=2,max=4"`
}

type AnswerInput struct {
	Description string `json:"description"`
	IsTrue      bool   `json:"is_true"`
}