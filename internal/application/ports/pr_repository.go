package ports

import (
	"context"

	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type PRRepository interface {
	Save(ctx context.Context, pr *entities.PullRequest) error
	GetByID(ctx context.Context, id string) (*entities.PullRequest, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	GetByReviewer(ctx context.Context, userID string) ([]*entities.PullRequest, error)
}
