package commands

import (
	"context"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type MergePRCommand struct {
	prRepo ports.PRRepository
}

func NewMergePRCommand(prRepo ports.PRRepository) *MergePRCommand {
	return &MergePRCommand{prRepo: prRepo}
}

func (c *MergePRCommand) Execute(ctx context.Context, prID string) (*entities.PullRequest, error) {
	pr, err := c.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("getting pr: %w", err)
	}
	if pr == nil {
		return nil, entities.ErrPRNotFound
	}

	if pr.IsMerged() {
		return pr, nil
	}

	pr.Merge()

	err = c.prRepo.Save(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("saving pr: %w", err)
	}

	return pr, nil
}
