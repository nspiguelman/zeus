package domain

type AnswerMessage struct {
	QuestionId int    `json:"questionId"`
	AnswerId   int    `json:"answerId"`
	Token      string `json:"token"`
	IsTimeout  bool   `json:"isTimeout"`
}

type ScoreMessage struct {
	PartialScore int    `json:"partialScore"`
	IsCorrect    bool   `json:"isCorrect"`
	TypeMessage  string `json:"typeMessage"`
}

type QuestionMessage struct {
	QuestionId  int    `json:"questionId"`
	AnswerIds   []int  `json:"answerIds"`
	TypeMessage string `json:"typeMessage"`
}