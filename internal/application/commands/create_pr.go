package commands

import (
	"context"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/services"
)

type CreatePRCommand struct {
	teamRepo          ports.TeamRepository
	userRepo          ports.UserRepository
	prRepo            ports.PRRepository
	assignmentService *services.ReviewerAssignmentService
}

func NewCreatePRCommand(
	teamRepo ports.TeamRepository,
	userRepo ports.UserRepository,
	prRepo ports.PRRepository,
	assignmentService *services.ReviewerAssignmentService,
) *CreatePRCommand {
	return &CreatePRCommand{
		teamRepo:          teamRepo,
		userRepo:          userRepo,
		prRepo:            prRepo,
		assignmentService: assignmentService,
	}
}

func (c *CreatePRCommand) Execute(ctx context.Context, prID, prName, authorID string) (*entities.PullRequest, error) {
	exists, err := c.prRepo.ExistsByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("checking pr exists: %w", err)
	}
	if exists {
		return nil, entities.ErrPRExists
	}

	author, err := c.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("getting author: %w", err)
	}
	if author == nil {
		return nil, entities.ErrUserNotFound
	}

	team, err := c.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("getting team: %w", err)
	}
	if team == nil {
		return nil, entities.ErrTeamNotFound
	}

	reviewers, err := c.assignmentService.SelectReviewers(team, authorID)
	if err != nil {
		return nil, fmt.Errorf("selecting reviewers: %w", err)
	}

	pr := entities.NewPullRequest(prID, prName, authorID, reviewers)

	err = c.prRepo.Save(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("saving pr: %w", err)
	}

	return pr, nil
}
