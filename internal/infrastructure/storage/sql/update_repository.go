package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLUpdateRepository struct {
	db *pgxpool.Pool
}

func NewSQLUpdateRepository(db *pgxpool.Pool) *SQLUpdateRepository {
	return &SQLUpdateRepository{db: db}
}

func (r *SQLUpdateRepository) SaveUpdate(ctx context.Context, update *models.Update) error {
	const op = "repository.SQLUpdateRepository.SaveUpdate"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	const query = `
    INSERT INTO updates (link_id, update_type, title, username, body_preview, created_at, sent)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.Exec(ctx, query,
		update.LinkID, update.UpdateType, update.Title,
		update.UserName, update.BodyPreview, update.CreatedAt,
		update.Sent)
	if err != nil {
		return fmt.Errorf("%s: failed to insert update: %w", op, err)
	}

	return tx.Commit(ctx)
}

func (r *SQLUpdateRepository) GetUpdateByLink(ctx context.Context, linkID int64, limit, offset int) ([]models.Update, error) {
	const op = "repository.SQLUpdateRepository.GetUpdateByLink"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	const query = `
    SELECT id, link_id, update_type, title, username, body_preview, created_at, sent
    FROM updates
    WHERE link_id = $1
    ORDER BY created_at DESC
    LIMIT $2 OFFSET $3`

	rows, err := tx.Query(ctx, query, linkID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query updates: %w", op, err)
	}

	defer rows.Close()

	var updates []models.Update

	for rows.Next() {
		var update models.Update
		if err := rows.Scan(
			&update.ID, &update.LinkID, &update.UpdateType,
			&update.Title, &update.UserName, &update.BodyPreview,
			&update.CreatedAt, &update.Sent); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		updates = append(updates, update)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error during row iteration: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return updates, nil
}
