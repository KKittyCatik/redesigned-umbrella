package queries

import (
	"context"
	"log/slog"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
)

type MigrationStatus struct {
	Version uint
	Dirty   bool
}

type GetMigrationStatusQuery struct {
	migrationRepo ports.MigrationRepository
	logger        *slog.Logger
}

func NewGetMigrationStatusQuery(migrationRepo ports.MigrationRepository, logger *slog.Logger) *GetMigrationStatusQuery {
	return &GetMigrationStatusQuery{
		migrationRepo: migrationRepo,
		logger:        logger,
	}
}

func (q *GetMigrationStatusQuery) Execute(ctx context.Context) (*MigrationStatus, error) {
	q.logger.Info("Executing query: getting migration status")

	version, dirty, err := q.migrationRepo.GetCurrentVersion(ctx)
	if err != nil {
		q.logger.Error("Failed to get migration status", "error", err)
		return nil, err
	}

	status := &MigrationStatus{
		Version: version,
		Dirty:   dirty,
	}

	q.logger.Info("Query executed successfully",
		"version", version,
		"dirty", dirty,
	)

	return status, nil
}
