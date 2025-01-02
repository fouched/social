package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	UserID    int64      `json:"user_id"`
	Tags      string     `json:"tags"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Comments  *[]Comment `json:"comments"`
}

type PostsRepo struct {
	db *sql.DB
}

// CreatePost creates a post in the database
func (r *PostsRepo) CreatePost(post *Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		insert into posts (content, title, user_id, tags)
		values ($1, $2, $3, $4) returning id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		post.Content,
		post.Title,
		post.UserID,
		post.Tags,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePost creates a post in the database
func (r *PostsRepo) UpdatePost(post *Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		update posts 
		set content = $1, 
		    title = $2, 
		    tags = $3,
			updated_at = $4
		where id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		post.Content,
		post.Title,
		post.Tags,
		time.Now(),
		post.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeletePost creates a post in the database
func (r *PostsRepo) DeletePost(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		delete from posts 
		where id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// GetPostById retrieves a post from the database by ID
func (r *PostsRepo) GetPostById(id int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select id, user_id, title, content, tags, created_at, updated_at 
		from posts
		where id = $1
	`
	var post Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.Tags,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}
