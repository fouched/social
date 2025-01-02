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
	User      User      `json:"user"`
}

type CommentsRepo struct {
	db *sql.DB
}

func (r *CommentsRepo) GetCommentsByPostId(id int64) (*[]Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, 
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
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt,
			&c.User.ID, &c.User.Username, &c.User.Email, &c.User.CreatedAt, &c.User.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return &comments, nil
}
