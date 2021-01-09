package controllers

import "github.com/nspiguelman/zeus/domain"

type KahootController struct {
	kahootGames []domain.KahootGame
}

func NewKahootController() KahootController {
	kahootGames := make([]domain.KahootGame, 0)
	return KahootController{kahootGames}
}

func (kc *KahootController) CreateNewKahootGame() domain.KahootGame {
	game := domain.NewKahootGame()
	kc.kahootGames = append(kc.kahootGames, game)
	return game
}

func (kc *KahootController) GetUsers(pin string) []domain.User {
	for _, game := range kc.kahootGames {
		if game.GetPin() == pin {
			return game.GetUsers()
		}
	}
	panic("Game room not found")
}