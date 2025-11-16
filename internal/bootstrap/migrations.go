package bootstrap

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/commands"
	"github.com/KKittyCatik/redesigned-umbrella/internal/application/queries"
	"github.com/KKittyCatik/redesigned-umbrella/internal/infrastructure/config"
	"github.com/KKittyCatik/redesigned-umbrella/internal/infrastructure/repositories"
)

func runMigrationCommand(db *sql.DB, logger *slog.Logger, cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	migrationsPath := cfg.DB.MigrationsPath
	if migrationsPath == "" {
		migrationsPath = "./migrations"
	}

	migrationRepo, err := repositories.NewPostgresMigrationRepository(db, migrationsPath)
	if err != nil {
		logger.Error("failed to create migration repository", "error", err)
		os.Exit(1)
	}

	switch cfg.Command.Name {
	case "migrate":
		logger.Info("Starting database migrations...")
		applyCmd := commands.NewApplyMigrationsCommand(migrationRepo, logger)
		if err := applyCmd.Execute(ctx); err != nil {
			logger.Error("migration failed", "error", err)
			os.Exit(1)
		}
		logger.Info("Migrations applied successfully")

	case "rollback":
		logger.Info("Starting database rollback...")
		rollbackCmd := commands.NewRollbackMigrationsCommand(migrationRepo, logger)
		if err := rollbackCmd.Execute(ctx); err != nil {
			logger.Error("rollback failed", "error", err)
			os.Exit(1)
		}
		logger.Info("Migrations rolled back successfully")

	case "migration-status":
		logger.Info("Checking migration status...")
		statusQuery := queries.NewGetMigrationStatusQuery(migrationRepo, logger)
		status, err := statusQuery.Execute(ctx)
		if err != nil {
			logger.Error("failed to get migration status", "error", err)
			os.Exit(1)
		}
		logger.Info("Migration status", "version", status.Version, "dirty", status.Dirty)
	}
}
