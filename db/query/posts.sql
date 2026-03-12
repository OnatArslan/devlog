-- name: CreatePost :one
INSERT INTO posts (author_id, title, content)
VALUES ($1, $2, $3)
RETURNING *
;


-- name: GetAllPosts :many
SELECT p.*, u.username FROM posts p JOIN users u ON u.id = p.author_id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;


-- name: GetPostById :one
SELECT p.*, u.username FROM posts p JOIN users u ON u.id = p.author_id
WHERE p.id = $1;
