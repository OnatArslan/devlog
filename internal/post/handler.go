package post

import (
	"net/http"

	"github.com/OnatArslan/devlog/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type PostHandler struct {
	svc      *PostService
	validate *validator.Validate
}

func NewPostHandler(svc *PostService, validate *validator.Validate) *PostHandler {

	return &PostHandler{
		svc:      svc,
		validate: validate,
	}
}

func (h *PostHandler) Example(w http.ResponseWriter, r *http.Request) {

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *PostHandler) Routes(r chi.Router) chi.Router {
	r.Get("/", h.Example)

	return r
}
