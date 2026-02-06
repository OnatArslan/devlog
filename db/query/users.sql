
-- name: CreateUser :one
INSERT INTO users (email, username, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetByEmail :one
SELECT
  u.id,
  u.email,
  u.username,
  u.password_hash,
  u.is_active,
  u.token_invalid_before,
  u.created_at,
  u.updated_at
FROM users u
WHERE u.email = $1;
