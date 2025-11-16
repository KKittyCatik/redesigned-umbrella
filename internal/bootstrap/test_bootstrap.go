package bootstrap

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/commands"
	"github.com/KKittyCatik/redesigned-umbrella/internal/application/queries"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/services"
	apphttp "github.com/KKittyCatik/redesigned-umbrella/internal/infrastructure/http"
	"github.com/KKittyCatik/redesigned-umbrella/internal/infrastructure/repositories"
)

type TestApplication struct {
	DB     *sql.DB
	Router http.Handler
}

func NewTestApplication(db *sql.DB) *TestApplication {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	teamRepo := repositories.NewPostgresTeamRepository(db)
	userRepo := repositories.NewPostgresUserRepository(db)
	prRepo := repositories.NewPostgresPRRepository(db)

	randomizer := services.NewDefaultRandomizer()
	assignmentService := services.NewReviewerAssignmentService(randomizer)

	createTeamCmd := commands.NewCreateTeamCommand(teamRepo, userRepo)
	createPRCmd := commands.NewCreatePRCommand(teamRepo, userRepo, prRepo, assignmentService)
	mergePRCmd := commands.NewMergePRCommand(prRepo)
	reassignReviewerCmd := commands.NewReassignReviewerCommand(teamRepo, userRepo, prRepo, assignmentService)
	setUserActiveCmd := commands.NewSetUserActiveCommand(userRepo)

	getTeamQuery := queries.NewGetTeamQuery(teamRepo)
	getUserReviewsQuery := queries.NewGetUserReviewsQuery(prRepo)

	router := apphttp.NewRouter(logger, apphttp.RouterDeps{
		CreateTeam:       createTeamCmd,
		CreatePR:         createPRCmd,
		MergePR:          mergePRCmd,
		ReassignReviewer: reassignReviewerCmd,
		SetUserActive:    setUserActiveCmd,
		GetTeam:          getTeamQuery,
		GetUserReviews:   getUserReviewsQuery,
		UserRepo:         userRepo,
	})

	return &TestApplication{
		DB:     db,
		Router: router,
	}
}
