-- +goose Up
CREATE TABLE IF NOT EXISTS users(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  email TEXT NOT NULL UNIQUE,
  username TEXT UNIQUE,

  password_hash TEXT NOT NULL,

  is_active BOOLEAN NOT NULL DEFAULT true,

  token_invalid_before TIMESTAMPTZ NOT NULL DEFAULT now(),

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


-- +goose Down
DROP TABLE IF EXISTS users;
