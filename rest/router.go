package rest

import (
	"github.com/nspiguelman/zeus/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
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
	// List all paths
	s.router.GET("/ping", s.ping)
	s.router.GET("/room/:pin/users", s.getUsers)

	go s.router.Run()
}

func (s *Server) getUsers(c *gin.Context) {
	pin := c.Param("pin")
	c.JSON(http.StatusOK, s.kahootController.GetUsers(pin))
}

func (s *Server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "message": "pong" })
}