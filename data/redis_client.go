package data

import "github.com/gomodule/redigo/redis"

func getRedisConnection() (redis.Conn, error) {
	dns := "zeus_db_redis:6379"
	return redis.Dial("tcp", dns)
}