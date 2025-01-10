package repo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var ErrDuplicateUser = errors.New("a user with that email of username already exits")

type UsersRepo struct {
	db *sql.DB
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (r *UsersRepo) GetById(id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `
		select id, email, username, password, created_at, updated_at 
		from users
		where id = $1
		and is_active = true
	`
	user := User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &user, ErrNotFound
		default:
			return &user, err
		}
	}
	return &user, nil
}

func (r *UsersRepo) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `
		select id, email, username, password, created_at, updated_at 
		from users
		where email = $1
		and is_active = true
	`
	user := User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &user, ErrNotFound
		default:
			return &user, err
		}
	}
	return &user, nil
}

func (r *UsersRepo) CreateAndInvite(user *User, token string, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	return withTx(ctx, r.db, func(tx *sql.Tx) error {
		if err := r.create(user, tx); err != nil {
			return err
		}

		if err := r.createUserInvitation(ctx, tx, user.ID, token, expiry); err != nil {
			return err
		}

		return nil
	})
}

func (r *UsersRepo) create(user *User, tx *sql.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `
		insert into users(username, password, email) 
		values($1, $2, $3) returning id, created_at, updated_at
	`

	err := tx.QueryRowContext(ctx, query,
		&user.Username,
		&user.Password.hash,
		&user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrDuplicateUser
		} else {
			return err
		}
	}

	return nil
}

func (r *UsersRepo) createUserInvitation(ctx context.Context, tx *sql.Tx, userID int64, token string, expiry time.Duration) error {
	query := `insert into user_invitations (user_id, token, expiry) values ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, DbTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID, []byte(token), time.Now().Add(expiry))
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepo) Activate(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	return withTx(ctx, r.db, func(tx *sql.Tx) error {
		// find the user for the token
		user, err := r.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		// update the user
		user.IsActive = true
		if err := r.update(ctx, tx, user); err != nil {
			return err
		}

		// delete the invitation
		if err := r.deleteUserInvitationsById(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})

}

func (r *UsersRepo) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		select u.id, u.username, u.email, u.is_active, u.created_at, u.updated_at
		from users u
			join user_invitations ui on ui.user_id = u.id
		where ui.token = $1
			and ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (r *UsersRepo) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := "update users set username = $1, email = $2, is_active = $3, updated_at = $4 where id = $5"

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, time.Now(), user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepo) deleteUserInvitationsById(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := "delete from user_invitations where user_id = $1"

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepo) Delete(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	return withTx(ctx, r.db, func(tx *sql.Tx) error {
		if err := r.deleteUserInvitationsById(ctx, tx, userID); err != nil {
			return err
		}

		if err := r.deleteUserById(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (r *UsersRepo) deleteUserById(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := "delete from users where id = $1"

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
