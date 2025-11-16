package ports

import "context"

type MigrationRepository interface {
	ApplyMigrations(ctx context.Context) error
	RollbackMigrations(ctx context.Context) error
	GetCurrentVersion(ctx context.Context) (uint, bool, error)
}
