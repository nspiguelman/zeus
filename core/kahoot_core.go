package core

import (
	"github.com/nspiguelman/zeus/domain"
	"gorm.io/gorm"
)

type KahootInput struct {
	Name string `json:"name" validate:"min=1,max=50"`
}

func CreateKahootCore (db *gorm.DB, kahoot *KahootInput) (int, error) {
	return domain.CreateKahootDomain(db, &domain.Kahoot{ Name: kahoot.Name })
}