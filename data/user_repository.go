package data

import "github.com/nspiguelman/zeus/domain"

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