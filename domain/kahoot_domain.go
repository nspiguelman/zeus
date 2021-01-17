package domain

import "gorm.io/gorm"

type Kahoot struct {
	id  int `gorm:"primary_key;not null"`
	name string `gorm:"<-:create"`
}

func CreateKahootDomain(db *gorm.DB) {
	newKahoot := Kahoot{id: 20, name: "nahuel"}
	result := db.Debug().Create(&newKahoot)
	if result.Error != nil {
		panic(result.Error.Error())
	}
}
