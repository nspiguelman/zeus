package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nspiguelman/zeus/controllers"
	"gopkg.in/olahol/melody.v1"
)

type Server struct {
	router           *gin.Engine
	kahootController *controllers.KahootController
	socket           *melody.Melody
}

func NewServer() *Server {
	router := gin.Default()
	socket := melody.New()
	kahootController := controllers.NewKahootController()
	return &Server{router, &kahootController, socket }
}

func (s *Server) StartServer() {
	s.router.GET("/ping", s.kahootController.Ping())

	/*END POINTS*/
	//ABM PREGUNTAS KAHOOT
	s.router.POST("/room", s.kahootController.CreateKahoot())
	s.router.POST("/room/:pin/question", s.kahootController.CreateQuestion())

	//LOGIN ; Devuelve un token al usuario.
	s.router.POST("/room/:pin/name/:name/login", s.kahootController.Login())

	//WEB SOCKET ;
	s.router.GET("/room/:pin/ws", s.kahootController.HandShake(s.socket))
	s.kahootController.HandleMessage(s.socket)
	s.router.GET("/room/:pin/send_question", s.kahootController.SendQuestion(s.socket))

	s.router.Run()
}