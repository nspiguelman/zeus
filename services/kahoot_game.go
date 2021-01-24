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
		TotalQuestions:  0,
		IsTimeout:       false,
		IsStarted:       false,
		IsScoreSent:     false,
	}
}

func NewKahootGame(rm *data.RepositoryManager) *KahootGame {
	once.Do(func() {
		once.Do(initGame)
		game.rm = rm
		game.IsStarted = true
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

	if err := kg.searchAnswers(); err != nil {
		return errors.New("An error occurred while getting answers: " + err.Error())
	}

	kg.CurrentQuestion = kg.questions[0].ID
	kg.IsStarted = true
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

func (kg *KahootGame) searchAnswers() error {
	answers, err := kg.rm.AnswerRepository.GetAllByKahootID(kg.kahoot.ID)
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
