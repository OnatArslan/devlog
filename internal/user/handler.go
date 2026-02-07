package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

type SignUpResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// decode request body
	var req SignUpRequest
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.svc.SignUp(r.Context(), SignUpInput{Email: req.Email, Username: req.Username, Password: req.Password})
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken), errors.Is(err, ErrUsernameTaken), errors.Is(err, ErrConflict):
			httpx.WriteError(w, http.StatusConflict, err)
		default:
			httpx.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, SignUpResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

type SignInRequest struct {
	Email string `json:"email"`
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.svc.SignIn(r.Context(), req.Email)

	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, user)

}

func (h *UserHandler) Routes(r chi.Router) chi.Router {
	r.Post("/signup", h.SignUp)
	r.Post("/signin", h.SignIn)
	return r
}
