package repo

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user"`
}

type CommentsRepo struct {
	db *sql.DB
}

func (r *CommentsRepo) GetByPostId(id int64) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `
		select c.id, c.post_id, c.content, c.created_at, c.updated_at, 
		       u.id, u.username, u.email, u.created_at, u.updated_at 
		from comments c 
		join users u on u.id = c.user_id
		where c.post_id = $1
		order by c.created_at desc
	`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		var u User
		c.User = &u
		err := rows.Scan(&c.ID, &c.PostID, &c.Content, &c.CreatedAt, &c.UpdatedAt,
			&c.User.ID, &c.User.Username, &c.User.Email, &c.User.CreatedAt, &c.User.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentsRepo) Create(comment *Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `
		insert into comments (post_id, user_id, content)
		values ($1, $2, $3) returning id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
