package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/KKittyCatik/redesigned-umbrella/internal/infrastructure/config"
	"github.com/KKittyCatik/redesigned-umbrella/internal/infrastructure/http"
	"github.com/KKittyCatik/redesigned-umbrella/internal/infrastructure/repositories"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/commands"
	"github.com/KKittyCatik/redesigned-umbrella/internal/application/queries"

	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/services"
)

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("Starting PR Reviewer Service")

	cfg := config.Load(logger)

	db, err := config.NewPostgresConnection(&cfg.DB)
	if err != nil {
		logger.Error("DB connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if cfg.Command.IsMigrationCommand() {
		runMigrationCommand(db, logger, cfg)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- Repositories ---
	teamRepo := repositories.NewPostgresTeamRepository(db)
	userRepo := repositories.NewPostgresUserRepository(db)
	prRepo := repositories.NewPostgresPRRepository(db)

	// --- Domain Services ---
	randomizer := services.NewDefaultRandomizer()
	assignmentService := services.NewReviewerAssignmentService(randomizer)

	// --- Application Layer ---
	createTeamCmd := commands.NewCreateTeamCommand(teamRepo, userRepo)
	createPRCmd := commands.NewCreatePRCommand(teamRepo, userRepo, prRepo, assignmentService)
	mergePRCmd := commands.NewMergePRCommand(prRepo)
	reassignReviewerCmd := commands.NewReassignReviewerCommand(teamRepo, userRepo, prRepo, assignmentService)
	setUserActiveCmd := commands.NewSetUserActiveCommand(userRepo)

	getTeamQuery := queries.NewGetTeamQuery(teamRepo)
	getUserReviewsQuery := queries.NewGetUserReviewsQuery(prRepo)

	// --- HTTP API ---
	router := http.NewRouter(logger, http.RouterDeps{
		CreateTeam:       createTeamCmd,
		CreatePR:         createPRCmd,
		MergePR:          mergePRCmd,
		ReassignReviewer: reassignReviewerCmd,
		SetUserActive:    setUserActiveCmd,
		GetTeam:          getTeamQuery,
		GetUserReviews:   getUserReviewsQuery,
		UserRepo:         userRepo,
	})

	http.StartServer(ctx, logger, cfg.Server, router)
}
