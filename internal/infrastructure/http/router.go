package http

import (
	"log/slog"
	"net/http"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/commands"
	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/application/queries"
)

type RouterDeps struct {
	CreateTeam       *commands.CreateTeamCommand
	CreatePR         *commands.CreatePRCommand
	MergePR          *commands.MergePRCommand
	ReassignReviewer *commands.ReassignReviewerCommand
	SetUserActive    *commands.SetUserActiveCommand
	GetTeam          *queries.GetTeamQuery
	GetUserReviews   *queries.GetUserReviewsQuery
	UserRepo         ports.UserRepository
}

func NewRouter(logger *slog.Logger, deps RouterDeps) http.Handler {
	handler := NewHandler(
		deps.CreateTeam,
		deps.CreatePR,
		deps.MergePR,
		deps.ReassignReviewer,
		deps.SetUserActive,
		deps.GetTeam,
		deps.GetUserReviews,
		deps.UserRepo,
		logger,
	)

	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("POST /login", handler.Login)
	mux.HandleFunc("GET /health", handler.Health)

	// Protected endpoints
	mux.HandleFunc("POST /team/add", AuthMiddleware(logger, handler.CreateTeam))
	mux.HandleFunc("GET /team/get", AuthMiddleware(logger, handler.GetTeam))
	mux.HandleFunc("POST /users/setIsActive", AuthMiddleware(logger, handler.SetUserActive))
	mux.HandleFunc("POST /pullRequest/create", AuthMiddleware(logger, handler.CreatePR))
	mux.HandleFunc("POST /pullRequest/merge", AuthMiddleware(logger, handler.MergePR))
	mux.HandleFunc("POST /pullRequest/reassign", AuthMiddleware(logger, handler.ReassignReviewer))
	mux.HandleFunc("GET /users/getReview", AuthMiddleware(logger, handler.GetUserReviews))

	return mux
}
