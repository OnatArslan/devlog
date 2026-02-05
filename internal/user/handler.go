package user

import (
	"net/http"

	"github.com/OnatArslan/devlog/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	svc *userService
}

func NewUserHandler(svc *userService) *UserHandler {

	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) Example(w http.ResponseWriter, r *http.Request) {

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"STATUS": "OK"})
}

func (h *UserHandler) Routes(r chi.Router) chi.Router {
	r.Get("/", h.Example)
	return r
}
