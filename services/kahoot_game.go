package services

import (
	"errors"
	"github.com/nspiguelman/zeus/data"
	"github.com/nspiguelman/zeus/domain"
	"sync"
)

var (
	game *KahootGame
	once sync.Once
)

type KahootGame struct {
	kahoot          *domain.Kahoot
	questions       []domain.Question
	answers         []domain.Answer
	host            domain.User
	users           []domain.User
	pin             string
	CurrentQuestion int  // pregunta actual
	CurrentQuestionIndex int // para iterar las questions
	TotalQuestions  int  // total de preguntas
	IsTimeout       bool // timeout de la pregunta
	IsStarted       bool // ya empezo el kahoot, no se puede suscribir nadie
	IsScoreSent     bool // la siguiente pregunta solo puede ser enviada cuando el score sea notificado
	rm              *data.RepositoryManager
}

func initGame() {
	game = &KahootGame{
		questions:       make([]domain.Question, 0),
		answers:         make([]domain.Answer, 0),
		host:            domain.User{}, // TODO: Revisar si vamos a tener un host
		users:           make([]domain.User, 0),
		pin:             "",
		CurrentQuestion: 0,
		CurrentQuestionIndex: 0,
		TotalQuestions:  0,
		IsTimeout:       false,
		IsStarted:       false,
		IsScoreSent:     false,
	}
}

func NewKahootGame(rm *data.RepositoryManager) *KahootGame {
	once.Do(func() {
		initGame()
		game.rm = rm
		game.IsStarted = false
	})

	return game
}

func (kg *KahootGame) Start(pin string) error {
	kg.pin = pin
	if err := kg.searchKahoot(); err != nil {
		return errors.New("An error occurred while getting kahoot: " + err.Error())
	}

	if err := kg.searchQuestions(); err != nil {
		return errors.New("An error occurred while getting questions: " + err.Error())
	}

	if err := kg.searchAnswers(kg.questions[kg.CurrentQuestionIndex].ID); err != nil {
		return errors.New("An error occurred while getting answers: " + err.Error())
	}

	kg.CurrentQuestion = kg.questions[kg.CurrentQuestionIndex].ID
	kg.IsStarted = true
	kg.IsScoreSent = false
	kg.IsTimeout = false
	return nil
}

func (kg *KahootGame) NextQuestion(pin string) error {
	if kg.CurrentQuestionIndex + 1 >= kg.TotalQuestions {
		//TODO: enviar algún evento para game over, por ahora devolvemos error
		return errors.New("Game Over")
	}

	kg.CurrentQuestionIndex++
	if err := kg.searchAnswers(kg.questions[kg.CurrentQuestionIndex].ID); err != nil {
		return errors.New("An error occurred while getting answers for question " + string(kg.CurrentQuestion) +  ": " + err.Error())
	}

	kg.CurrentQuestion = kg.questions[kg.CurrentQuestionIndex].ID
	return nil
}

func (kg *KahootGame) searchKahoot() error {
	kahoot, err := kg.rm.KahootRepository.GetByPin(kg.pin);
	if kahoot == nil && err == nil {
		return errors.New("Room " + kg.pin + " not found")
	}
	if err != nil {
		return err
	}

	kg.kahoot = kahoot
	return nil
}

func (kg *KahootGame) searchQuestions() error {
	questions, err := kg.rm.QuestionRepository.GetAllByKahootID(kg.kahoot.ID)
	if len(questions) == 0 && err == nil {
		return errors.New("Questions not found for kahoot: " + kg.pin)
	}
	if err != nil {
		return err
	}

	kg.questions = questions
	kg.TotalQuestions = len(questions)
	return nil
}

func (kg *KahootGame) searchAnswers(questionID int) error {
	answers, err := kg.rm.AnswerRepository.GetAllByQuestionID(questionID)
	if len(answers) == 0 && err == nil {
		return errors.New("Answers not found for kahoot: " + kg.pin)
	}
	if err != nil {
		return err
	}

	kg.answers = answers
	return nil
}

func (kg *KahootGame) GetCurrentAnswerIds() []int {
	var ids = make([]int, 0)

	for _, answer := range kg.answers {
		ids = append(ids, answer.ID)
	}

	return ids
}
