package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/salvatoreolivieri/go-api/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewRedisStorage(redisDb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{redisDb: redisDb},
	}
}
