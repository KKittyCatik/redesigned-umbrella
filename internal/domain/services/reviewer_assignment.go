package services

import (
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type ReviewerAssignmentService struct {
	rnd Randomizer
}

func NewReviewerAssignmentService(rnd Randomizer) *ReviewerAssignmentService {
	if rnd == nil {
		rnd = NewDefaultRandomizer()
	}
	return &ReviewerAssignmentService{rnd: rnd}
}

func (s *ReviewerAssignmentService) SelectReviewers(team *entities.Team, authorID string) ([]string, error) {
	active := team.GetActiveMembers()
	candidates := make([]string, 0, len(active))
	for _, u := range active {
		if u.ID == authorID {
			continue
		}
		candidates = append(candidates, u.ID)
	}

	if len(candidates) == 0 {
		return nil, entities.ErrNoCandidateFound
	}

	perm := s.rnd.Perm(len(candidates))
	limit := 2

	if len(candidates) < limit {
		limit = len(candidates)
	}

	result := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, candidates[perm[i]])
	}
	return result, nil
}

func (s *ReviewerAssignmentService) FindReplacement(team *entities.Team, authorID string, currentReviewers []string) (string, error) {
	active := team.GetActiveMembers()
	current := make(map[string]struct{}, len(currentReviewers))
	for _, id := range currentReviewers {
		current[id] = struct{}{}
	}

	candidates := make([]string, 0)
	for _, u := range active {
		if u.ID == authorID {
			continue
		}
		if _, ok := current[u.ID]; ok {
			continue
		}
		candidates = append(candidates, u.ID)
	}

	if len(candidates) == 0 {
		return "", entities.ErrNoCandidateFound
	}

	return candidates[s.rnd.Intn(len(candidates))], nil
}
