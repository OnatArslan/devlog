package user

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/OnatArslan/devlog/internal/httpx"
	"github.com/golang-jwt/jwt/v5"
)

// ctxKey is a private type used to prevent context key collisions.
type ctxKey string

const authUserKey ctxKey = "auth_user"

// AuthUser contains minimal authenticated identity passed through request context.
type AuthUser struct {
	ID       int64
	Email    string
	Username string
}

// AuthUserFromContext returns authenticated user data if middleware populated it.
func AuthUserFromContext(ctx context.Context) (AuthUser, bool) {
	user, ok := ctx.Value(authUserKey).(AuthUser)
	return user, ok
}

// AuthMiddleware validates Bearer JWT tokens and injects auth user data into context.
func (h *UserHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and validate the Authorization header format.
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			httpx.WriteError(w, http.StatusUnauthorized, ErrInvalidCredentials)
			return
		}

		// Strip "Bearer " prefix and trim trailing/leading spaces.
		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if tokenStr == "" {
			httpx.WriteError(w, http.StatusUnauthorized, ErrInvalidCredentials)
			return
		}

		// Ensure JWT secret exists before attempting signature verification.
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			httpx.WriteError(w, http.StatusInternalServerError, ErrJWTSecretNotSet)
			return
		}

		// Parse and validate token signature and standard claims into custom claims.
		claims := &CustomClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			httpx.WriteError(w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		// Build context-safe auth payload for downstream protected handlers.
		authUser := AuthUser{
			ID:       claims.UserID,
			Email:    claims.Email,
			Username: claims.Username,
		}

		ctx := context.WithValue(r.Context(), authUserKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
