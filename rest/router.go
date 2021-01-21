package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nspiguelman/zeus/controllers"
)

type Server struct {
	router       *gin.Engine
	kahootController *controllers.KahootController
}

func NewServer(kahootController *controllers.KahootController) *Server {
	router := gin.Default()
	return &Server{router, kahootController}
}

func (s *Server) StartServer() {
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})

	})
	s.router.POST("/room", s.kahootController.CreateKahoot())
	s.router.POST("/room/:pin/question", s.kahootController.CreateQuestion())

	s.router.Run()
}
/*

s.router.GET("/room/:pin/users", s.getUsers)

// CREATE Game
s.router.POST("/room", controllers.KahootController.CreateKahoot) // Create game
s.router.POST("/room/:pin/question") // Create question
s.router.POST("/room/:pin/question/:id_question/answer") // Create answers

// Play game
s.router.POST("/room/:pin/login") // Login users
s.router.POST("/room/:pin/question/:id_question/user/:user_id") // Answer question
s.router.GET("/room/:pin/next_question") // Get next question
s.router.GET("/room/:pin/score") // Get total score
s.router.GET("/room/:pin/user/:user_id/score") // Get user score

func (s *Server) getUsers(c *gin.Context) {
	pin := c.Param("pin")
	c.JSON(http.StatusOK, s.kahootController.GetUsers(pin))
}

func (s *Server) createKahoot(c *gin.Context) {
	c.JSON(http.StatusCreated, s.kahootController.CreateNewKahootGame())
}
 */