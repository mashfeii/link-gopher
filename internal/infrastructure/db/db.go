package db

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewPool(cfg *config.Database) (*pgxpool.Pool, error) {
	const op = "db.NewPool"

	pool, err := createPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := ApplyMigrations(pool.Config().ConnConfig); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pool, nil
}

func createPool(cfg *config.Database) (*pgxpool.Pool, error) {
	const op = "db.createPool"

	connConfig, err := pgxpool.ParseConfig(cfg.ToDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	connConfig.MaxConns = 32
	connConfig.MaxConnIdleTime = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create connection pool: %w", op, err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	return pool, nil
}

func ApplyMigrations(cfg *pgx.ConnConfig) error {
	const op = "db.ApplyMigrations"

	goose.SetBaseFS(embedMigrations)

	db := stdlib.OpenDB(*cfg)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("%s: failed to set dialect: %w", op, err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("%s: failed to apply migrations: %w", op, err)
	}

	return nil
}
