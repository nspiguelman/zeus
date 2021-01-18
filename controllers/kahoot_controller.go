package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nspiguelman/zeus/core"
	"github.com/nspiguelman/zeus/data"
	"github.com/nspiguelman/zeus/domain"
	"gopkg.in/validator.v2"
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

func (kc *KahootController) CreateKahoot() func(c *gin.Context) {
	return func (c *gin.Context) {
		var kahoot core.KahootInput
		if err := c.ShouldBindJSON(&kahoot); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
		}

		if err := validator.Validate(kahoot); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
			return
		}

		kahootCore := &core.KahootCore{
			KahootRepository: &data.KahootRepository{
				Data: data.New(),
			},
			UserRepository: &data.UserRepository{
				Data: data.New(),
			},
		}
		result, err := kahootCore.CreateKahootCore(&kahoot);
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func (kc *KahootController) Login() func(c *gin.Context) {
	return func (c *gin.Context) {

	}
}