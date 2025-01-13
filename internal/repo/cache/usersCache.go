package cache

import (
	"github.com/fouched/social/internal/repo"
	"github.com/redis/go-redis/v9"
)

type UsersCache struct {
	rdb *redis.Client
}

func (c *UsersCache) Get(id int64) (*repo.User, error) {
	//TODO
	return nil, nil
}

func (c *UsersCache) Set(user *repo.User) error {
	//TODO
	return nil
}
