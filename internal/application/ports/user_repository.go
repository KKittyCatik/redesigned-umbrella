package ports

import (
	"context"

	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	GetByTeamName(ctx context.Context, teamName string) ([]*entities.User, error)
}
