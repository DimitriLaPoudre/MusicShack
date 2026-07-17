package sqlite

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type SQLiteRepository struct {
	DB *sqlx.DB
}

func New(filepath string) (SQLiteRepository, error) {
	db, err := sqlx.Open(
		"sqlite",
		fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_time_format=sqlite", filepath),
	)
	if err != nil {
		return SQLiteRepository{}, fmt.Errorf("open connection to SQLite file: %s: %w", filepath, err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return SQLiteRepository{}, err
	}

	return SQLiteRepository{DB: db}, nil
}

func (r *SQLiteRepository) Migrate(filepath string) error {
	driver, err := migratesqlite.WithInstance(r.DB.DB, &migratesqlite.Config{})
	if err != nil {
		return fmt.Errorf("create sqlite driver: %w", err)
	}

	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("create iofs: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		source,
		"sqlite",
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrator instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("load migration: %w", err)
	}

	slog.Info("sqlite migration completed successfully")
	return nil
}
