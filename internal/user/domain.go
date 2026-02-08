package user

import "time"

// User is the core domain model used across handler, service, and repository layers.
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
