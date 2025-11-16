package repositories

import (
	"context"
	"database/sql"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresMigrationRepository struct {
	migrate *migrate.Migrate
}

func NewPostgresMigrationRepository(db *sql.DB, migrationsPath string) (ports.MigrationRepository, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return &PostgresMigrationRepository{
		migrate: m,
	}, nil
}

func (r *PostgresMigrationRepository) ApplyMigrations(ctx context.Context) error {
	if err := r.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (r *PostgresMigrationRepository) RollbackMigrations(ctx context.Context) error {
	if err := r.migrate.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (r *PostgresMigrationRepository) GetCurrentVersion(ctx context.Context) (uint, bool, error) {
	return r.migrate.Version()
}
