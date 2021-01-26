package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nspiguelman/zeus/data"
	"github.com/nspiguelman/zeus/domain"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/olahol/melody.v1"
	"log"
	"sync"
	"time"
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
	ArrivalOrder int // se setea en 0 en cada broadcast y se usa como un contador para tener el orden de llegada
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

func (kg *KahootGame) GenerateToken(name string) string {
	// pasar a jwt
	hash, err := bcrypt.GenerateFromPassword([]byte(name), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}


func (kg *KahootGame) ProcessAnswer(answer domain.AnswerMessage) error  {
	fmt.Println("processAnswer: Procesando respuesta:")
	fmt.Printf("%+v\n", answer)

	var score = calculateScore(answer)
	err :=saveDbScore(score,answer.Token)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

func saveDbScore(score int, user string) error{
	//aca guardamos los datos . Redis
	return nil
}

func calculateScore(answer domain.AnswerMessage ) int{
	//solo chequeo rta correcta, pero no estoy contemplando por ahora el orden de llegada.
	var score = 0
	if ( isAnswerCorrect(answer.QuestionId,answer.AnswerId) ){
		score += 10
	}
	return score
}

func isAnswerCorrect(questionId int,answerId int ) bool{
	//chequear en base rta correcta
	return true
}

func (kg *KahootGame) BroadCastQuestion(m *melody.Melody, question domain.QuestionMessage) {
	time.Sleep(2 * time.Second)
	msg, _ := json.Marshal(question)
	m.Broadcast(msg)
}

