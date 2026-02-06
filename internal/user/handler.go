package user

import (
	"encoding/json"
	"errors"
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
type SignUpRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,alphanum"`
	Password        string `json:"password" validate:"required,min=8,max=64"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
}

func (h *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	if err := h.validate.Struct(req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.svc.SignUpUser(r.Context(), SignUpInput{Email: req.Email, Username: req.Username, Password: req.Password})
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken), errors.Is(err, ErrUsernameTaken), errors.Is(err, ErrConflict):
			httpx.WriteError(w, http.StatusConflict, err)
		default:
			httpx.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, ResponseUser{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (h *UserHandler) Routes(r chi.Router) chi.Router {
	r.Post("/signup", h.SignUpHandler)
	return r
}
