package user

import (
	"encoding/json"
	"net/http"

	"github.com/OnatArslan/devlog/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	svc      *userService
	validate *validator.Validate
}

func NewUserHandler(svc *userService, validate *validator.Validate) *UserHandler {

	return &UserHandler{
		svc:      svc,
		validate: validate,
	}
}

// REGISTER ---------------------
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,alphanum"`
	Password        string `json:"password" validate:"required,min=8,max=64"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	if err := h.validate.Struct(req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, req)
}

func (h *UserHandler) Routes(r chi.Router) chi.Router {
	r.Post("/signup", h.Register)
	return r
}
