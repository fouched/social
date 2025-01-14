package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fouched/social/internal/repo"
	"github.com/redis/go-redis/v9"
	"time"
)

type UsersCache struct {
	rdb *redis.Client
}

const UserExpTime = time.Hour

func (c *UsersCache) Get(id int64) (*repo.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
	defer cancel()

	cacheKey := fmt.Sprintf("user-%d", id)
	data, err := c.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		// not in cache yet
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user repo.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (c *UsersCache) Set(user *repo.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
	defer cancel()

	cacheKey := fmt.Sprintf("user-%d", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.rdb.SetEx(ctx, cacheKey, json, UserExpTime).Err()
}
