package data

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func newPool() *redis.Pool {
	dns := "zeus_db_redis:6379"
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", dns)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}