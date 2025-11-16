package http

import (
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

func MapTeamToResponse(team *entities.Team) TeamResponse {
	members := make([]UserResponse, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, MapUserToResponse(member))
	}
	return TeamResponse{
		Name:    team.Name,
		Members: members,
	}
}

func MapUserToResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func MapPRToResponse(pr *entities.PullRequest) PRResponse {
	reviewers := pr.AssignedReviewers
	if reviewers == nil {
		reviewers = []string{}
	}
	return PRResponse{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status.String(),
		AssignedReviewers: reviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func MapCreateTeamRequestToUsers(req CreateTeamRequest) []*entities.User {
	users := make([]*entities.User, 0, len(req.Members))
	for _, member := range req.Members {
		user := entities.NewUser(member.UserID, member.Username, req.TeamName, member.IsActive)
		users = append(users, user)
	}
	return users
}
