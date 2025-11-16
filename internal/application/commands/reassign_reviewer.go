package commands

import (
	"context"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/services"
)

type ReassignReviewerCommand struct {
	teamRepo          ports.TeamRepository
	userRepo          ports.UserRepository
	prRepo            ports.PRRepository
	assignmentService *services.ReviewerAssignmentService
}

func NewReassignReviewerCommand(
	teamRepo ports.TeamRepository,
	userRepo ports.UserRepository,
	prRepo ports.PRRepository,
	assignmentService *services.ReviewerAssignmentService,
) *ReassignReviewerCommand {
	return &ReassignReviewerCommand{
		teamRepo:          teamRepo,
		userRepo:          userRepo,
		prRepo:            prRepo,
		assignmentService: assignmentService,
	}
}

type ReassignReviewerResult struct {
	PR         *entities.PullRequest
	ReplacedBy string
}

func (c *ReassignReviewerCommand) Execute(ctx context.Context, prID, oldReviewerID string) (*ReassignReviewerResult, error) {
	pr, err := c.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("getting pr: %w", err)
	}
	if pr == nil {
		return nil, entities.ErrPRNotFound
	}

	if pr.IsMerged() {
		return nil, entities.ErrPRMerged
	}

	found := false
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, entities.ErrReviewerNotAssigned
	}

	oldReviewer, err := c.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		return nil, fmt.Errorf("getting old reviewer: %w", err)
	}
	if oldReviewer == nil {
		return nil, entities.ErrUserNotFound
	}

	team, err := c.teamRepo.GetByName(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, fmt.Errorf("getting team: %w", err)
	}
	if team == nil {
		return nil, entities.ErrTeamNotFound
	}

	newReviewerID, err := c.assignmentService.FindReplacement(team, pr.AuthorID, pr.AssignedReviewers)
	if err != nil {
		return nil, fmt.Errorf("finding replacement: %w", err)
	}

	err = pr.ReassignReviewer(oldReviewerID, newReviewerID)
	if err != nil {
		return nil, fmt.Errorf("reassigning reviewer: %w", err)
	}

	err = c.prRepo.Save(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("saving pr: %w", err)
	}

	return &ReassignReviewerResult{
		PR:         pr,
		ReplacedBy: newReviewerID,
	}, nil
}
