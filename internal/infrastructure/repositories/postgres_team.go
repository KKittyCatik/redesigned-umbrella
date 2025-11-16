package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KKittyCatik/redesigned-umbrella/internal/application/ports"
	"github.com/KKittyCatik/redesigned-umbrella/internal/domain/entities"
)

type PostgresTeamRepository struct {
	db *sql.DB
}

func NewPostgresTeamRepository(db *sql.DB) ports.TeamRepository {
	return &PostgresTeamRepository{db: db}
}

func (r *PostgresTeamRepository) Save(ctx context.Context, team *entities.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
        INSERT INTO teams (name) VALUES ($1)
        ON CONFLICT DO NOTHING
    `, team.Name)
	if err != nil {
		return fmt.Errorf("insert team: %w", err)
	}

	for _, member := range team.Members {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO users (id, username, team_name, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE SET is_active = EXCLUDED.is_active
        `, member.ID, member.Username, member.TeamName, member.IsActive)
		if err != nil {
			return fmt.Errorf("insert user %s: %w", member.ID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresTeamRepository) GetByName(ctx context.Context, name string) (*entities.Team, error) {
	team := &entities.Team{
		Name:    name,
		Members: make([]*entities.User, 0),
	}

	rows, err := r.db.QueryContext(ctx, `
        SELECT id, username, team_name, is_active 
        FROM users 
        WHERE team_name = $1
    `, name)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, username, teamName string
		var isActive bool
		if err := rows.Scan(&id, &username, &teamName, &isActive); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		user := entities.NewUser(id, username, teamName, isActive)
		team.Members = append(team.Members, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(team.Members) == 0 {
		return nil, nil
	}

	return team, nil
}

func (r *PostgresTeamRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)
    `, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query team exists: %w", err)
	}
	return exists, nil
}
