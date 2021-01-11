package main

import (
	"github.com/nspiguelman/zeus/controllers"
	"github.com/nspiguelman/zeus/rest"
)

func main() {
	/*var user *domain.User
	var nickname string
	var password string*/
	kahootController := controllers.NewKahootController()

	server := rest.NewServer(&kahootController)
	server.StartServer()
/*
	shell := ishell.New()
	shell.SetPrompt("Kahoot >> ")
	shell.Print("Type 'help' to get commands\n")

	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "login to Kahoot",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			c.Print("Nickname: ")
			nickname = c.ReadLine()
			for nickname == "" {
				c.Println("Nickname is empty, please login with valid user")
				c.Print("Nickname: ")
				nickname = c.ReadLine()
			}

			c.Print("Password: ")
			password = c.ReadLine()

			user = domain.NewUser(nickname, password)

			return
		},
	})

	shell.Run()
*/
}
