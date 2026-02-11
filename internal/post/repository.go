package post

import (
	"context"
	"fmt"
	"time"

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

type GetAllPostsRow struct {
	ID        int64
	AuthorID  int64
	Username  string
	Title     string
	Content   string
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (r *PostRepository) GetAllPosts(ctx context.Context) ([]GetAllPostsRow, error) {

	rows, err := r.q.GetAllPosts(ctx)
	if err != nil {
		return []GetAllPostsRow{}, err
	}

	posts := make([]GetAllPostsRow, 0, len(rows))

	for _, row := range rows {
		posts = append(posts, GetAllPostsRow{
			ID:        row.ID,
			AuthorID:  row.AuthorID,
			Username:  row.Username,
			Title:     row.Title,
			Content:   row.Content,
			UpdatedAt: row.UpdatedAt.Time,
			CreatedAt: row.CreatedAt.Time,
		})
	}
	return posts, nil
}
