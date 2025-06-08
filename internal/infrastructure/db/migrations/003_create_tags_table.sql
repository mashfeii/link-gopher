-- +goose Up
CREATE TABLE IF NOT EXISTS tags (
  id SERIAL PRIMARY KEY,
  chat_id BIGINT REFERENCES users(chat_id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  UNIQUE (chat_id, name)
);

-- +goose Down
DROP TABLE IF EXISTS tags;
