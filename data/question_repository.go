package data

import (
	"errors"
	"github.com/nspiguelman/zeus/domain"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	Data *Data
}

func (qr *QuestionRepository) Create(question *domain.Question) error {
	result := qr.Data.DB.Create(question)
	return result.Error
}

func (qr *QuestionRepository) GetAllByKahootID(kahootID int) ([]domain.Question, error) {
	var questions []domain.Question
	result := qr.Data.DB.Where("kahoot_id = ?", kahootID).Order("id asc").Find(&questions)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []domain.Question{}, nil
	}

	if result.Error != nil {
		return []domain.Question{}, result.Error
	}

	return questions, nil
}

func (qr *QuestionRepository) GetById(questionId int) (*domain.Question, error) {
	var question domain.Question
	result := qr.Data.DB.Where("id = ?", questionId).First(&question)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}