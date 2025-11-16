package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/commands"
	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/application/queries"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type Handler struct {
	// Commands
	createTeamCmd       *commands.CreateTeamCommand
	createPRCmd         *commands.CreatePRCommand
	mergePRCmd          *commands.MergePRCommand
	reassignReviewerCmd *commands.ReassignReviewerCommand
	setUserActiveCmd    *commands.SetUserActiveCommand

	// Queries
	getTeamQuery        *queries.GetTeamQuery
	getUserReviewsQuery *queries.GetUserReviewsQuery

	// Repository
	userRepo ports.UserRepository

	logger *slog.Logger
}

func NewHandler(
	createTeamCmd *commands.CreateTeamCommand,
	createPRCmd *commands.CreatePRCommand,
	mergePRCmd *commands.MergePRCommand,
	reassignReviewerCmd *commands.ReassignReviewerCommand,
	setUserActiveCmd *commands.SetUserActiveCommand,
	getTeamQuery *queries.GetTeamQuery,
	getUserReviewsQuery *queries.GetUserReviewsQuery,
	userRepo ports.UserRepository,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		createTeamCmd:       createTeamCmd,
		createPRCmd:         createPRCmd,
		mergePRCmd:          mergePRCmd,
		reassignReviewerCmd: reassignReviewerCmd,
		setUserActiveCmd:    setUserActiveCmd,
		getTeamQuery:        getTeamQuery,
		getUserReviewsQuery: getUserReviewsQuery,
		userRepo:            userRepo,
		logger:              logger,
	}
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.TeamName == "" {
		h.logger.Error("validation error", "error", "team_name is empty")
		h.respondWithError(w, http.StatusBadRequest, "team_name cannot be empty")
		return
	}

	if len(req.Members) == 0 {
		h.logger.Error("validation error", "error", "members is empty")
		h.respondWithError(w, http.StatusBadRequest, "members cannot be empty")
		return
	}

	for _, member := range req.Members {
		if member.UserID == "" || member.Username == "" {
			h.logger.Error("validation error", "error", "user_id or username is empty")
			h.respondWithError(w, http.StatusBadRequest, "user_id and username cannot be empty")
			return
		}
	}

	members := MapCreateTeamRequestToUsers(req)

	team, err := h.createTeamCmd.Execute(r.Context(), req.TeamName, members)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(MapTeamToResponse(team))
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		h.logger.Error("validation error", "error", "team_name is empty")
		h.respondWithError(w, http.StatusBadRequest, "team_name is required")
		return
	}

	team, err := h.getTeamQuery.Execute(r.Context(), teamName)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MapTeamToResponse(team))
}

func (h *Handler) SetUserActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.UserID == "" {
		h.logger.Error("validation error", "error", "user_id is empty")
		h.respondWithError(w, http.StatusBadRequest, "user_id cannot be empty")
		return
	}

	user, err := h.setUserActiveCmd.Execute(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MapUserToResponse(user))
}

func (h *Handler) CreatePR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.PRID == "" || req.PRName == "" || req.AuthorID == "" {
		h.logger.Error("validation error", "error", "missing required fields")
		h.respondWithError(w, http.StatusBadRequest, "pull_request_id, name, and author_id are required")
		return
	}

	pr, err := h.createPRCmd.Execute(r.Context(), req.PRID, req.PRName, req.AuthorID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(MapPRToResponse(pr))
}

func (h *Handler) MergePR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.PRID == "" {
		h.logger.Error("validation error", "error", "pull_request_id is empty")
		h.respondWithError(w, http.StatusBadRequest, "pull_request_id is required")
		return
	}

	pr, err := h.mergePRCmd.Execute(r.Context(), req.PRID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MapPRToResponse(pr))
}

func (h *Handler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ReassignReviewerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.PRID == "" || req.OldReviewerID == "" {
		h.logger.Error("validation error", "error", "missing required fields")
		h.respondWithError(w, http.StatusBadRequest, "pull_request_id and old_reviewer_id are required")
		return
	}

	result, err := h.reassignReviewerCmd.Execute(r.Context(), req.PRID, req.OldReviewerID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pull_request": MapPRToResponse(result.PR),
		"replaced_by":  result.ReplacedBy,
	})
}

func (h *Handler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.logger.Error("validation error", "error", "user_id is empty")
		h.respondWithError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	prs, err := h.getUserReviewsQuery.Execute(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	responses := make([]PRResponse, 0, len(prs))
	for _, pr := range prs {
		responses = append(responses, MapPRToResponse(pr))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pull_requests": responses,
	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.UserID == "" {
		h.logger.Error("validation error", "error", "user_id is empty")
		h.respondWithError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), req.UserID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	if user == nil {
		h.respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	token, err := GenerateToken(req.UserID)
	if err != nil {
		h.logger.Error("failed to generate token", "error", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func (h *Handler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	h.logger.Error("error occurred", "error", err)

	switch err {
	case entities.ErrTeamExists:
		h.respondWithError(w, http.StatusConflict, "Team already exists")
	case entities.ErrTeamNotFound:
		h.respondWithError(w, http.StatusNotFound, "Team not found")
	case entities.ErrUserNotFound:
		h.respondWithError(w, http.StatusNotFound, "User not found")
	case entities.ErrPRExists:
		h.respondWithError(w, http.StatusConflict, "Pull request already exists")
	case entities.ErrPRNotFound:
		h.respondWithError(w, http.StatusNotFound, "Pull request not found")
	case entities.ErrPRMerged:
		h.respondWithError(w, http.StatusConflict, "Pull request is already merged")
	case entities.ErrReviewerNotAssigned:
		h.respondWithError(w, http.StatusBadRequest, "Reviewer is not assigned to this pull request")
	case entities.ErrNoCandidateFound:
		h.respondWithError(w, http.StatusConflict, "No active reviewer available")
	case entities.ErrMemberExists:
		h.respondWithError(w, http.StatusConflict, "Member already exists in team")
	default:
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
