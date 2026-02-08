package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already taken")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrConflict           = errors.New("conflict")
	ErrJWTSecretNotSet    = errors.New("jwt secret not set")
)
