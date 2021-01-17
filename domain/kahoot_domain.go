package domain

import (
	"gorm.io/gorm"
)

type Kahoot struct {
	ID int `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func CreateKahootDomain(db *gorm.DB, kahoot *Kahoot) (int, error) {
	result := db.Create(kahoot)
	if result.Error != nil {
		return 0, result.Error
	}
	return kahoot.ID, nil
}