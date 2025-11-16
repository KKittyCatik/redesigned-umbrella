package ports

import (
	"context"

	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type TeamRepository interface {
	Save(ctx context.Context, team *entities.Team) error
	GetByName(ctx context.Context, name string) (*entities.Team, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}
