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
	"time"
)

type Server struct {
	router       *gin.Engine
	kahootController *controllers.KahootController
	socket *melody.Melody
}

type Score struct {
	partialScore int
	isCorrect bool
	typeMessage string `default:"score"`
}

type Question struct {
	QuestionId int
	answerIds []int
	typeMessage string `default:"question"`
}

type Answer struct {
	QuestionId int `json:"questionId"`
	AnswerId int `json:"answerId"`
	Token string `json:"token,omitempty"`
	IsTimeout bool `json:"token,omitempty"`
}

func NewServer(kahootController *controllers.KahootController) *Server {
	router := gin.Default()
	socket := melody.New()
	return &Server{router, kahootController, socket}
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
		// setear timeout en false
		go broadCastQuestion(s.socket)
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



func broadCastQuestion(m *melody.Melody){
	time.Sleep(2 * time.Second)
	b := []byte("{\"typeMessage\": \"question\", \"questionId\": 23, \"answerIds\": [123, 234, 345]}")
	m.Broadcast(b)

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


