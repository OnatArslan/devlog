package user

import "errors"

// Domain-level user/auth errors shared across repository, service, and handler layers.
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already taken")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrConflict           = errors.New("conflict")
	ErrJWTSecretNotSet    = errors.New("jwt secret not set")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUnknownClaimsType  = errors.New("unknown claims type")
	ErrWeakPassword       = errors.New("password is weak")
)
