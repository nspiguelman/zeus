package api

import (
	"github.com/nspiguelman/zeus/api/controllers"
)

// Run a
func Run() {
	var server = controllers.Server{}
	server.Initialize()
	server.Run(":8080")
}
