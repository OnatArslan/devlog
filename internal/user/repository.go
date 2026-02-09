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

// userRepository provides persistence operations backed by sqlc queries.
type userRepository struct {
	q *sqlc.Queries
}

// NewUserRepository wires repository methods to generated sqlc query implementations.
func NewUserRepository(q *sqlc.Queries) *userRepository {
	// Return a repository instance backed by generated sqlc queries.
	return &userRepository{
		q: q,
	}
}

// CreateUserParams defines the input fields required to insert a new user row.
type CreateUserParams struct {
	Email        string
	Username     string
	PasswordHash string
}

// CreateUser inserts a user and maps database constraint errors to domain errors.
func (r *userRepository) CreateUser(ctx context.Context, input CreateUserParams) (User, error) {
	// Execute the insert query and capture the created row.
	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: input.PasswordHash,
	})

	if err != nil {
		// Translate known Postgres unique violations into domain-level conflicts.
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

	// Map sqlc row fields into the package domain model.
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

// GetByEmail returns an active user by email or a domain not-found error.
func (r *userRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	// Query one active user by email from the database.
	row, err := r.q.GetByEmail(ctx, email)
	if err != nil {
		// Convert no-row result into a domain not-found error.
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("repository get by email: %w", err)
	}

	// Map query result into the package domain model.
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
