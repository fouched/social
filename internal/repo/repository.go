package repo

import (
	"context"
	"database/sql"
)

type Repository struct {
	Posts interface {
		Create(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Posts: &PostsRepo{db},
		Users: &UsersRepo{db},
	}
}
