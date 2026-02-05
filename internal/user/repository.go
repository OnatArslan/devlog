package user

import "github.com/OnatArslan/devlog/internal/sqlc"

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *userRepository {

	return &userRepository{
		q: q,
	}
}
