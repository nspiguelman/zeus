package data

import "github.com/nspiguelman/zeus/domain"

type QuestionRepository struct {
	Data *Data
}

func (qr *QuestionRepository) Create(question *domain.Question) error {
	result := qr.Data.DB.Create(question)
	return result.Error
}