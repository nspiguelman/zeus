package main

import (
	"github.com/nspiguelman/zeus/rest"
)

func main() {
	server := rest.NewServer()
	server.StartServer()
	// TODO: free resources if panic
}
