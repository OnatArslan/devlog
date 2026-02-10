package post

import (
	"context"
	"fmt"

	"github.com/OnatArslan/devlog/internal/sqlc"
)

type PostRepository struct {
	q *sqlc.Queries
}

func NewPostRepository(q *sqlc.Queries) *PostRepository {
	return &PostRepository{
		q: q,
	}
}

type CreatePostParams struct {
	AuthorID int64
	Title    string
	Content  string
}

func (r *PostRepository) CreatePost(ctx context.Context, params CreatePostParams) (Post, error) {

	row, err := r.q.CreatePost(ctx, sqlc.CreatePostParams{
		AuthorID: params.AuthorID,
		Title:    params.Title,
		Content:  params.Content,
	})
	if err != nil {
		return Post{}, fmt.Errorf("error on repo: %w", err)
	}

	return Post{
		ID:        row.ID,
		AuthorID:  row.AuthorID,
		Title:     row.Title,
		Content:   row.Content,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
