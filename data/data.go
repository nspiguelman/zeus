package data

import (
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	data *Data
	once sync.Once
)

type Data struct {
	DB *gorm.DB
}

func initDB() {
	db, err := getDBConnection()
	if err != nil {
		log.Panic(err.Error())
	}

	data = &Data {
		DB: db,
	}
}

func New() *Data {
	once.Do(initDB)

	return data
}

func Close() error {
	if data == nil {
		return nil
	}

	sqlDB, _ := data.DB.DB()
	return sqlDB.Close()
}