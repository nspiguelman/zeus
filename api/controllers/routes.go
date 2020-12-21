package controllers

func (server *Server) initializeRoutes() {
	// el GET está como ejemplo de que puede hacerse
	server.Router.HandleFunc("/users", server.CreateUser).Methods("POST", "GET")
}
