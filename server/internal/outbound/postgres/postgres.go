package postgres

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func New(dsn string) (PostgresRepository, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return PostgresRepository{}, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return PostgresRepository{}, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return PostgresRepository{}, err
	}

	return PostgresRepository{Pool: pool}, nil
}

func (r *PostgresRepository) Migrate(dsn string) error {
	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs: %v", err)
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		source,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to load migration: %v", err)
	}

	slog.Info("postgres migration completed successfully")
	return nil
}
