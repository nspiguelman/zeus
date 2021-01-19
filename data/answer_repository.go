package data

import "github.com/nspiguelman/zeus/domain"

type AnswerRepository struct {
	Data *Data
}

func (ar AnswerRepository) Create(answer *domain.Answer) error {
	result := ar.Data.DB.Create(answer)
	return result.Error
}


