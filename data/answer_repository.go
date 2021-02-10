package data

import (
	"errors"
	"github.com/nspiguelman/zeus/domain"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	Data *Data
}

func (ar AnswerRepository) CreateAnswers(answers []*domain.Answer) error {
	result := ar.Data.DB.CreateInBatches(answers, len(answers))
	return result.Error
}

func (ar AnswerRepository) GetAllByQuestionID(questionID int) ([]domain.Answer, error) {
	var answers []domain.Answer
	result := ar.Data.DB.Where("question_id = ?", questionID).Order("id asc").Find(&answers)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []domain.Answer{}, nil
	}

	if result.Error != nil {
		return []domain.Answer{}, result.Error
	}

	return answers, nil
}

func (ar AnswerRepository) CheckAnswer(answer_id int) (domain.Answer, error) {
	var answers domain.Answer
	result := ar.Data.DB.Where("id = ?", answer_id).First(&answers)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.Answer{}, nil
	}

	if result.Error != nil {
		return domain.Answer{}, result.Error
	}

	return answers, nil
}

