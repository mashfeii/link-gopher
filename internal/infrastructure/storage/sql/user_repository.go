package storage

import (
	"context"
	"errors"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type SQLUserRepository struct {
	db *pgxpool.Pool
}

func NewSQLUserRepository(db *pgxpool.Pool) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) AddUser(ctx context.Context, user *models.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	_, err = tx.Exec(ctx, "INSERT INTO users (chat_id) VALUES ($1)", user.ChatID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *SQLUserRepository) GetUser(ctx context.Context, chatID int64) (*models.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	row := tx.QueryRow(ctx, "SELECT chat_id FROM users WHERE chat_id = $1", chatID)

	var user models.User

	err = row.Scan(&user.ChatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.NewErrUserNotFound(chatID)
		}

		return nil, err
	}

	return &user, tx.Commit(ctx)
}

func (r *SQLUserRepository) DeleteUser(ctx context.Context, chatID int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	_, err = tx.Exec(ctx, "DELETE FROM users WHERE chat_id = $1", chatID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
