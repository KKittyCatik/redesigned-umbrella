package commands

import (
	"context"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type CreateTeamCommand struct {
	teamRepo ports.TeamRepository
	userRepo ports.UserRepository
}

func NewCreateTeamCommand(teamRepo ports.TeamRepository, userRepo ports.UserRepository) *CreateTeamCommand {
	return &CreateTeamCommand{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (c *CreateTeamCommand) Execute(ctx context.Context, teamName string, members []*entities.User) (*entities.Team, error) {
	exists, err := c.teamRepo.ExistsByName(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("checking team exists: %w", err)
	}
	if exists {
		return nil, entities.ErrTeamExists
	}

	team := entities.NewTeam(teamName, members)

	err = c.teamRepo.Save(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("saving team: %w", err)
	}

	for _, user := range members {
		err = c.userRepo.Save(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("saving user %s: %w", user.ID, err)
		}
	}

	return team, nil
}
