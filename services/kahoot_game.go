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
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	game *KahootGame
	once sync.Once
)

type KahootGame struct {
	kahoot               *domain.Kahoot
	questions            []domain.Question
	answers              []domain.Answer
	host                 domain.Kahootuser
	users                []domain.Kahootuser
	pin                  string
	CurrentQuestion      int  // pregunta actual
	CurrentQuestionIndex int  // para iterar las questions
	TotalQuestions       int  // total de preguntas
	IsTimeout            bool // timeout de la pregunta
	IsStarted            bool // ya empezo el kahoot, no se puede suscribir nadie
	IsScoreSent          bool // la siguiente pregunta solo puede ser enviada cuando el score sea notificado
	rm                   *data.RepositoryManager
	ArrivalOrder         int // se setea en 0 en cada broadcast y se usa como un contador para tener el orden de llegada
	answerChannel        chan domain.AnswerMessage
}

func initGame() {
	game = &KahootGame{
		questions:       make([]domain.Question, 0),
		answers:         make([]domain.Answer, 0),
		host:            domain.Kahootuser{}, // TODO: Revisar si vamos a tener un host
		users:           make([]domain.Kahootuser, 0),
		pin:             "",
		CurrentQuestion: 0,
		CurrentQuestionIndex: 0,
		TotalQuestions:  0,
		IsTimeout:       false,
		IsStarted:       false,
		IsScoreSent:     false,
		ArrivalOrder: 0,
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

	kg.ArrivalOrder = 0
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

	kg.ArrivalOrder = 0
	kg.CurrentQuestionIndex++
	if err := kg.searchAnswers(kg.questions[kg.CurrentQuestionIndex].ID); err != nil {
		return errors.New("An error occurred while getting answers for question " + string(kg.CurrentQuestion) +  ": " + err.Error())
	}

	kg.CurrentQuestion = kg.questions[kg.CurrentQuestionIndex].ID
	return nil
}

func (kg *KahootGame) searchKahoot() error {
	kahoot, err := kg.rm.KahootRepository.GetByPin(kg.pin)
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

func (kg *KahootGame) setArrivalOrder() {
	totalUsers := kg.rm.UserRepository.CountByPin(kg.pin)
	kg.ArrivalOrder = totalUsers
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

func (kg *KahootGame) saveDbScore(token string, scoreInput *domain.ScoreMessage) error {
	score, err := kg.rm.KahootRepository.GetScore(token)
	if err != nil {
		return err
	}

	score.Score += scoreInput.Score
	score.IsCorrect = scoreInput.IsCorrect
	err = kg.rm.KahootRepository.SetScore(token, score)
	if err != nil {
		return err
	}
	return nil
}

// TODO: agregar el puntaje parametrizado
func (kg *KahootGame) calculateScore(answer domain.AnswerMessage) *domain.ScoreMessage {
	var score int
	isCorrect := kg.isAnswerCorrect(answer)

	if isCorrect {
		if score = 100; kg.ArrivalOrder > 10 {
			score = 1000 + kg.ArrivalOrder * 30
		}
		kg.ArrivalOrder++
		log.Println("respuesta: correcta - puntaje : ", score , " - respuesta id: " , answer.AnswerId)
	}
	return domain.NewScoreMessage(score, isCorrect)
}

func (kg *KahootGame) isAnswerCorrect(answer domain.AnswerMessage) bool {
	answers, err := kg.rm.AnswerRepository.CheckAnswer(answer.AnswerId)
	if err != nil {
		panic(err)
	}
	return answers.IsTrue
}

func (kg *KahootGame) BroadCastQuestion(m *melody.Melody, question domain.QuestionMessage) {
	msg, _ := json.Marshal(question)

	kg.setRound(5)
	err := m.Broadcast(msg)
	if err != nil {
		panic(err.Error())
	}

	go kg.processAnswers(m)
}

func (kg *KahootGame) Answer(answer domain.AnswerMessage) {
	if kg.IsTimeout {
		log.Println("timeout:", answer)
	} else {
		kg.answerChannel <- answer
	}
}

func (kg *KahootGame) setRound(timeout int) {
	// setea parametros de la maquina de estados
	// crea el canal de respuestas y timea el cierre
	kg.IsTimeout = false
	kg.IsScoreSent = false
	kg.ArrivalOrder = 100

	kg.answerChannel = make(chan domain.AnswerMessage, 1000)
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	go func(){
		<-timer.C
		kg.IsTimeout = true
		close(kg.answerChannel)
	}()
}

func (kg *KahootGame) processAnswers(m *melody.Melody) {
	log.Println("begin processing answers")

	for answer := range kg.answerChannel {
		var score = kg.calculateScore(answer)
		err := kg.saveDbScore(answer.Token, score)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	log.Println("end processing answers")

	go kg.sendScores(m)
}

func (kg *KahootGame) sendScores(m *melody.Melody) {
	log.Println("begin sending scores")
	var users []domain.Kahootuser
	users, err := kg.rm.UserRepository.GetAllUsers(kg.kahoot.ID)

	if err != nil {
		panic(err.Error())
	}

	var scoreMessages = make(map[string]domain.ScoreMessage)
	for _, user := range users {
		scoreMessage, err := kg.rm.KahootRepository.GetScore(user.Token)
		if err != nil {
			panic(err.Error())
		}

		scoreMessages[user.Token] = *scoreMessage
	}

	allScoreMessage := &domain.AllScoreMessage{
		ScoreMessages: scoreMessages,
		TypeMessage:   "score",
	}

	msg, _ := json.Marshal(allScoreMessage)
	err = m.Broadcast(msg)
	if err != nil {
		panic(err.Error())
	}

	kg.IsScoreSent = true
	log.Println("end sending scores")
}

func (kg *KahootGame) InitScore (token string) error {
	return kg.rm.KahootRepository.SetScore(token, domain.NewScoreMessage(0, false))
}

func (kg *KahootGame) CreateKahoot(kahootInput domain.KahootInput) (*domain.Kahoot, error) {
	kahootDomain := domain.NewKahoot(kahootInput)
	err := kg.rm.KahootRepository.Create(kahootDomain)
	return kahootDomain, err
}

func (kg *KahootGame) CreateQuestion(questionInput domain.QuestionInput, pin string) (*domain.Question, int, string) {
	kahoot, err := kg.rm.KahootRepository.GetByPin(pin)
	if kahoot == nil && err == nil {
		return nil, http.StatusNotFound, "Room " + pin + " not found"
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err.Error()
	}

	questionDomain := domain.NewQuestion(kahoot.ID, questionInput)
	if err := kg.rm.QuestionRepository.Create(questionDomain); err != nil {
		return nil, http.StatusInternalServerError, "Error saving question: " + err.Error()
	}
	return questionDomain, 200, ""
}

func (kg *KahootGame) CreateAnswers(answerInput []domain.AnswerInput, questionId int) ([]*domain.Answer, int, string) {
	question, err := kg.rm.QuestionRepository.GetById(questionId)
	fmt.Println(question)
	fmt.Println(err)
	if question == nil && err == nil {
		return nil, http.StatusNotFound, "Question " + strconv.Itoa(questionId) + " not found"
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err.Error()
	}

	trueFound := false
	var answersDomain []*domain.Answer
	for _, answer := range answerInput {
		answerDomain := domain.NewAnswer(questionId, answer)
		if answerDomain.IsTrue && !trueFound {
			trueFound = answerDomain.IsTrue
		} else if answerDomain.IsTrue && trueFound {
			return nil, http.StatusBadRequest, "Only one answer can be true"
		}
		answersDomain = append(answersDomain, answerDomain)
	}

	if !trueFound {
		return nil, http.StatusBadRequest, "Almost one answer must be true"
	}

	if err := kg.rm.AnswerRepository.CreateAnswers(answersDomain); err != nil {
		return nil, http.StatusInternalServerError, err.Error()
	}

	return answersDomain, http.StatusOK, "ok"
}