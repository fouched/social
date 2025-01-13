package cache

import (
	"github.com/fouched/social/internal/repo"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Users interface {
		Get(int64) (*repo.User, error)
		Set(*repo.User) error
	}
}

func NewRedisRepository(rdb *redis.Client) Repository {
	return Repository{
		Users: &UsersCache{rdb: rdb},
	}
}
