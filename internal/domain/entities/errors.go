package entities

import "errors"

var (
	// PullRequest errors
	ErrPRMerged            = errors.New("pull request is already merged")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this pull request")

	// Team errors
	ErrTeamNotFound     = errors.New("team not found")
	ErrTeamExists       = errors.New("team already exists")
	ErrNoCandidateFound = errors.New("no active replacement candidate found")
	ErrMemberExists     = errors.New("member already exists")

	// User errors
	ErrUserNotFound = errors.New("user not found")

	// PR errors
	ErrPRExists        = errors.New("pull request already exists")
	ErrPRNotFound      = errors.New("pull request not found")
	ErrInvalidPRStatus = errors.New("invalid pull request status")
)
