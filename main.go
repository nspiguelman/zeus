package main

import (
	"github.com/nspiguelman/zeus/controllers"
	"github.com/nspiguelman/zeus/rest"
)

func main() {
	kahootController := controllers.NewKahootController()

	server := rest.NewServer(&kahootController)
	server.StartServer()
}
