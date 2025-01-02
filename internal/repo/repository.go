package repo

import (
	"database/sql"
	"errors"
	"time"
)

const dbTimeout = time.Second * 60

var (
	ErrNotFound = errors.New("record not found")
)

type Repository struct {
	Posts interface {
		CreatePost(*Post) error
		GetPostById(int64) (*Post, error)
	}
	Comments interface {
		GetCommentsByPostId(int64) (*[]Comment, error)
	}
	Users interface {
		Create(*User) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Posts:    &PostsRepo{db},
		Comments: &CommentsRepo{db},
		Users:    &UsersRepo{db},
	}
}
