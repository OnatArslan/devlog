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

// GetAllPosts returns all posts with their author usernames.
func (s *Service) GetAllPosts(ctx context.Context) ([]Row, error) {

	posts, err := s.repo.GetAllPosts(ctx)
	if err != nil {
		return []Row{}, err
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
