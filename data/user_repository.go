package data

import "github.com/nspiguelman/zeus/domain"

type UserRepository struct {
	Data *Data
}

func (usr *UserRepository) Create(kahootUser *domain.Kahootuser) error {
	result := usr.Data.DB.Create(kahootUser)
	return result.Error
}