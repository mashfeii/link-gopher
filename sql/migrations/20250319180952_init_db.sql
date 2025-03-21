-- +goose Up
-- +goose StatementBegin
begin;

CREATE TABLE IF NOT EXISTS users (
    chat_id BIGINT PRIMARY KEY,
);

CREATE TABLE IF NOT EXISTS links (
  link_id SERIAL PRIMARY KEY,
  chat_id BIGINT REFERENCES users(chat_id) ON DELETE CASCADE,
  url TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE (chat_id, url)
);

CREATE TABLE IF NOT EXISTS tags (
  tag_id SERIAL PRIMARY KEY,
  tag_name TEXT NOT NULL UNIQUE,
);

CREATE TABLE IF NOT EXISTS links_tags (
  link_id INT REFERENCES links(link_id) ON DELETE CASCADE,
  tag_id INT REFERENCES tags(tag_id) ON DELETE CASCADE,
  PRIMARY KEY (link_id, tag_id)
);

CREATE TABLE IF NOT EXISTS filters (
  filter_id SERIAL PRIMARY KEY,
  link_id INT REFERENCES links(link_id) ON DELETE CASCADE,
  filter_key TEXT NOT NULL,
  filter_value TEXT NOT NULL,
);

end;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
begin;

DROP TABLE IF EXISTS filters;
DROP TABLE IF EXISTS links_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS users;

end;
-- +goose StatementEnd
