package domain

type Question struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Question    string `gorm:"not null" json:"question"`
	Description string `gorm:"not null" json:"description"`
	KahootID    int    `gorm:"not null" json:"kahootId"`
}

type Answer struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Description string `json:"description"`
	QuestionID  int    `gorm:"not null" json:"questionId"`
	IsTrue      bool   `json:"isTrue"`
}

func NewQuestion(kahootID int, input QuestionInput) *Question{
	return &Question{
		Question:    input.Question,
		Description: input.Description,
		KahootID:    kahootID,
	}
}

func NewAnswer(questionID int, input AnswerInput) *Answer {
	return &Answer{
		Description: input.Description,
		QuestionID:  questionID,
		IsTrue:      input.IsTrue,
	}
}
