package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLTagRepository struct {
	db *pgxpool.Pool
}

func NewSQLTagRepository(db *pgxpool.Pool) *SQLTagRepository {
	return &SQLTagRepository{db: db}
}

func (r *SQLTagRepository) GetTagsByUser(ctx context.Context, chatID int64) ([]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(ctx))
		}
	}()

	rows, err := tx.Query(ctx, `SELECT name FROM tags WHERE chat_id = $1`, chatID)
	if err != nil {
		return nil, err
	}

	var tags []string

	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return nil, err
		}

		tags = append(tags, tagName)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *SQLTagRepository) AddTagToUser(ctx context.Context, userID int64, tagName string) (*int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	var tagID int64

	err = tx.QueryRow(ctx, `INSERT INTO tags (chat_id, name) VALUES ($1, $2) RETURNING id`, userID, tagName).Scan(&tagID) //nolint:lll // straight sql request
	if err != nil {
		return nil, err
	}

	return &tagID, tx.Commit(ctx)
}

func (r *SQLTagRepository) AddTagToLink(ctx context.Context, chatID int64, linkURL, tagName string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	_, err = tx.Exec(ctx, `
    INSERT INTO links_tags (tag_id, link_id)
    SELECT t.id, l.id
    FROM tags t
    JOIN links l ON l.url = $2 AND l.chat_id = $3
    WHERE t.name = $1 AND t.chat_id = $3
    ON CONFLICT DO NOTHING`,
		tagName, linkURL, chatID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *SQLTagRepository) RemoveTagsFromLink(ctx context.Context, linkID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(ctx))
		}
	}()

	_, err = tx.Exec(ctx, `DELETE FROM links_tags WHERE link_id = $1`, linkID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *SQLTagRepository) GetTagsByLink(ctx context.Context, linkID int64) ([]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(ctx))
		}
	}()

	rows, err := tx.Query(ctx,
		`SELECT t.name
    FROM links_tags lt
    JOIN tags t ON lt.tag_id = t.id
    WHERE lt.link_id = $1`,
		linkID,
	)
	if err != nil {
		return nil, err
	}

	var tags []string

	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return nil, err
		}

		tags = append(tags, tagName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return tags, nil
}
