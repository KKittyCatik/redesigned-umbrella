package repositories

import (
	"context"
	"sync"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type InMemoryTeamRepository struct {
	mu    sync.RWMutex
	teams map[string]*entities.Team
}

func NewInMemoryTeamRepository() ports.TeamRepository {
	return &InMemoryTeamRepository{
		teams: make(map[string]*entities.Team),
	}
}

func (r *InMemoryTeamRepository) Save(ctx context.Context, team *entities.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.teams[team.Name] = team
	return nil
}

func (r *InMemoryTeamRepository) GetByName(ctx context.Context, name string) (*entities.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.teams[name], nil
}

func (r *InMemoryTeamRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.teams[name]
	return exists, nil
}
