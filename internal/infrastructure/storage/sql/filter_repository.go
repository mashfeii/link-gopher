package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLFilterRepository struct {
	db *pgxpool.Pool
}

func NewSQLFilterRepository(db *pgxpool.Pool) *SQLFilterRepository {
	return &SQLFilterRepository{
		db: db,
	}
}

func (r *SQLFilterRepository) AddFilterToLink(ctx context.Context, linkID int64, key, value string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(ctx))
		}
	}()

	err = tx.QueryRow(ctx, `INSERT INTO filters (link_id, key, value) VALUES ($1, $2, $3) RETURNING NULL`, linkID, key, value).Scan(nil)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *SQLFilterRepository) RemoveFiltersFromLink(ctx context.Context, linkID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(ctx))
		}
	}()

	_, err = tx.Exec(ctx, `DELETE FROM filters WHERE link_id = $1`, linkID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *SQLFilterRepository) GetFiltersByLink(ctx context.Context, linkID int64) ([]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(ctx))
		}
	}()

	rows, err := tx.Query(ctx, `SELECT key || ':' || value FROM filters WHERE link_id = $1`, linkID)
	if err != nil {
		return nil, err
	}

	filters := make([]string, 0)

	for rows.Next() {
		var filterKeyValue string

		if err := rows.Scan(&filterKeyValue); err != nil {
			return nil, err
		}

		filters = append(filters, filterKeyValue)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return filters, nil
}
