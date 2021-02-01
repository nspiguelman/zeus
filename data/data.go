package data

import (
	"github.com/gomodule/redigo/redis"
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
	RedisDB redis.Conn
}

func initDB() {
	db, err := getDBConnection()
	if err != nil {
		log.Panic(err.Error())
	}

	conn, err := getRedisConnection()
	if err != nil {
		log.Panic(err.Error())
	}

	data = &Data {
		DB: db,
		RedisDB: conn,
	}
}

func New() *Data {
	once.Do(initDB)
	return data
}

func Close() (error, error) {
	if data == nil {
		return nil, nil
	}

	sqlDB, _ := data.DB.DB()
	return sqlDB.Close(), data.RedisDB.Close()
}