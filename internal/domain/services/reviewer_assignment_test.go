package services

import (
	"testing"

	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type MockRandomizer struct {
	permResult []int
	intnResult int
}

func (m *MockRandomizer) Perm(n int) []int {
	return m.permResult
}

func (m *MockRandomizer) Intn(n int) int {
	return m.intnResult
}

func TestSelectReviewers(t *testing.T) {
	tests := []struct {
		name           string
		team           *entities.Team
		authorID       string
		mockRandomizer *MockRandomizer
		expectedCount  int
		expectError    bool
	}{
		{
			name: "Select 2 reviewers from 3 active members (excluding author)",
			team: &entities.Team{
				Name: "Backend",
				Members: []*entities.User{
					entities.NewUser("user1", "Alice", "Backend", true),
					entities.NewUser("user2", "Bob", "Backend", true),
					entities.NewUser("user3", "Charlie", "Backend", true),
				},
			},
			authorID: "user1",
			mockRandomizer: &MockRandomizer{
				permResult: []int{0, 1, 2},
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "Select 1 reviewer from 2 active members (excluding author)",
			team: &entities.Team{
				Name: "Backend",
				Members: []*entities.User{
					entities.NewUser("user1", "Alice", "Backend", true),
					entities.NewUser("user2", "Bob", "Backend", true),
				},
			},
			authorID: "user1",
			mockRandomizer: &MockRandomizer{
				permResult: []int{0, 1},
			},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "No active members except author",
			team: &entities.Team{
				Name: "Backend",
				Members: []*entities.User{
					entities.NewUser("user1", "Alice", "Backend", true),
					entities.NewUser("user2", "Bob", "Backend", false),
				},
			},
			authorID:       "user1",
			mockRandomizer: &MockRandomizer{},
			expectedCount:  0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewReviewerAssignmentService(tt.mockRandomizer)

			reviewers, err := service.SelectReviewers(tt.team, tt.authorID)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(reviewers) != tt.expectedCount {
				t.Errorf("expected %d reviewers, got %d", tt.expectedCount, len(reviewers))
			}

			for _, reviewer := range reviewers {
				if reviewer == tt.authorID {
					t.Error("author should not be in reviewers list")
				}
			}
		})
	}
}

func TestFindReplacement(t *testing.T) {
	tests := []struct {
		name             string
		team             *entities.Team
		authorID         string
		currentReviewers []string
		mockRandomizer   *MockRandomizer
		expectError      bool
	}{
		{
			name: "Find replacement from 3 active members (excluding author and current reviewers)",
			team: &entities.Team{
				Name: "Backend",
				Members: []*entities.User{
					entities.NewUser("user1", "Alice", "Backend", true),
					entities.NewUser("user2", "Bob", "Backend", true),
					entities.NewUser("user3", "Charlie", "Backend", true),
				},
			},
			authorID:         "user1",
			currentReviewers: []string{"user2"},
			mockRandomizer: &MockRandomizer{
				intnResult: 0,
			},
			expectError: false,
		},
		{
			name: "No replacement available",
			team: &entities.Team{
				Name: "Backend",
				Members: []*entities.User{
					entities.NewUser("user1", "Alice", "Backend", true),
					entities.NewUser("user2", "Bob", "Backend", false),
				},
			},
			authorID:         "user1",
			currentReviewers: []string{"user2"},
			mockRandomizer:   &MockRandomizer{},
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewReviewerAssignmentService(tt.mockRandomizer)

			replacement, err := service.FindReplacement(tt.team, tt.authorID, tt.currentReviewers)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError {
				if replacement == tt.authorID {
					t.Error("replacement should not be author")
				}

				for _, reviewer := range tt.currentReviewers {
					if replacement == reviewer {
						t.Error("replacement should not be in current reviewers")
					}
				}
			}
		})
	}
}
