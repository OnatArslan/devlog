package post

import "time"

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
)

type Post struct {
	ID        int64
	AuthorID  int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
