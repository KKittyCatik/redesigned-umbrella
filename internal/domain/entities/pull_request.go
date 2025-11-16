package entities

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

func (s PRStatus) String() string {
	return string(s)
}

func (s PRStatus) IsValid() bool {
	return s == PRStatusOpen || s == PRStatusMerged
}

func ParsePRStatus(s string) (PRStatus, error) {
	status := PRStatus(s)
	if !status.IsValid() {
		return "", ErrInvalidPRStatus
	}
	return status, nil
}

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}

func NewPullRequest(id, name, authorID string, assignedReviewers []string) *PullRequest {
	if id == "" || name == "" || authorID == "" {
		return nil
	}

	for _, reviewer := range assignedReviewers {
		if reviewer == authorID {
			assignedReviewers = removeElement(assignedReviewers, authorID)
			break
		}
	}

	if len(assignedReviewers) > 2 {
		assignedReviewers = assignedReviewers[:2]
	}

	return &PullRequest{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            PRStatusOpen,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         time.Now(),
		MergedAt:          nil,
	}
}

func removeElement(slice []string, elem string) []string {
	result := []string{}
	for _, v := range slice {
		if v != elem {
			result = append(result, v)
		}
	}
	return result
}

func (pr *PullRequest) Merge() {
	if pr.Status != PRStatusMerged {
		pr.Status = PRStatusMerged
		now := time.Now()
		pr.MergedAt = &now
	}
}

func (pr *PullRequest) IsMerged() bool {
	return pr.Status == PRStatusMerged
}

func (pr *PullRequest) ReassignReviewer(oldReviewerID, newReviewerID string) error {
	if pr.IsMerged() {
		return ErrPRMerged
	}

	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldReviewerID {
			pr.AssignedReviewers[i] = newReviewerID
			return nil
		}
	}

	return ErrReviewerNotAssigned
}
