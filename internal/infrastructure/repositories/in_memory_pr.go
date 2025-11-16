package repositories

import (
	"context"
	"sync"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type InMemoryPRRepository struct {
	mu  sync.RWMutex
	prs map[string]*entities.PullRequest
}

func NewInMemoryPRRepository() ports.PRRepository {
	return &InMemoryPRRepository{
		prs: make(map[string]*entities.PullRequest),
	}
}

func (r *InMemoryPRRepository) Save(ctx context.Context, pr *entities.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prs[pr.ID] = pr
	return nil
}

func (r *InMemoryPRRepository) GetByID(ctx context.Context, id string) (*entities.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.prs[id], nil
}

func (r *InMemoryPRRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.prs[id]
	return exists, nil
}

func (r *InMemoryPRRepository) GetByReviewer(ctx context.Context, userID string) ([]*entities.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entities.PullRequest
	for _, pr := range r.prs {
		for _, reviewer := range pr.AssignedReviewers {
			if reviewer == userID {
				result = append(result, pr)
				break
			}
		}
	}
	return result, nil
}
