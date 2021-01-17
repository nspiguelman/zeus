package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nspiguelman/zeus/core"
	"github.com/nspiguelman/zeus/domain"
	"gorm.io/gorm"
	"net/http"
)

type KahootController struct {
	kahootGames []domain.KahootGame
}

func NewKahootController() KahootController {
	kahootGames := make([]domain.KahootGame, 0)
	return KahootController{kahootGames}
}

/*
func (kc *KahootController) CreateNewKahootGame() domain.KahootGame {
	game := domain.NewKahootGame()
	kc.kahootGames = append(kc.kahootGames, game)
	return game
}

func (kc *KahootController) _getUsers(pin string) []domain.User {
	for _, game := range kc.kahootGames {
		if game.GetPin() == pin {
			return game.GetUsers()
		}
	}
	panic("Game room not found")
}
*/

// ---------------

func (kc *KahootController) CreateKahoot(db *gorm.DB) func(c *gin.Context) {
	return func (c *gin.Context) {
		// validations
		// call to domain
		core.CreateKahootCore(db)
		// response
		c.JSON(http.StatusOK, gin.H{"data": ""})
	}
}