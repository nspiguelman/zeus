package domain

type AnswerMessage struct {
	QuestionId int    `json:"questionId"`
	AnswerId   int    `json:"answerId"`
	Token      string `json:"token"`
	IsTimeout  bool   `json:"isTimeout"`
}

type ScoreMessage struct {
	Score     int  `json:"score"`
	IsCorrect bool `json:"isCorrect"`
}

type AllScoreMessage struct {
	ScoreMessages map[string]ScoreMessage `json:"scores"`
	TypeMessage   string                  `json:"typeMessage"`
}

type QuestionMessage struct {
	QuestionId  int    `json:"questionId"`
	AnswerIds   []int  `json:"answerIds"`
	TypeMessage string `json:"typeMessage"`
}

func NewScoreMessage(score int, isCorrect bool) *ScoreMessage {
	return &ScoreMessage{
		Score:     score,
		IsCorrect: isCorrect,
	}
}
