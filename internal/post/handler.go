package post

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/OnatArslan/devlog/internal/httpx"
	"github.com/OnatArslan/devlog/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type PostHandler struct {
	svc      *PostService
	validate *validator.Validate
	authMW   func(http.Handler) http.Handler
}

func NewPostHandler(svc *PostService, validate *validator.Validate, authMW func(http.Handler) http.Handler) *PostHandler {

	return &PostHandler{
		svc:      svc,
		validate: validate,
		authMW:   authMW,
	}
}

func (h *PostHandler) Example(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type CreatePostResponse struct {
	ID        int64     `json:"id"`
	AuthorId  int64     `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=1"`
	Content string `json:"content" validate:"required,min=1"`
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get auth user
	user, ok := user.AuthUserFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusNonAuthoritativeInfo, errors.New("auth user can not found"))
		return
	}
	// Get req data
	var req CreatePostRequest
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		fmt.Println("xd")
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate request
	if err := h.validate.Struct(req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	post, err := h.svc.CreatePost(r.Context(), CreatePostInput{
		AuthorID: user.ID,
		Title:    req.Title,
		Content:  req.Content,
	})

	if err != nil {
		httpx.WriteError(w, http.StatusConflict, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, CreatePostResponse{
		ID:        post.ID,
		AuthorId:  post.AuthorID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

func (h *PostHandler) Routes(r chi.Router) chi.Router {
	r.Get("/", h.Example)
	r.Group(func(r chi.Router) {
		r.Use(h.authMW)
		r.Post("/", h.CreatePost)
	})

	return r
}
