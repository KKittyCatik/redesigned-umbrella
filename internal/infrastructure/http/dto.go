package http

import "time"

// Request DTOs

type LoginRequest struct {
	UserID string `json:"user_id"`
}

type CreateTeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []CreateUserRequest `json:"members"`
}

type CreateUserRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type CreatePRRequest struct {
	PRID     string `json:"pull_request_id"`
	PRName   string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type MergePRRequest struct {
	PRID string `json:"pull_request_id"`
}

type ReassignReviewerRequest struct {
	PRID          string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

// Response DTOs

type LoginResponse struct {
	Token string `json:"token"`
}

type TeamResponse struct {
	Name    string         `json:"team_name"`
	Members []UserResponse `json:"members"`
}

type UserResponse struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PRResponse struct {
	ID                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"created_at"`
	MergedAt          *time.Time `json:"merged_at,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	Status string `json:"status"`
}
