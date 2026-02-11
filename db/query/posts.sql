-- name: CreatePost :one
INSERT INTO posts (author_id, title, content)
VALUES ($1, $2, $3)
RETURNING *
;


-- name: GetAllPosts :many
SELECT p.*, u.username FROM posts p JOIN users u ON u.id = p.author_id;


-- name: GetPostById :one
SELECT p.*, u.username FROM posts p JOIN users u ON u.id = p.author_id
WHERE p.id = $1;

-- name: UpdatePost :one
UPDATE posts
SET
  title = COALESCE(sqlc.narg(title), title),
  content = COALESCE(sqlc.narg(content), content),
  updated_at = now()
WHERE id = $1
RETURNING *;
