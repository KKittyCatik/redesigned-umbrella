package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) ports.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *entities.User) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO users (id, username, team_name, is_active) 
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET 
            username = EXCLUDED.username,
            is_active = EXCLUDED.is_active
    `, user.ID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var userID, username, teamName string
	var isActive bool

	err := r.db.QueryRowContext(ctx, `
        SELECT id, username, team_name, is_active 
        FROM users 
        WHERE id = $1
    `, id).Scan(&userID, &username, &teamName, &isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	return entities.NewUser(userID, username, teamName, isActive), nil
}

func (r *PostgresUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)
    `, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query user exists: %w", err)
	}
	return exists, nil
}

func (r *PostgresUserRepository) GetByTeamName(ctx context.Context, teamName string) ([]*entities.User, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, username, team_name, is_active 
        FROM users 
        WHERE team_name = $1
    `, teamName)
	if err != nil {
		return nil, fmt.Errorf("query users by team: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var id, username, team string
		var isActive bool
		if err := rows.Scan(&id, &username, &team, &isActive); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, entities.NewUser(id, username, team, isActive))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}
