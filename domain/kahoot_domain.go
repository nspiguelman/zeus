package domain

import (
	"crypto/rand"
	"log"
	"math/big"
	"strconv"
)


type Kahoot struct {
	ID int `gorm:"primaryKey" json:"-"`
	PIN string `json:"pin"`
	Name string `gorm:"not null" json:"name"`
}

func NewKahoot(name string) *Kahoot {
	pin := strconv.FormatInt(generatePin(), 10)
	return &Kahoot{Name: name, PIN: pin}
}

func generatePin() int64 {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return n.Int64()
}