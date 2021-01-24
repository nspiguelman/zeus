package data

import (
	"errors"
	"github.com/nspiguelman/zeus/domain"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	Data *Data
}

func (ar AnswerRepository) Create(answer *domain.Answer) error {
	result := ar.Data.DB.Create(answer)
	return result.Error
}

func (ar AnswerRepository) GetAllByKahootID(kahootID int) ([]domain.Answer, error) {
	var answers []domain.Answer
	result := ar.Data.DB.Where("kahoot_id = ?", kahootID).Order("question_id asc, id asc").Find(&answers)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []domain.Answer{}, nil
	}

	if result.Error != nil {
		return []domain.Answer{}, result.Error
	}

	return answers, nil
}


