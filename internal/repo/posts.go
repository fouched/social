package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
	Version   int       `json:"version"`
}

type PostFeed struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostsRepo struct {
	db *sql.DB
}

// Create creates a post in the database
func (r *PostsRepo) Create(post *Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
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

// Update updates a post in the database
func (r *PostsRepo) Update(post *Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `
		update posts 
		set content = $1, 
		    title = $2, 
		    tags = $3,
			updated_at = $4,
			version = version + 1
		where id = $5 and version = $6
		returning version
	`

	err := r.db.QueryRowContext(ctx, query,
		post.Content,
		post.Title,
		post.Tags,
		time.Now(),
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

// Delete deletes a post from the database
func (r *PostsRepo) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
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

// GetById retrieves a post from the database by ID
func (r *PostsRepo) GetById(id int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `
		select p.id, p.user_id, p.title, p.content, p.tags, p.created_at, p.updated_at, p.version,
			u.username
		from posts p 
			join users u on u.id = p.user_id
		where p.id = $1
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
		&post.Version,
		&post.User.Username,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &post, ErrNotFound
		default:
			return &post, err
		}
	}
	return &post, nil
}

func (r *PostsRepo) GetUserFeed(userId int64, pq PaginatedQuery) ([]PostFeed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	//Note: Jetbrains incorrectly complains that created_at should be in the group by below
	query := `
		select 
			p.id, p.user_id, p.title, p.content, p.tags, p.version, p.created_at, p.updated_at,
		  	u.username,
		  	count (c.id) as comments_count  
		from posts p
		    left join comments c on c.post_id = p.id
		  	left join users u on u.id = p.user_id
		where p.user_id in (select user_id from followers where follower_id = $1)
		  	and (
				p.title ilike '%' || $2 || '%'
				OR p.content ilike '%' || $2 || '%'
		  	)
			and p.tags ilike '%' || $3 || '%'
			and p.created_at >= $4
		group by p.id, u.username
		order by p.created_at ` + pq.Sort + `
		limit $5 offset $6
	`
	fmt.Println(pq)
	//Note: above  placeholders for variables to be substituted can only be used
	//for actual data values and not for SQL syntax itself, hence the string concat

	rows, err := r.db.QueryContext(ctx, query, userId, pq.Search, pq.Tag, pq.Since, pq.Limit, pq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostFeed
	for rows.Next() {
		var pf PostFeed
		err := rows.Scan(
			&pf.ID,
			&pf.UserID,
			&pf.Title,
			&pf.Content,
			&pf.Tags,
			&pf.Version,
			&pf.CreatedAt,
			&pf.UpdatedAt,
			&pf.User.Username,
			&pf.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, pf)
	}

	return feed, nil
}
