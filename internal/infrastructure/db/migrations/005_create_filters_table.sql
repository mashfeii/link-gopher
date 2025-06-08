-- +goose Up
CREATE TABLE IF NOT EXISTS filters (
  id SERIAL PRIMARY KEY,
  link_id INT REFERENCES links(id) ON DELETE CASCADE,
  key TEXT NOT NULL,
  value TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS filters;
