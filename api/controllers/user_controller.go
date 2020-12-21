package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/nspiguelman/zeus/api/responses"
)

// CreateUser crea un usuario
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	response := &responses.EmptyResponse{Success: true, Message: ""}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
