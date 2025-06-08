package storage

import (
	"context"
	defaulterrors "errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ORMLinksRepository struct {
	db      *pgxpool.Pool
	dialect goqu.DialectWrapper
}

func NewORMLinksRepository(db *pgxpool.Pool) *ORMLinksRepository {
	return &ORMLinksRepository{
		db:      db,
		dialect: goqu.Dialect("postgres"),
	}
}

func (r *ORMLinksRepository) AddLink(ctx context.Context, link *models.Link) (*int64, error) {
	const op = "ORMLinksRepository.AddLink"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.Insert("links").Rows(
		goqu.Record{
			"chat_id":      link.ChatID,
			"url":          link.URL,
			"type":         link.Type,
			"last_checked": link.LastChecked,
			"last_updated": link.LastUpdated,
		},
	).Returning("id")

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var id int64
	err = tx.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to add link: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return &id, nil
}

func (r *ORMLinksRepository) DeleteLink(ctx context.Context, chatID int64, url string) (*models.Link, error) {
	const op = "ORMLinksRepository.DeleteLink"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.Delete("links").
		Where(goqu.Ex{"chat_id": chatID, "url": url}).
		Returning("id", "chat_id", "url", "type", "last_checked", "last_updated")

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var link models.Link
	err = tx.QueryRow(ctx, query, args...).Scan(
		&link.LinkID,
		&link.ChatID,
		&link.URL,
		&link.Type,
		&link.LastChecked,
		&link.LastUpdated,
	)
	if err != nil {
		if defaulterrors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errors.NewErrLinkNotFound())
		}
		return nil, fmt.Errorf("%s: failed to delete link: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return &link, nil
}

func (r *ORMLinksRepository) GetLinksByUser(ctx context.Context, chatID int64) ([]models.Link, error) {
	const op = "ORMLinksRepository.GetLinksByUser"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.From("links").
		Select("id", "chat_id", "url", "type", "last_checked", "last_updated").
		Where(goqu.C("chat_id").Eq(chatID))

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get links: %w", op, err)
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var link models.Link
		if err := rows.Scan(
			&link.LinkID,
			&link.ChatID,
			&link.URL,
			&link.Type,
			&link.LastChecked,
			&link.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan link: %w", op, err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error iterating over rows: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return links, nil
}

func (r *ORMLinksRepository) GetLinkByURL(ctx context.Context, chatID int64, url string) (*models.Link, error) {
	const op = "ORMLinksRepository.GetLinkByURL"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.From("links").
		Select("id", "chat_id", "url", "type", "last_checked", "last_updated").
		Where(goqu.C("chat_id").Eq(chatID), goqu.C("url").Eq(url))

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var link models.Link
	err = tx.QueryRow(ctx, query, args...).Scan(
		&link.LinkID,
		&link.ChatID,
		&link.URL,
		&link.Type,
		&link.LastChecked,
		&link.LastUpdated,
	)
	if err != nil {
		if defaulterrors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errors.NewErrLinkNotFound())
		}
		return nil, fmt.Errorf("%s: failed to get link by URL: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return &link, nil
}

func (r *ORMLinksRepository) GetAllActiveLinks(ctx context.Context) ([]models.Link, error) {
	const op = "ORMLinksRepository.GetAllActiveLinks"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.From("links").
		Select("id", "chat_id", "url", "type", "last_checked", "last_updated").
		Where(goqu.C("last_checked").Gt(goqu.L("NOW() - INTERVAL '1 day'")))

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get all active links: %w", op, err)
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var link models.Link
		if err := rows.Scan(
			&link.LinkID,
			&link.ChatID,
			&link.URL,
			&link.Type,
			&link.LastChecked,
			&link.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan link: %w", op, err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error iterating over rows: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return links, nil
}

func (r *ORMLinksRepository) UpdateLink(ctx context.Context, link *models.Link) error {
	const op = "ORMLinksRepository.UpdateLink"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.Update("links").
		Set(goqu.Record{
			"last_checked": link.LastChecked,
			"last_updated": link.LastUpdated,
		}).
		Where(
			goqu.C("chat_id").Eq(link.ChatID),
			goqu.C("url").Eq(link.URL),
		)

	query, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to update link: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errors.NewErrLinkNotFound())
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

