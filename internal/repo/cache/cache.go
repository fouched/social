package cache

import (
	"github.com/fouched/social/internal/repo"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Users interface {
		Get(int64) (*repo.User, error)
		Set(*repo.User) error
	}
}

func NewRedisCache(rdb *redis.Client) Cache {
	return Cache{
		Users: &UsersCache{rdb: rdb},
	}
}
