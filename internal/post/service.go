package post

type PostService struct {
	repo *PostRepository
}

func NewPostService(repo *PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}
