package core

import (
	"github.com/nspiguelman/zeus/data"
	"github.com/nspiguelman/zeus/domain"
)

type KahootCore struct {
	KahootRepository *data.KahootRepository
	UserRepository *data.UserRepository
}

type KahootInput struct {
	Name string `json:"name" validate:"min=1,max=50"`
}


type UserInput struct {
	Username string `json:"username" validate:"min=4, max=15"`
	//PIN
}

func (kc *KahootCore) CreateKahootCore (kahootInput *KahootInput) (*domain.Kahoot, error) {
	kahoot := domain.NewKahoot(kahootInput.Name)
	return kc.KahootRepository.Create(kahoot)
}