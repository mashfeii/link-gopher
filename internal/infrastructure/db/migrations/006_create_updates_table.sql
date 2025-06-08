-- +goose Up
CREATE TABLE IF NOT EXISTS updates (
  id BIGSERIAL PRIMARY KEY,
  link_id BIGINT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
  update_type TEXT NOT NULL CHECK (update_type IN ('answer', 'comment', 'issue', 'pull_request')),
  title TEXT NOT NULL,
  username TEXT NOT NULL,
  body_preview TEXT NOT NULL,
  created_at TIMESTAMP,
  sent BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_updates_link_id_sent ON updates(link_id, sent);

-- +goose Down
DROP INDEX IF EXISTS idx_updates_link_id_sent;
DROP TABLE IF EXISTS updates;
