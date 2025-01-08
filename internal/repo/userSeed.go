package repo

import (
	"context"
	"database/sql"
)

type UserSeedRepo struct {
	db *sql.DB
}

func (r *UserSeedRepo) Create(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
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
