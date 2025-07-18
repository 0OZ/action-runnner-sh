package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github-runner-manager/internal/types"
)

// GetRunnerToken fetches a new runner registration token from GitHub API
func GetRunnerToken(repo types.Repository, githubToken string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runners/registration-token", repo.Owner, repo.Name)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp types.GitHubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	return tokenResp.Token, nil
}
