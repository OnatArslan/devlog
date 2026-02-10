-- name: CreatePost :one
INSERT INTO posts (author_id, title, content)
VALUES ($1, $2, $3)
RETURNING *
;
