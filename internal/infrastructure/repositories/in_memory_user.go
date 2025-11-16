package repositories

import (
	"context"
	"sync"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*entities.User
}

func NewInMemoryUserRepository() ports.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (r *InMemoryUserRepository) Save(ctx context.Context, user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.users[id], nil
}

func (r *InMemoryUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.users[id]
	return exists, nil
}

func (r *InMemoryUserRepository) GetByTeamName(ctx context.Context, teamName string) ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entities.User
	for _, user := range r.users {
		if user.TeamName == teamName {
			result = append(result, user)
		}
	}
	return result, nil
}
