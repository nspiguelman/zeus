package rest

import (
	"crypto/md5"
	"encoding/hex"
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
	channelkahoot chan answer
	answerToProcess []answer
}

type answer struct {
	kahootId string
	timeout bool
	token string
	IdQuestion int
	answer string
}

func NewServer(kahootController *controllers.KahootController) *Server {
	router := gin.Default()
	socket := melody.New()
	answerToProcess := []answer{}
	channelkahoot := make(chan answer)
	return &Server{router, kahootController, socket, channelkahoot, answerToProcess}
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


	//WEB SOCKET ; DONDE SE RECIBE LAS RESPUESTAS DE LOS CLIENTES
	s.router.GET("/room/:pin/ws", func(c *gin.Context) {

		s.socket.HandleRequest(c.Writer, c.Request)
		s.socket.HandleMessage(func(x *melody.Session, msg []byte) {
			go proccesAnswer(msg, s.channelkahoot)
			s.answerToProcess = append(s.answerToProcess,<-s.channelkahoot)
		})

	})

	//CALCULAR PUNTAJES
	s.router.GET("/room/:pin/calculateScores", func(c *gin.Context) {
		go calculateScores(s.answerToProcess)
	})

	//MANDA BROADCAST
	s.router.GET("/room/:pin/start", func(c *gin.Context) {
		go broadCastQuestion(s.socket)
	})


	s.router.Run()
}

func calculateScores(answers []answer)()  {
	//calcula puntajes desencolando la lista FILO y los serializa.
	//y lanzar broadcast con pùntajes.

	//var lastAnswer = answers[len(answers)-1] //guardo el ultimo
	//answers[len(answers)-1] = "" // Erase element (write zero value)
	//answers = answers[:len(answers)-1] // elimino el ultimo


	fmt.Printf("%v", answers)
}


func proccesAnswer(msg []byte, c chan answer){
	//convertir msg en struct y mandarlo al channel
	var message = answer{kahootId: "PING",
		timeout:           true,
		token : "ASDASDASDAD",
		IdQuestion:         1,
		answer : "A",
	}
	c <- message
}


func broadCastQuestion(m *melody.Melody){
	time.Sleep(2 * time.Second)
	b := []byte("{question: 1 = 1 ?}")
	m.Broadcast(b)

}

func GenerateToken(name string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(name), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash to store:", string(hash))

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}


