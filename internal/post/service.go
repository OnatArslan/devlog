package post

import (
	"context"
	"fmt"
)

// Service contains business logic for post operations.
type Service struct {
	repo *Repository
}

// NewPostService creates a Service wired to the given repository.
func NewPostService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreatePostInput defines the fields required to create a new post.
type CreatePostInput struct {
	AuthorID int64
	Title    string
	Content  string
}

// CreatePost creates a new post and returns the persisted domain model.
func (s *Service) CreatePost(ctx context.Context, input CreatePostInput) (Post, error) {

	post, err := s.repo.CreatePost(ctx, CreatePostParams(input))
	if err != nil {
		return Post{}, fmt.Errorf("create post service : %w", err)
	}
	return post, nil
}

const (
	defaultPageLimit = 20
	maxPageLimit     = 100
)

// ListPostsInput defines pagination parameters for listing posts.
type ListPostsInput struct {
	Limit  int32
	Offset int32
}

// NormalizeListInput clamps limit to valid bounds and ensures offset is non-negative.
func NormalizeListInput(input ListPostsInput) ListPostsInput {
	if input.Limit <= 0 {
		input.Limit = defaultPageLimit
	}
	if input.Limit > maxPageLimit {
		input.Limit = maxPageLimit
	}
	if input.Offset < 0 {
		input.Offset = 0
	}
	return input
}

// GetAllPosts returns paginated posts with their author usernames.
func (s *Service) GetAllPosts(ctx context.Context, input ListPostsInput) ([]Row, error) {
	input = NormalizeListInput(input)

	posts, err := s.repo.GetAllPosts(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("get all posts service: %w", err)
	}

	return posts, nil
}

// GetPostByID returns a single post with author username by ID.
func (s *Service) GetPostByID(ctx context.Context, id int64) (Row, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return Row{}, err
	}

	return post, nil
}
