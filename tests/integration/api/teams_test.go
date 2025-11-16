package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KKittyCatik/redesigned-umbrella/internal/bootstrap"
	"github.com/KKittyCatik/redesigned-umbrella/tests/testutils"
)

func TestTeamCRUD(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	app := bootstrap.NewTestApplication(testDB.DB)
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	t.Run("Create team successfully", func(t *testing.T) {
		createTestUsers(t, testDB.DB)

		token := getAuthToken(t, ts.URL)

		teamData := map[string]interface{}{
			"team_name": "integration-test-team",
			"members": []map[string]interface{}{
				{
					"user_id":   "test-user-1",
					"username":  "john_backend",
					"is_active": true,
				},
				{
					"user_id":   "test-user-2",
					"username":  "jane_backend",
					"is_active": true,
				},
			},
		}

		body, _ := json.Marshal(teamData)
		req, _ := http.NewRequest("POST", ts.URL+"/team/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		if response["team_name"] != "integration-test-team" {
			t.Errorf("Expected team name 'integration-test-team', got %v", response["team_name"])
		}
	})
}

func createTestUsers(t *testing.T, db *sql.DB) {
	users := []struct {
		id       string
		username string
	}{
		{"test-user-1", "testuser1"},
		{"test-user-2", "testuser2"},
		{"test-user-3", "testuser3"},
	}

	for _, user := range users {
		_, err := db.Exec(
			"INSERT INTO users (id, username, team_name, is_active, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING",
			user.id, user.username, "default-team", true, "2023-01-01 00:00:00",
		)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}
}

func getAuthToken(t *testing.T, baseURL string) string {
	loginData := map[string]string{
		"user_id": "test-admin",
	}

	body, _ := json.Marshal(loginData)
	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed with status: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	return result["token"]
}
