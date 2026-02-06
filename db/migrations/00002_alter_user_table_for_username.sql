-- +goose Up
ALTER TABLE users ALTER COLUMN username SET NOT NULL;

-- +goose Down
ALTER TABLE users ALTER COLUMN username SET NULL;
