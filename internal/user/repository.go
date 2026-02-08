package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/OnatArslan/devlog/internal/sqlc"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *userRepository {
	return &userRepository{
		q: q,
	}
}

type CreateUserParams struct {
	Email        string
	Username     string
	PasswordHash string
}

func (r *userRepository) CreateUser(ctx context.Context, input CreateUserParams) (User, error) {

	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: input.PasswordHash,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return User{}, ErrEmailTaken
			case "users_username_key":
				return User{}, ErrUsernameTaken
			default:
				return User{}, ErrConflict
			}
		}
		return User{}, err
	}

	return User{
		ID:                 row.ID,
		Email:              row.Email,
		Username:           row.Username,
		PasswordHash:       row.PasswordHash,
		IsActive:           row.IsActive,
		TokenInvalidBefore: row.TokenInvalidBefore.Time,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	row, err := r.q.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("repo get user by email : %w", err)
	}

	return User{
		ID:                 row.ID,
		Email:              row.Email,
		Username:           row.Username,
		PasswordHash:       row.PasswordHash,
		IsActive:           row.IsActive,
		TokenInvalidBefore: row.TokenInvalidBefore.Time,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
	}, nil
}
