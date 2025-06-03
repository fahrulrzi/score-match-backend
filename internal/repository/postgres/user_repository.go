package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fahrulrzi/score-match-backend/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUserByEmail implements repository.UserRepository.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, password, username, created_at, updated_at FROM users WHERE email = $1`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return user, err
}

// GetUserById implements repository.UserRepository.
func (r *UserRepository) GetUserById(ctx context.Context, id uint) (*entity.User, error) {
	query := `SELECT id, email, password, username, created_at, updated_at FROM users WHERE id = $1`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
	INSERT INTO users (email, password, username, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Username, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	return err
}
