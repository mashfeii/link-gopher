-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  chat_id BIGINT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
