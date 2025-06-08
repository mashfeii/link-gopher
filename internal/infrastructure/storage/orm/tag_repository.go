package storage

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ORMTagsRepository struct {
	db      *pgxpool.Pool
	dialect goqu.DialectWrapper
}

func NewORMTagsRepository(db *pgxpool.Pool) *ORMTagsRepository {
	return &ORMTagsRepository{
		db:      db,
		dialect: goqu.Dialect("postgres"),
	}
}

func (r *ORMTagsRepository) GetTagsByUser(ctx context.Context, chatID int64) ([]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, _, err := r.dialect.From("tags").
		Select("name").
		Where(goqu.C("chat_id").Eq(chatID)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return tags, tx.Commit(ctx)
}

func (r *ORMTagsRepository) AddTagToUser(ctx context.Context, userID int64, tagName string) (*int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, _, err := r.dialect.Insert("tags").
		Rows(goqu.Record{
			"chat_id": userID,
			"name":    tagName,
		}).
		Returning("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	var tagID int64
	err = tx.QueryRow(ctx, query).Scan(&tagID)
	if err != nil {
		return nil, err
	}

	return &tagID, tx.Commit(ctx)
}

func (r *ORMTagsRepository) AddTagToLink(ctx context.Context, chatID int64, linkURL, tagName string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tagSubquerySQL, _, err := r.dialect.From("tags").
		Select("id").
		Where(goqu.C("name").Eq(tagName), goqu.C("chat_id").Eq(chatID)).
		ToSQL()
	if err != nil {
		return err
	}

	linkSubquerySQL, _, err := r.dialect.From("links").
		Select("id").
		Where(goqu.C("url").Eq(linkURL), goqu.C("chat_id").Eq(chatID)).
		ToSQL()
	if err != nil {
		return err
	}

	insertQuery, _, err := r.dialect.Insert("links_tags").
		Cols("tag_id", "link_id").
		FromQuery(
			goqu.Select(goqu.I("t.id").As("tag_id"), goqu.I("l.id").As("link_id")).
				From(goqu.T("t")).
				CrossJoin(goqu.T("l")).
				Where(goqu.I("t.id").IsNotNull(), goqu.I("l.id").IsNotNull()),
		).
		OnConflict(goqu.DoNothing()).
		ToSQL()
	if err != nil {
		return err
	}

	finalSQL := fmt.Sprintf(`
        WITH t AS (%s),
             l AS (%s)
        %s`, tagSubquerySQL, linkSubquerySQL, insertQuery)

	_, err = tx.Exec(ctx, finalSQL)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *ORMTagsRepository) RemoveTagsFromLink(ctx context.Context, linkID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, _, err := r.dialect.Delete("links_tags").
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

func (r *ORMTagsRepository) GetTagsByLink(ctx context.Context, linkID int64) ([]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, _, err := r.dialect.From("links_tags").
		Join(goqu.T("tags"), goqu.On(goqu.Ex{"links_tags.tag_id": goqu.I("tags.id")})).
		Select("tags.name").
		Where(goqu.C("links_tags.link_id").Eq(linkID)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return tags, tx.Commit(ctx)
}
