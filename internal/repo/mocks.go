package repo

import "time"

func NewMockRepo() Repository {
	return Repository{
		Users: &MockUserRepo{},
	}
}

type MockUserRepo struct{}

func (r *MockUserRepo) CreateAndInvite(user *User, token string, expiry time.Duration) error {
	return nil
}

func (r *MockUserRepo) GetById(id int64) (*User, error) {
	return &User{}, nil
}

func (r *MockUserRepo) GetByEmail(email string) (*User, error) {
	return &User{}, nil
}

func (r *MockUserRepo) Activate(token string) error {
	return nil
}

func (r *MockUserRepo) Delete(id int64) error {
	return nil
}
