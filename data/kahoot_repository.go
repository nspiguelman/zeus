package data

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/nspiguelman/zeus/domain"
	"gorm.io/gorm"
)

type KahootRepository struct {
	Data *Data
}

func (kr *KahootRepository) Create(kahoot *domain.Kahoot) error {
	result := kr.Data.DB.Create(kahoot)
	return result.Error
}

func (kr *KahootRepository) GetByPin(pin string) (*domain.Kahoot, error) {
	var kahoot domain.Kahoot
	result := kr.Data.DB.Where("pin = ?", pin).First(&kahoot)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &kahoot, nil
}

func (kr *KahootRepository) SetScore(token string, score *domain.ScoreMessage) error {
	_, err := kr.Data.RedisPool.Get().Do("HMSET", redis.Args{token}.AddFlat(score)...)
	return err
}

func (kr *KahootRepository) GetScore (token string) (*domain.ScoreMessage, error){
	values, err := redis.Values(kr.Data.RedisPool.Get().Do("HGETALL", token))
	if err != nil {
		return nil, err
	}

	var score domain.ScoreMessage
	err = redis.ScanStruct(values, &score)
	if err != nil {
		return nil, err
	}

	return &score, nil
}