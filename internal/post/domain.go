// Package post placeholder
package post

import "time"

// Status represents the publication state of a post.
type Status string

// Post is the core domain model for a blog post.
type Post struct {
	ID        int64
	AuthorID  int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
