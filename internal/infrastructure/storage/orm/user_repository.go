package storage

import (
	"context"
	"fmt"

	defaulterrors "errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ORMUserRepository struct {
	db      *pgxpool.Pool
	dialect goqu.DialectWrapper
}

func NewORMUserRepository(db *pgxpool.Pool) *ORMUserRepository {
	return &ORMUserRepository{
		db:      db,
		dialect: goqu.Dialect("postgres"),
	}
}

func (r *ORMUserRepository) AddUser(ctx context.Context, user *models.User) error {
	const op = "ORMUserRepository.AddUser"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.Insert("users").Rows(
		goqu.Record{"chat_id": user.ChatID},
	).Returning("chat_id")

	query, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var chatID int64
	err = tx.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, errors.NewErrUserAlreadyExists(user.ChatID))
		}
		return fmt.Errorf("%s: failed to add user: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *ORMUserRepository) GetUser(ctx context.Context, chatID int64) (*models.User, error) {
	const op = "ORMUserRepository.GetUser"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.From("users").
		Select("chat_id").
		Where(goqu.C("chat_id").Eq(chatID))

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var user models.User
	err = tx.QueryRow(ctx, query, args...).Scan(&user.ChatID)
	if err != nil {
		if defaulterrors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errors.NewErrUserNotFound(chatID))
		}
		return nil, fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return &user, nil
}

func (r *ORMUserRepository) DeleteUser(ctx context.Context, chatID int64) error {
	const op = "ORMUserRepository.DeleteUser"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	ds := r.dialect.Delete("users").
		Where(goqu.C("chat_id").Eq(chatID))

	query, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to delete user: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errors.NewErrUserNotFound(chatID))
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

