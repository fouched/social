package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersRepo struct {
	db *sql.DB
}

func (r *FollowersRepo) Follow(userID, followerID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := "insert into followers (user_id, follower_id) values ($1, $2)"

	_, err := r.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			// ignore - the user is already following
		} else {
			return err
		}
	}

	return nil
}

func (r *FollowersRepo) Unfollow(userID, followerID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := "delete from followers where user_id = $1 and follower_id = $2"

	_, err := r.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
