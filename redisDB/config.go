package redisDB

import "github.com/go-redis/redis/v8"

var Rds *redis.Client

func init() {
	// initialize redis
	Rds = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "",
		DB:       0,
	})
}
