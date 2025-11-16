package commands

import (
	"context"
	"log/slog"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
)

type RollbackMigrationsCommand struct {
	migrationRepo ports.MigrationRepository
	logger        *slog.Logger
}

func NewRollbackMigrationsCommand(migrationRepo ports.MigrationRepository, logger *slog.Logger) *RollbackMigrationsCommand {
	return &RollbackMigrationsCommand{
		migrationRepo: migrationRepo,
		logger:        logger,
	}
}

func (c *RollbackMigrationsCommand) Execute(ctx context.Context) error {
	c.logger.Info("Executing command: rolling back migrations")

	if err := c.migrationRepo.RollbackMigrations(ctx); err != nil {
		c.logger.Error("Rollback failed", "error", err)
		return err
	}

	c.logger.Info("Command executed: migrations rolled back successfully")
	return nil
}
