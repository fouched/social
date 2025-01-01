package repo

import (
	"context"
	"database/sql"
	"time"
)

type UsersRepo struct {
	db *sql.DB
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *UsersRepo) Create(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		insert into users(username, password, email) 
		values($1, $2, $3) returning id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
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
