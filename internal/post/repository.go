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

type PostRow struct {
	ID        int64     `json:"id"`
	AuthorID  int64     `json:"author_id"`
	Username  string    `json:"author_username"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *PostRepository) GetAllPosts(ctx context.Context) ([]PostRow, error) {
	rows, err := r.q.GetAllPosts(ctx)
	if err != nil {
		return []PostRow{}, err
	}
	posts := make([]PostRow, 0, len(rows))
	for _, row := range rows {
		posts = append(posts, PostRow{
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

func (r *PostRepository) GetPostById(ctx context.Context, id int64) (PostRow, error) {

	row, err := r.q.GetPostById(ctx, id)
	if err != nil {
		return PostRow{}, err
	}

	return PostRow{
		ID:        row.ID,
		AuthorID:  row.AuthorID,
		Username:  row.Username,
		Title:     row.Title,
		Content:   row.Content,
		UpdatedAt: row.UpdatedAt.Time,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}
