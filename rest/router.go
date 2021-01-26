package rest

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nspiguelman/zeus/controllers"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net/http"
	"time"
)

type Server struct {
	router           *gin.Engine
	kahootController *controllers.KahootController
	socket           *melody.Melody
	channelKahoot    chan Answer
	answerToProcess  []Answer
}

type Score struct {
	PartialScore int    `json:"partialScore"`
	IsCorrect    bool   `json:"isCorrect"`
	TypeMessage  string `json:"typeMessage"`
}

type Question struct {
	QuestionId  int    `json:"questionId"`
	AnswerIds   []int  `json:"answerIds"`
	TypeMessage string `json:"typeMessage"`
}

type Answer struct {
	QuestionId int    `json:"questionId"`
	AnswerId   int    `json:"answerId"`
	Token      string `json:"token"`
	IsTimeout  bool   `json:"isTimeout"`
}

func NewServer(kahootController *controllers.KahootController) *Server {
	router := gin.Default()
	socket := melody.New()
	answerToProcess := []Answer{}
	channelKahoot := make(chan Answer)
	return &Server{router, kahootController, socket, channelKahoot, answerToProcess}
}

func (s *Server) StartServer() {
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})

	})

	/*END POINTS*/
	s.router.POST("/room", s.kahootController.CreateKahoot())
	s.router.POST("/room/:pin/question", s.kahootController.CreateQuestion())

	//LOGIN ; Devuelve un token al usuario.
	s.router.POST("/room/:pin/name/:name/login", func(c *gin.Context) {
		name := c.Param("name")
		var token = GenerateToken(name)
		c.JSON(200, gin.H{
			"token": token,
		})
	}) // Login users


	//WEB SOCKET ;
	s.router.GET("/room/:pin/ws", func(c *gin.Context)  {
		s.socket.HandleRequest(c.Writer, c.Request)
	})

	s.socket.HandleMessage(func(x *melody.Session, msg []byte) {
		//pin := x.Request.Header.Get("pin")
		token := x.Request.Header.Get("token")

		answer := Answer{}
		answer.Token = token

		err := json.Unmarshal([]byte(msg), &answer)
		if err != nil {
			panic(err.Error())
		}
		answer.Token = token

		if s.kahootController.KahootGames.IsTimeout {
			answer.IsTimeout = true
		}
		processAnswer(answer)

	})

	//MANDA BROADCAST
	s.router.GET("/room/:pin/send_question", func(c *gin.Context) {
		// manejar el timeout con una nueva go routine
		//if !s.kahootController.KahootGames.IsScoreSent {
		//	log.Panic("no fue enviado")
		//}
		pin := c.Param("pin")
		if !s.kahootController.KahootGames.IsStarted {
			if err := s.kahootController.KahootGames.Start(pin); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			if !s.kahootController.KahootGames.IsScoreSent {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send next question. Score must be sent."})
				return
			}

			if err := s.kahootController.KahootGames.NextQuestion(pin); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// setear timeout en false
		go broadCastQuestion(s.socket, Question{
			QuestionId: s.kahootController.KahootGames.CurrentQuestion,
			AnswerIds: s.kahootController.KahootGames.GetCurrentAnswerIds(),
			TypeMessage: "question",
		})
	})

	s.router.Run()
}



func processAnswer(answer Answer) {
	fmt.Println("processAnswer: Procesando respuesta:")
	fmt.Printf("%+v\n", answer)

	var score = calculateScore(answer)
	err :=saveDbScore(score,answer.Token)
	if err != nil {
		panic(err.Error())
	}
}

func saveDbScore(score int, user string) error{
	//aca guardamos los datos . Redis
	return nil
}

func calculateScore(answer Answer ) int{
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


func broadCastScore(m *melody.Melody){
	time.Sleep(2 * time.Second)
	b := []byte("estos son los puntajes")
	m.Broadcast(b)

}


func broadCastQuestion(m *melody.Melody, question Question) {
	time.Sleep(2 * time.Second)
	msg, _ := json.Marshal(question)
	m.Broadcast(msg)
}

func GenerateToken(name string) string {
	// pasar a jwt
	hash, err := bcrypt.GenerateFromPassword([]byte(name), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}
