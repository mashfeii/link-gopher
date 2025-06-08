-- +goose Up
CREATE TABLE IF NOT EXISTS links (
  id SERIAL PRIMARY KEY,
  chat_id BIGINT NOT NULL REFERENCES users(chat_id) ON DELETE CASCADE,
  url TEXT NOT NULL,
  type TEXT NOT NULL,
  last_checked TIMESTAMP DEFAULT NOW(),
  last_updated TIMESTAMP DEFAULT NOW(),
  UNIQUE(chat_id, url)
);

CREATE INDEX IF NOT EXISTS idx_links_chat_id ON links(chat_id);
CREATE INDEX IF NOT EXISTS idx_links_url_chat_id ON links(url, chat_id);

-- +goose Down
DROP INDEX IF EXISTS idx_links_chat_id;
DROP INDEX IF EXISTS idx_links_url_chat_id;
DROP TABLE IF EXISTS links;
