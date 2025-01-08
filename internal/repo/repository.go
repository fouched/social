package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("record not found")
	DbTimeout   = time.Second * 60
)

type Repository struct {
	Posts interface {
		Create(*Post) error
		GetById(int64) (Post, error)
		Update(Post) error
		Delete(int64) error
		GetUserFeed(int64, PaginatedQuery) ([]PostFeed, error)
	}
	Comments interface {
		Create(*Comment) error
		GetByPostId(int64) ([]Comment, error)
	}
	Users interface {
		CreateAndInvite(user *User, token string, expiry time.Duration) error
		GetById(int64) (User, error)
		Activate(string) error
		Delete(int64) error
	}
	// special repo since we don't want to impl Saga pattern for seeds
	UserSeed interface {
		Create(*User) error
	}
	Followers interface {
		Follow(userID, followerID int64) error
		Unfollow(userID, followerID int64) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Posts:     &PostsRepo{db},
		Comments:  &CommentsRepo{db},
		Users:     &UsersRepo{db},
		UserSeed:  &UserSeedRepo{db},
		Followers: &FollowersRepo{db},
	}
}

func withTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
