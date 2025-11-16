package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	DB        *sql.DB
	Container testcontainers.Container
}

func SetupTestDB(t *testing.T) *TestDB {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get connection string: %s", err)
	}

	connStr = connStr + "&sslmode=disable"

	var db *sql.DB
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		t.Fatalf("Failed to connect to test database: %s", err)
	}

	if err := runMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %s", err)
	}

	return &TestDB{
		DB:        db,
		Container: pgContainer,
	}
}

func (tdb *TestDB) Cleanup(t *testing.T) {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
	if tdb.Container != nil {
		if err := tdb.Container.Terminate(context.Background()); err != nil {
			t.Logf("Warning: failed to terminate container: %s", err)
		}
	}
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS teams (
			name VARCHAR(255) PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			team_name VARCHAR(255) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (team_name) REFERENCES teams(name) ON DELETE CASCADE
		)`,

		`CREATE TABLE IF NOT EXISTS pull_requests (
			id VARCHAR(255) PRIMARY KEY,
			pull_request_name VARCHAR(255) NOT NULL,
			author_id VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'open',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			merged_at TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		`CREATE TABLE IF NOT EXISTS pull_request_reviewers (
			pull_request_id VARCHAR(255) NOT NULL,
			reviewer_id VARCHAR(255) NOT NULL,
			assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (pull_request_id, reviewer_id),
			FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
			FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		`CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_requests_status ON pull_requests(status)`,
		`CREATE INDEX IF NOT EXISTS idx_pull_request_reviewers_reviewer_id ON pull_request_reviewers(reviewer_id)`,

		`INSERT INTO teams (name) VALUES 
			('backend-team'),
			('frontend-team'),
			('mobile-team'),
			('admin-team')
		ON CONFLICT (name) DO NOTHING`,

		`INSERT INTO users (id, username, team_name, is_active) VALUES 
		 	('test-admin', 'admin', 'admin-team', true),
			('test-user-1', 'john_backend', 'backend-team', true),
			('test-user-2', 'jane_backend', 'backend-team', true),
			('test-user-3', 'bob_frontend', 'frontend-team', true),
			('test-user-4', 'alice_mobile', 'mobile-team', true),
			('test-user-5', 'charlie_backend', 'backend-team', false)
		ON CONFLICT (id) DO NOTHING`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %v\nSQL: %s", err, migration)
		}
	}

	return nil
}
