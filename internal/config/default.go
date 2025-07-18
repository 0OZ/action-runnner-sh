package config

import (
	"fmt"
	"github-runner-manager/internal/types"
	"log"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// RepoConfig defines the configuration for generating a repository
type RepoConfig struct {
	Owner       string
	Name        string
	AnchorName  string
	RunnerCount int
}

// getDefaultConfig returns a dynamically generated configuration
func getDefaultConfig() types.Config {
	// Define your repositories and how many runners each should have
	repoConfigs := []RepoConfig{
		{
			Owner:       "0OZ",
			Name:        "monorepo-frontend",
			AnchorName:  "f-api-env",
			RunnerCount: 3,
		},
		{
			Owner:       "0OZ",
			Name:        "auftrag-select-backend",
			AnchorName:  "b-api-env",
			RunnerCount: 3,
		},
		{
			Owner:       "0OZ",
			Name:        "Toolbox",
			AnchorName:  "toolbox-api-env",
			RunnerCount: 1,
		},
		{
			Owner:       "0OZ",
			Name:        "auftrag-ai-frontend",
			AnchorName:  "aai-fd",
			RunnerCount: 2,
		},
	}

	return generateConfig(repoConfigs)
}

// generateConfig creates a dynamic configuration based on repo configs
func generateConfig(repoConfigs []RepoConfig) types.Config {
	config := types.Config{
		Repositories: make([]types.Repository, len(repoConfigs)),
	}

	usedDirs := make(map[string]bool)
	runnerCounter := 1

	for i, rc := range repoConfigs {
		repo := types.Repository{
			Owner:      rc.Owner,
			Name:       rc.Name,
			AnchorName: rc.AnchorName,
			Runners:    make([]types.Runner, rc.RunnerCount),
		}

		for j := 0; j < rc.RunnerCount; j++ {
			runner := types.Runner{
				ServiceName: fmt.Sprintf("service-%d", runnerCounter),
				RunnerName:  fmt.Sprintf("runner-%d", runnerCounter),
				WorkDir:     generateUniqueWorkDir(usedDirs),
			}
			repo.Runners[j] = runner
			runnerCounter++
		}

		config.Repositories[i] = repo
	}

	return config
}

// generateUniqueWorkDir creates a unique work directory under /tmp
func generateUniqueWorkDir(usedDirs map[string]bool) string {
	for {
		var id, err = gonanoid.New()
		if err != nil {
			log.Fatalf("Failed to generate ID: %v", err)
		}
		dirName := fmt.Sprintf("/tmp/runner/%s", id)

		// Check if it's already used
		if !usedDirs[dirName] {
			usedDirs[dirName] = true
			return dirName
		}
	}
}
