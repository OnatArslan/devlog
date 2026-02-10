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

// UserHandler maps HTTP requests to user service operations.
type UserHandler struct {
	svc      *userService
	validate *validator.Validate
}

// NewUserHandler constructs a handler with service and validator dependencies.
func NewUserHandler(svc *userService, validate *validator.Validate) *UserHandler {
	// Return an HTTP handler that delegates business logic to the service.
	return &UserHandler{
		svc:      svc,
		validate: validate,
	}
}

// REGISTER ---------------------
// SignUpRequest is the expected JSON payload for user registration.
type SignUpRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,alphanum"`
	Password        string `json:"password" validate:"required,strong-password"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`
}

// SignUpResponse is the public response body returned after successful registration.
type SignUpResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SignUp handles user registration requests and writes a safe public user payload.
func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Decode and validate the incoming JSON request body.
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

	// Call the service layer to create a new user account.
	user, err := h.svc.SignUp(r.Context(), SignUpInput{Email: req.Email, Username: req.Username, Password: req.Password})
	if err != nil {
		// Map domain conflicts and unexpected failures to HTTP status codes.
		switch {
		case errors.Is(err, ErrEmailTaken), errors.Is(err, ErrUsernameTaken), errors.Is(err, ErrConflict):
			httpx.WriteError(w, http.StatusConflict, err)
		default:
			httpx.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Return only safe public user fields in the response body.
	httpx.WriteJSON(w, http.StatusCreated, SignUpResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// SignInRequest is the expected JSON payload for password-based authentication.
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// SignInResponse returns user-safe identity fields with access token details.
type SignInResponse struct {
	User        SignInUserResult `json:"user"`
	AccessToken string           `json:"access_token"`
	TokenType   string           `json:"token_type"` // "Bearer"
	ExpiresAt   time.Time        `json:"expires_at"`
}

// SignInUserResult contains only non-sensitive user fields for login responses.
type SignInUserResult struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// SignIn validates request input, authenticates credentials, and returns token payload.
func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	// Decode and validate signin credentials from request body.
	var req SignInRequest

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

	signInInput := SignInInput{
		Email:    req.Email,
		Password: req.Password,
	}

	// Delegate authentication and token creation to the service.
	signInOutput, err := h.svc.SignIn(r.Context(), signInInput)
	if err != nil {
		// Map authentication and server-side errors to proper HTTP statuses.
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			httpx.WriteError(w, http.StatusUnauthorized, err)
		case errors.Is(err, ErrJWTSecretNotSet):
			httpx.WriteError(w, http.StatusInternalServerError, err)
		default:
			httpx.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Build a safe user payload without sensitive fields.
	userResponse := SignInUserResult{
		ID:       signInOutput.User.ID,
		Email:    signInOutput.User.Email,
		Username: signInOutput.User.Username,
	}

	// Build the final signin response containing token metadata.
	signInResponse := SignInResponse{
		User:        userResponse,
		AccessToken: signInOutput.Token,
		TokenType:   "Bearer",
		ExpiresAt:   signInOutput.ExpiresAt,
	}

	httpx.WriteJSON(w, http.StatusOK, signInResponse)
}

type GetMeResponse struct {
	ID                 int64     `json:"id"`
	Email              string    `json:"email"`
	Username           string    `json:"username"`
	IsActive           bool      `json:"is_active"`
	TokenInvalidBefore time.Time `json:"token_invalid_before"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	ctxUser, ok := AuthUserFromContext(r.Context())

	if !ok {
		httpx.WriteError(w, http.StatusNotFound, ErrUserNotFound)
		return
	}

	user, err := h.svc.GetMe(r.Context(), ctxUser.Email)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			httpx.WriteError(w, http.StatusNotFound, err)
		default:
			httpx.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	response := GetMeResponse{
		ID:                 user.ID,
		Email:              user.Email,
		Username:           user.Username,
		IsActive:           user.IsActive,
		TokenInvalidBefore: user.TokenInvalidBefore,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}

	httpx.WriteJSON(w, http.StatusOK, response)
}

// Routes registers user HTTP routes under the provided chi router.
func (h *UserHandler) Routes(r chi.Router) chi.Router {
	// Register user auth endpoints on the provided router.
	r.Post("/signup", h.SignUp)
	r.Post("/signin", h.SignIn)
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Get("/me", h.GetMe)
	})
	return r
}
