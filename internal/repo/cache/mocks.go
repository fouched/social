package cache

import "github.com/fouched/social/internal/repo"

func NewMockCache() Cache {
	return Cache{
		Users: &MockUserCache{},
	}
}

type MockUserCache struct{}

func (c *MockUserCache) Get(id int64) (*repo.User, error) {
	return nil, nil
}

func (c *MockUserCache) Set(user *repo.User) error {
	return nil
}

func (c *MockUserCache) Delete(id int64) {

}
