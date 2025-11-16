package queries

import (
	"context"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type GetUserReviewsQuery struct {
	prRepo ports.PRRepository
}

func NewGetUserReviewsQuery(prRepo ports.PRRepository) *GetUserReviewsQuery {
	return &GetUserReviewsQuery{prRepo: prRepo}
}

func (q *GetUserReviewsQuery) Execute(ctx context.Context, userID string) ([]*entities.PullRequest, error) {
	prs, err := q.prRepo.GetByReviewer(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting prs by reviewer: %w", err)
	}
	return prs, nil
}
