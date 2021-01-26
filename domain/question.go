package domain

type Question struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Question    string `json:"question"`
	Description string `json:"description"`
	KahootID    int    `gorm:"not null" json:"kahoot_id"`
}

type Answer struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	AnswerID    int    `gorm:"primaryKey" json:"answer_id"`
	Description string `json:"description"`
	QuestionID  int    `gorm:"not null" json:"question_id"`
	IsTrue      bool   `json:"is_true"`
}

func NewQuestion(kahootID int, input QuestionInput) *Question{
	return &Question{
		Question:    input.Question,
		Description: input.Description,
		KahootID:    kahootID,
	}
}

func NewAnswer(questionID int, answerID int, input AnswerInput) *Answer {
	return &Answer{
		AnswerID:    answerID,
		Description: input.Description,
		QuestionID:  questionID,
		IsTrue:      input.IsTrue,
	}
}