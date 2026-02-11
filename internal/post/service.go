package post

import (
	"context"
	"fmt"
)

type PostService struct {
	repo *PostRepository
}

func NewPostService(repo *PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

type CreatePostInput struct {
	AuthorID int64
	Title    string
	Content  string
}

func (s *PostService) CreatePost(ctx context.Context, input CreatePostInput) (Post, error) {

	post, err := s.repo.CreatePost(ctx, CreatePostParams{
		AuthorID: input.AuthorID,
		Title:    input.Title,
		Content:  input.Content,
	})
	if err != nil {
		return Post{}, fmt.Errorf("create post service : %w", err)
	}
	return post, nil
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]PostRow, error) {

	posts, err := s.repo.GetAllPosts(ctx)
	if err != nil {
		return []PostRow{}, err
	}

	return posts, nil

}

func (s *PostService) GetPostById(ctx context.Context, id int64) (PostRow, error) {

	post, err := s.repo.GetPostById(ctx, id)
	if err != nil {
		return PostRow{}, err
	}

	return post, nil
}
