package user

import "errors"

var (
	ErrEmailTaken    = errors.New("email already taken")
	ErrUsernameTaken = errors.New("username already taken")
	ErrConflict      = errors.New("conflict")
)
