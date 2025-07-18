package utils

import (
	"fmt"
	"github-runner-manager/internal/types"
	"strings"
)

// ParseRepo parses a repository string like "github.com/owner/repo" or "owner/repo"
func ParseRepo(repoStr string) (owner, name string, err error) {
	repoStr = strings.TrimPrefix(repoStr, "github.com/")
	repoStr = strings.TrimPrefix(repoStr, "https://github.com/")

	parts := strings.Split(repoStr, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s (expected owner/repo)", repoStr)
	}

	return parts[0], parts[1], nil
}

// GenerateAnchorName creates a unique anchor name for a repository
func GenerateAnchorName(repoName string) string {
	// Simple conversion: replace non-alphanumeric with dash and append -api-env
	name := strings.ToLower(repoName)
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return fmt.Sprintf("%s-api-env", name)
}

// RepoExists checks if a repository already exists in the configuration
func RepoExists(config types.Config, owner, name string) bool {
	for _, repo := range config.Repositories {
		if repo.Owner == owner && repo.Name == name {
			return true
		}
	}
	return false
}

// CreateNewRepository creates a new repository with the specified number of runners
func CreateNewRepository(owner, name string, numRunners int) types.Repository {
	newRepo := types.Repository{
		Owner:      owner,
		Name:       name,
		AnchorName: GenerateAnchorName(name),
		Runners:    make([]types.Runner, 0, numRunners),
	}

	// Generate runners
	for i := 1; i <= numRunners; i++ {
		runner := types.Runner{
			ServiceName: fmt.Sprintf("%s-runner-%d", name, i),
			RunnerName:  fmt.Sprintf("%s-%d", name, i),
			WorkDir:     fmt.Sprintf("/tmp/runner/%s-%d", name, i),
		}
		newRepo.Runners = append(newRepo.Runners, runner)
	}

	return newRepo
}

// UpdateRepositoryInfo updates repository information (FullName and URL)
func UpdateRepositoryInfo(repo *types.Repository) {
	repo.FullName = fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
	repo.URL = fmt.Sprintf("https://github.com/%s/%s", repo.Owner, repo.Name)
}

// GetTotalRunnerCount returns the total number of runners across all repositories
func GetTotalRunnerCount(config types.Config) int {
	total := 0
	for _, repo := range config.Repositories {
		total += len(repo.Runners)
	}
	return total
}
