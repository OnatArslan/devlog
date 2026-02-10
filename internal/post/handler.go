package post

import (
	"errors"
	"fmt"
	"net/http"

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

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {

	user, ok := user.AuthUserFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusNonAuthoritativeInfo, errors.New("auth user can not found"))
		return
	}

	fmt.Println(user)

	httpx.WriteJSON(w, http.StatusCreated, user)
}

func (h *PostHandler) Routes(r chi.Router) chi.Router {
	r.Get("/", h.Example)
	r.Group(func(r chi.Router) {
		r.Use(h.authMW)
		r.Post("/", h.CreatePost)
	})

	return r
}
