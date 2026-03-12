package post

import (
	"context"
	"fmt"
	"time"

	"github.com/OnatArslan/devlog/internal/sqlc"
)

// Repository provides post persistence operations backed by sqlc queries.
type Repository struct {
	q *sqlc.Queries
}

// NewPostRepository creates a Repository wired to the given sqlc query set.
func NewPostRepository(q *sqlc.Queries) *Repository {
	return &Repository{
		q: q,
	}
}

// CreatePostParams defines the input fields required to insert a new post row.
type CreatePostParams struct {
	AuthorID int64
	Title    string
	Content  string
}

// CreatePost inserts a new post and returns the created domain model.
func (r *Repository) CreatePost(ctx context.Context, params CreatePostParams) (Post, error) {

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

// Row is a post enriched with the author's username, used in list/detail responses.
type Row struct {
	ID        int64     `json:"id"`
	AuthorID  int64     `json:"author_id"`
	Username  string    `json:"author_username"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// GetAllPosts returns all posts joined with their author username.
func (r *Repository) GetAllPosts(ctx context.Context) ([]Row, error) {
	rows, err := r.q.GetAllPosts(ctx)
	if err != nil {
		return []Row{}, err
	}
	posts := make([]Row, 0, len(rows))
	for _, row := range rows {
		posts = append(posts, Row{
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

// GetPostByID returns a single post with author username by post ID.
func (r *Repository) GetPostByID(ctx context.Context, id int64) (Row, error) {

	row, err := r.q.GetPostById(ctx, id)
	if err != nil {
		return Row{}, err
	}

	return Row{
		ID:        row.ID,
		AuthorID:  row.AuthorID,
		Username:  row.Username,
		Title:     row.Title,
		Content:   row.Content,
		UpdatedAt: row.UpdatedAt.Time,
		CreatedAt: row.CreatedAt.Time,
	}, nil
}
