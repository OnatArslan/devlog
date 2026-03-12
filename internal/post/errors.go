package post

import "errors"

// ErrPostNotFound is returned when a requested post does not exist.
var (
	ErrPostNotFound = errors.New("post not found")
)
