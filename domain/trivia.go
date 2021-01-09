package domain

import "time"

type Trivia interface {
	GetQuestion() string
	GetId() int
	GetDate() *time.Time
}
