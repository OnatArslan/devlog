package post

import "time"

type Status string

type Post struct {
	ID        int64
	AuthorID  int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
