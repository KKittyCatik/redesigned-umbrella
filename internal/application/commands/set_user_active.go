package commands

import (
	"context"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type SetUserActiveCommand struct {
	userRepo ports.UserRepository
}

func NewSetUserActiveCommand(userRepo ports.UserRepository) *SetUserActiveCommand {
	return &SetUserActiveCommand{userRepo: userRepo}
}

func (c *SetUserActiveCommand) Execute(ctx context.Context, userID string, isActive bool) (*entities.User, error) {
	user, err := c.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	if user == nil {
		return nil, entities.ErrUserNotFound
	}

	user.SetActive(isActive)

	err = c.userRepo.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("saving user: %w", err)
	}

	return user, nil
}
