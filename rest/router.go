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


	//WEB SOCKET ; DONDE SE RECIBE LAS RESPUESTAS DE LOS CLIENTES
	s.router.GET("/room/:pin/ws", func(c *gin.Context) {
		s.socket.HandleRequest(c.Writer, c.Request)
		go proccesAnswer(s.socket)

	})

	//MANDA BROADCAST
	s.router.GET("/room/:pin/start", func(c *gin.Context) {

		go broadCastQuestion(s.socket)

	})


	s.router.Run()
}


func proccesAnswer(m *melody.Melody) {

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		//solo imprimo el mensaje que envia el cliente.
		//aca deberìa ir la logica de chequear la respuesta y sumar puntaje etc.
		fmt.Printf("%s\n", msg)
	})

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


