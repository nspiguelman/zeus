package data

import "github.com/nspiguelman/zeus/domain"

type KahootRepository struct {
	Data *Data
}

func (kr *KahootRepository) Create(kahoot *domain.Kahoot) (*domain.Kahoot, error) {
	result := kr.Data.DB.Create(kahoot)
	if result.Error != nil {
		return nil, result.Error
	}
	return kahoot, nil
}