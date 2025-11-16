package queries

import (
	"context"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type GetTeamQuery struct {
	teamRepo ports.TeamRepository
}

func NewGetTeamQuery(teamRepo ports.TeamRepository) *GetTeamQuery {
	return &GetTeamQuery{teamRepo: teamRepo}
}

func (q *GetTeamQuery) Execute(ctx context.Context, teamName string) (*entities.Team, error) {
	team, err := q.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("getting team: %w", err)
	}
	if team == nil {
		return nil, entities.ErrTeamNotFound
	}
	return team, nil
}
