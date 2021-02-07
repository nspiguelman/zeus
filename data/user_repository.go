package data

import (
	"errors"
	"github.com/nspiguelman/zeus/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	Data *Data
}

func (usr *UserRepository) Create(kahootUser *domain.Kahootuser) error {
	result := usr.Data.DB.Create(kahootUser)
	return result.Error
}

func (usr *UserRepository) CountByPin(pin string) int {
	var result int64
	usr.Data.DB.Where("kahoot_id = ?", pin).Count(&result)
	return int(result)
}

func (usr *UserRepository) GetAllUsers(kahootId int) ([]domain.Kahootuser, error) {
	var users []domain.Kahootuser
	result := usr.Data.DB.Where("kahoot_id = ?", kahootId).Find(&users)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []domain.Kahootuser{}, nil
	}

	if result.Error != nil {
		return []domain.Kahootuser{}, result.Error
	}

	return users, nil
}