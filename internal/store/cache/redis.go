package cache

import "github.com/go-redis/redis/v8"

func NewRedisCLient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
