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

type answer struct {
	questionId int
	answerId int
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
			pin := c.Param("pin")
			token := c.GetHeader("token")
			if s.kahootController.KahootGames.IsTimeout {
				// no envia respuesta
			}
			go proccesAnswer(msg, pin, token, s.channelkahoot)
			s.answerToProcess = append(s.answerToProcess,<-s.channelkahoot)
		})
	})

	//CALCULAR PUNTAJES
	s.router.GET("/room/:pin/calculateScores", func(c *gin.Context) {
		go calculateScores(s.answerToProcess)
	})

	//MANDA BROADCAST
	s.router.GET("/room/:pin/send_question", func(c *gin.Context) {
		// manejar el timeout con una nueva go routine
		if !s.kahootController.KahootGames.IsScoreSent {
			log.Panic("no fue enviado")
		}
		// setear timeout en false
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

/**
Procesa un mensaje y lo encola en un channel
 */
func proccesAnswer(msg []byte, pin string, token string, c chan answer) {
	// questionId
	// answerId
	/**
	{ questionId: 1, answerIds: [200, 201, 202, 203] }
	 */

	/*var message = answer{
		kahootId: pin,
		timeout:           true,
		token : "ASDASDASDAD",
		questionId:         1,
		answerId : "A",
	}*/
	//	c <- message
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


