package post

import "github.com/OnatArslan/devlog/internal/sqlc"

type PostRepository struct {
	q *sqlc.Queries
}

func NewPostRepository(q *sqlc.Queries) *PostRepository {
	return &PostRepository{
		q: q,
	}
}
