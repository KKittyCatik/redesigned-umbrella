package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type PostgresPRRepository struct {
	db *sql.DB
}

func NewPostgresPRRepository(db *sql.DB) ports.PRRepository {
	return &PostgresPRRepository{db: db}
}

func (r *PostgresPRRepository) Save(ctx context.Context, pr *entities.PullRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
        INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at) 
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE SET 
            status = EXCLUDED.status,
            merged_at = EXCLUDED.merged_at
    `, pr.ID, pr.Name, pr.AuthorID, pr.Status.String(), pr.CreatedAt, pr.MergedAt)
	if err != nil {
		return fmt.Errorf("insert pr: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
        DELETE FROM pull_request_reviewers WHERE pull_request_id = $1
    `, pr.ID)
	if err != nil {
		return fmt.Errorf("delete reviewers: %w", err)
	}

	for _, reviewer := range pr.AssignedReviewers {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id) 
            VALUES ($1, $2)
        `, pr.ID, reviewer)
		if err != nil {
			return fmt.Errorf("insert reviewer: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresPRRepository) GetByID(ctx context.Context, id string) (*entities.PullRequest, error) {
	var prID, name, authorID, statusStr string
	var createdAt time.Time
	var mergedAt *time.Time

	err := r.db.QueryRowContext(ctx, `
        SELECT id, name, author_id, status, created_at, merged_at 
        FROM pull_requests 
        WHERE id = $1
    `, id).Scan(&prID, &name, &authorID, &statusStr, &createdAt, &mergedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query pr: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
        SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id = $1
    `, id)
	if err != nil {
		return nil, fmt.Errorf("query reviewers: %w", err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("scan reviewer: %w", err)
		}
		reviewers = append(reviewers, reviewerID)
	}

	status, err := entities.ParsePRStatus(statusStr)
	if err != nil {
		return nil, err
	}

	pr := &entities.PullRequest{
		ID:                prID,
		Name:              name,
		AuthorID:          authorID,
		Status:            status,
		AssignedReviewers: reviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}

	return pr, nil
}

func (r *PostgresPRRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1)
    `, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query pr exists: %w", err)
	}
	return exists, nil
}

func (r *PostgresPRRepository) GetByReviewer(ctx context.Context, userID string) ([]*entities.PullRequest, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT pr.id, pr.name, pr.author_id, pr.status, pr.created_at, pr.merged_at
        FROM pull_requests pr
        JOIN pull_request_reviewers prr ON pr.id = prr.pull_request_id
        WHERE prr.reviewer_id = $1
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("query prs by reviewer: %w", err)
	}
	defer rows.Close()

	var prs []*entities.PullRequest
	for rows.Next() {
		var id, name, authorID, statusStr string
		var createdAt time.Time
		var mergedAt *time.Time

		if err := rows.Scan(&id, &name, &authorID, &statusStr, &createdAt, &mergedAt); err != nil {
			return nil, fmt.Errorf("scan pr: %w", err)
		}

		status, err := entities.ParsePRStatus(statusStr)
		if err != nil {
			return nil, err
		}

		pr := &entities.PullRequest{
			ID:        id,
			Name:      name,
			AuthorID:  authorID,
			Status:    status,
			CreatedAt: createdAt,
			MergedAt:  mergedAt,
		}
		prs = append(prs, pr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return prs, nil
}
