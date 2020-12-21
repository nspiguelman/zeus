package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server es usado para exportar el Router
type Server struct {
	Router *mux.Router
}

// Initialize carga las rutas de la API
func (server *Server) Initialize() {
	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

// Run inicia la aplicación en el puerto pasado como parámetro
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
