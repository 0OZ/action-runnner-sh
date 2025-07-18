package types

import "time"

// Runner represents a single GitHub Actions runner
type Runner struct {
	ServiceName string
	RunnerName  string
	WorkDir     string
}

// Repository represents a GitHub repository with its runners
type Repository struct {
	Owner      string
	Name       string
	FullName   string
	URL        string
	Token      string
	AnchorName string
	Runners    []Runner
}

// Config holds the entire configuration
type Config struct {
	Repositories []Repository
}

// GitHubTokenResponse represents the API response for runner registration token
type GitHubTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
