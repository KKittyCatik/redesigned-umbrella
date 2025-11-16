package commands

import (
	"context"
	"log/slog"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
)

type ApplyMigrationsCommand struct {
	migrationRepo ports.MigrationRepository
	logger        *slog.Logger
}

func NewApplyMigrationsCommand(migrationRepo ports.MigrationRepository, logger *slog.Logger) *ApplyMigrationsCommand {
	return &ApplyMigrationsCommand{
		migrationRepo: migrationRepo,
		logger:        logger,
	}
}

func (c *ApplyMigrationsCommand) Execute(ctx context.Context) error {
	c.logger.Info("Executing command: applying migrations")

	if err := c.migrationRepo.ApplyMigrations(ctx); err != nil {
		c.logger.Error("Migration failed", "error", err)
		return err
	}

	c.logger.Info("Command executed: migrations applied successfully")
	return nil
}
