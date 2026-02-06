package user

import "time"

// User is the domain entity used by the user package.
// It intentionally avoids database-specific types.
type User struct {
	ID                 int64
	Email              string
	Username           string
	PasswordHash       string
	IsActive           bool
	TokenInvalidBefore time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type CreateUserInput struct {
	Email        string
	Username     string
	PasswordHash string
}

type SignUpInput struct {
	Email    string
	Username string
	Password string
}

type ResponseUser struct {
	ID        int64
	Email     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
