package storage

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ORMFiltersRepository struct {
	db      *pgxpool.Pool
	dialect goqu.DialectWrapper
}

func NewORMFiltersRepository(db *pgxpool.Pool) *ORMFiltersRepository {
	return &ORMFiltersRepository{
		db:      db,
		dialect: goqu.Dialect("postgres"),
	}
}

func (r *ORMFiltersRepository) AddFilterToLink(ctx context.Context, linkID int64, key, value string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, _, err := r.dialect.Insert("filters").
		Rows(goqu.Record{
			"link_id": linkID,
			"key":     key,
			"value":   value,
		}).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *ORMFiltersRepository) RemoveFiltersFromLink(ctx context.Context, linkID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, _, err := r.dialect.Delete("filters").
		Where(goqu.C("link_id").Eq(linkID)).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *ORMFiltersRepository) GetFiltersByLink(ctx context.Context, linkID int64) ([]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, _, err := r.dialect.From("filters").
		Select(goqu.L("key || ':' || value")).
		Where(goqu.C("link_id").Eq(linkID)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filters := []string{}
	for rows.Next() {
		var filter string
		if err := rows.Scan(&filter); err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return filters, tx.Commit(ctx)
}
