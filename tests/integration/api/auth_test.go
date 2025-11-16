package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KKittyCatik/redesigned-umbrella/internal/bootstrap"
	"github.com/KKittyCatik/redesigned-umbrella/tests/testutils"
)

func TestLogin(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	app := bootstrap.NewTestApplication(testDB.DB)
	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	t.Run("Successful login", func(t *testing.T) {
		loginData := map[string]string{
			"user_id": "test-admin",
		}

		body, _ := json.Marshal(loginData)
		resp, err := http.Post(ts.URL+"/login", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result["token"] == "" {
			t.Error("Expected token in response")
		}
	})

	t.Run("Login with non-existent user", func(t *testing.T) {
		loginData := map[string]string{
			"user_id": "non-existent-user",
		}

		body, _ := json.Marshal(loginData)
		resp, err := http.Post(ts.URL+"/login", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404 for non-existent user, got %d", resp.StatusCode)
		}
	})
}
