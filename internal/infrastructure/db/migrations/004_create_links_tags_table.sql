-- +goose Up
CREATE TABLE IF NOT EXISTS links_tags (
  link_id INT REFERENCES links(id) ON DELETE CASCADE,
  tag_id INT REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (link_id, tag_id)
);

-- +goose Down
DROP TABLE IF EXISTS links_tags;
