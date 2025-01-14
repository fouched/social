package cache

import (
	"github.com/fouched/social/internal/repo"
	"github.com/redis/go-redis/v9"
	"time"
)

var cacheTimeout = time.Second * 3

type Cache struct {
	Users interface {
		Get(int64) (*repo.User, error)
		Set(*repo.User) error
		Delete(id int64)
	}
}

func NewRedisCache(rdb *redis.Client) Cache {
	return Cache{
		Users: &UsersCache{rdb: rdb},
	}
}
