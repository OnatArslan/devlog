package post

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/OnatArslan/devlog/internal/httpx"
	"github.com/OnatArslan/devlog/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Handler maps HTTP requests to post service operations.
type Handler struct {
	svc      *Service
	validate *validator.Validate
	authMW   func(http.Handler) http.Handler
}

// NewPostHandler constructs a Handler with service, validator, and auth middleware dependencies.
func NewPostHandler(svc *Service, validate *validator.Validate, authMW func(http.Handler) http.Handler) *Handler {

	return &Handler{
		svc:      svc,
		validate: validate,
		authMW:   authMW,
	}
}

// CreatePostResponse is the JSON response body returned after a post is created.
type CreatePostResponse struct {
	ID        int64     `json:"id"`
	AuthorID  int64     `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePostRequest is the expected JSON payload for creating a post.
type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=1"`
	Content string `json:"content" validate:"required,min=1"`
}

// CreatePost handles authenticated post creation requests.
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get auth user
	authUser, ok := user.AuthUserFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("auth user can not found"))
		return
	}
	// Get req data
	var req CreatePostRequest
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate request
	if err := h.validate.Struct(req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	post, err := h.svc.CreatePost(r.Context(), CreatePostInput{
		AuthorID: authUser.ID,
		Title:    req.Title,
		Content:  req.Content,
	})

	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, CreatePostResponse{
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

// GetAllPostsResponse is the JSON response body for listing all posts.
type GetAllPostsResponse struct {
	Count int64 `json:"count"`
	Posts []Row `json:"posts"`
}

// GetAllPosts handles requests to list all posts.
func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {

	posts, err := h.svc.GetAllPosts(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, GetAllPostsResponse{Posts: posts, Count: int64(len(posts))})
}

// GetPostByIDRequest is the URL parameter type for post ID lookups.
type GetPostByIDRequest struct {
	ID int64 `json:"id"`
}

// GetPostByID handles requests to fetch a single post by its ID.
func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	post, err := h.svc.GetPostByID(r.Context(), id)

	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, post)
}

// Routes registers post HTTP routes under the provided chi router.
func (h *Handler) Routes(r chi.Router) chi.Router {
	r.Get("/", h.GetAllPosts)
	r.Get("/{id}", h.GetPostByID)
	r.Group(func(r chi.Router) {
		r.Use(h.authMW)
		r.Post("/", h.CreatePost)
	})

	return r
}
