package store

import (
	"context"
	"database/sql"
	"time"
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (s *UserStore) Create(ctx context.Context, user *User) error {

	query := `
		insert into users(username, password, email) 
		values($1, $2, $3) returning id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		&user.Username,
		&user.Password,
		&user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
