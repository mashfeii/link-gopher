package storage

import (
	"context"
	"errors"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLLinkRepository struct {
	db *pgxpool.Pool
}

func NewSQLLinkRepository(db *pgxpool.Pool) *SQLLinkRepository {
	return &SQLLinkRepository{
		db: db,
	}
}

func (r *SQLLinkRepository) AddLink(ctx context.Context, link *models.Link) (*int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	var linkID int64

	err = tx.QueryRow(ctx, "INSERT INTO links (chat_id, url, type, last_checked, last_updated) VALUES ($1, $2, $3, $4, $5) RETURNING id", link.ChatID, link.URL, link.Type, link.LastChecked, link.LastUpdated).Scan(&linkID) //nolint:lll // straight sql query to handle requirements
	if err != nil {
		return nil, err
	}

	return &linkID, tx.Commit(ctx)
}

func (r *SQLLinkRepository) DeleteLink(ctx context.Context, chatID int64, url string) (*models.Link, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	var link models.Link

	err = tx.QueryRow(ctx, "DELETE FROM links WHERE chat_id = $1 AND url = $2 RETURNING id, chat_id, url, type, last_checked, last_updated", chatID, url).Scan(&link.LinkID, &link.ChatID, &link.URL, &link.Type, &link.LastChecked, &link.LastUpdated) //nolint:lll // straight sql query to handle requirements
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &link, tx.Commit(ctx)
}

func (r *SQLLinkRepository) GetLinksByUser(ctx context.Context, chatID int64) ([]models.Link, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	rows, err := tx.Query(ctx, "SELECT id, chat_id, url, type, last_checked, last_updated FROM links WHERE chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link

	for rows.Next() {
		var link models.Link
		if err := rows.Scan(&link.LinkID, &link.ChatID, &link.URL, &link.Type, &link.LastChecked, &link.LastUpdated); err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, tx.Commit(ctx)
}

func (r *SQLLinkRepository) GetLinkByURL(ctx context.Context, chatID int64, url string) (*models.Link, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	var link models.Link

	err = tx.QueryRow(ctx, "SELECT id, chat_id, url, type, last_checked, last_updated FROM links WHERE chat_id = $1 AND url = $2", chatID, url).Scan(&link.LinkID, &link.ChatID, &link.URL, &link.Type, &link.LastChecked, &link.LastUpdated) //nolint:lll // straight sql query to handle requirements
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &link, tx.Commit(ctx)
}

func (r *SQLLinkRepository) GetAllActiveLinks(ctx context.Context) ([]models.Link, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	// TODO: validate interval: from config either hardcoded
	rows, err := tx.Query(ctx, "SELECT id, chat_id, url, type, last_checked, last_updated FROM links WHERE last_checked > NOW() - INTERVAL '1 day'") //nolint:lll // straight sql query to handle requirements
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link

	for rows.Next() {
		var link models.Link
		if err := rows.Scan(&link.LinkID, &link.ChatID, &link.URL, &link.Type, &link.LastChecked, &link.LastUpdated); err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, tx.Commit(ctx)
}

func (r *SQLLinkRepository) UpdateLink(ctx context.Context, link *models.Link) error {
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
        UPDATE links 
        SET 
            last_checked = $1, 
            last_updated = $2 
        WHERE url = $3 AND chat_id = $4`,
		link.LastChecked.UTC(),
		link.LastUpdated.UTC(),
		link.URL,
		link.ChatID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		return err
	}

	return tx.Commit(ctx)
}
