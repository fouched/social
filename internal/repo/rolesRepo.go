package repo

import (
	"context"
	"database/sql"
	"errors"
)

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RolesRepo struct {
	db *sql.DB
}

func (r *RolesRepo) GetByName(roleName string) (*Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := "select id, name, level, description from roles where name = $1"

	var role Role
	err := r.db.QueryRowContext(ctx, query, roleName).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &role, ErrNotFound
		default:
			return &role, err
		}
	}

	return &role, nil
}
