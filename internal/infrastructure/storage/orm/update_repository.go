package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ORMUpdateRepository struct {
	db *pgxpool.Pool
	ds goqu.DialectWrapper
}

func NewORMUpdateRepository(db *pgxpool.Pool) *ORMUpdateRepository {
	return &ORMUpdateRepository{
		db: db,
		ds: goqu.Dialect("postgres"),
	}
}

func (r *ORMUpdateRepository) SaveUpdate(ctx context.Context, update *models.Update) error {
	const op = "repository.ORMUpdateRepository.SaveUpdate"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	ds := r.ds.Insert("updates").Rows(
		goqu.Record{
			"link_id":      update.LinkID,
			"update_type":  update.UpdateType,
			"title":        update.Title,
			"username":     update.UserName,
			"body_preview": update.BodyPreview,
			"created_at":   update.CreatedAt,
			"sent":         update.Sent,
		},
	)

	query, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("%s: failed to build insert query: %w", op, err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to insert update: %w", op, err)
	}

	return tx.Commit(ctx)
}

func (r *ORMUpdateRepository) GetUpdateByLink(ctx context.Context, linkID int64, limit, offset uint) ([]models.Update, error) {
	const op = "repository.ORMUpdateRepository.GetUpdateByLink"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	ds := r.ds.From("updates").Select(
		"id",
		"link_id",
		"update_type",
		"title",
		"username",
		"body_preview",
		"created_at",
		"sent",
	).
		Where(goqu.Ex{"link_id": linkID}).
		Order(goqu.I("created_at").Desc()).
		Offset(offset).
		Limit(limit)

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build select query: %w", op, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query updates: %w", op, err)
	}

	defer rows.Close()

	var updates []models.Update
	for rows.Next() {
		var update models.Update
		if err := rows.Scan(
			&update.ID,
			&update.LinkID,
			&update.UpdateType,
			&update.Title,
			&update.UserName,
			&update.BodyPreview,
			&update.CreatedAt,
			&update.Sent,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan update: %w", op, err)
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
