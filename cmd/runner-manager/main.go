package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github-runner-manager/internal/config"
	"github-runner-manager/internal/docker"
	"github-runner-manager/internal/github"
	"github-runner-manager/internal/types"
	"github-runner-manager/internal/utils"
)

func main() {
	var (
		githubToken string
		repos       string
		configFile  string
		numRunners  int
		skipDocker  bool
	)

	flag.StringVar(&githubToken, "token", "", "GitHub Personal Access Token (required)")
	flag.StringVar(&repos, "repos", "", "Add repositories (format: owner/repo-1,owner/repo-2 or github.com/owner/repo-1,github.com/owner/repo-2)")
	flag.StringVar(&configFile, "config", "runners-config.json", "Configuration file path")
	flag.IntVar(&numRunners, "runners", 2, "Number of runners to create for each new repository")
	flag.BoolVar(&skipDocker, "skip-docker", false, "Skip Docker operations (only generate compose file)")
	flag.Parse()

	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			log.Fatal("Error: GitHub Personal Access Token required\nUsage: ./runner-manager -token <github_token> or set GITHUB_TOKEN env var")
		}
	}

	log.Println("===== Starting GitHub Runner Manager =====")

	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Handle adding new repositories
	if repos != "" {
		if err := handleAddRepositories(&cfg, repos, numRunners, configFile); err != nil {
			log.Fatalf("Error adding repositories: %v", err)
		}
	}

	// Update repository information and fetch tokens
	if err := updateRepositoriesAndFetchTokens(&cfg, githubToken); err != nil {
		log.Fatalf("Error updating repositories: %v", err)
	}

	// Generate and deploy docker-compose configuration
	if err := generateAndDeployDockerCompose(cfg, skipDocker); err != nil {
		log.Fatalf("Error generating/deploying docker-compose: %v", err)
	}

	log.Println("===== GitHub Runner Manager Completed =====")
}

// handleAddRepositories handles adding multiple repositories to the configuration
func handleAddRepositories(cfg *types.Config, repos string, numRunners int, configFile string) error {
	// Parse comma-separated repositories
	repoList := strings.Split(repos, ",")

	addedCount := 0
	for _, repoStr := range repoList {
		repoStr = strings.TrimSpace(repoStr)
		if repoStr == "" {
			continue
		}

		owner, name, err := utils.ParseRepo(repoStr)
		if err != nil {
			log.Printf("Warning: Failed to parse repository '%s': %v", repoStr, err)
			continue
		}

		if utils.RepoExists(*cfg, owner, name) {
			log.Printf("Repository %s/%s already exists in configuration", owner, name)
			continue
		}

		newRepo := utils.CreateNewRepository(owner, name, numRunners)
		cfg.Repositories = append(cfg.Repositories, newRepo)
		log.Printf("Added repository %s/%s with %d runners", owner, name, numRunners)
		addedCount++
	}

	if addedCount > 0 {
		// Save updated config
		if err := config.Save(configFile, *cfg); err != nil {
			log.Printf("Warning: Failed to save config: %v", err)
		}
		log.Printf("Successfully added %d new repositories to configuration", addedCount)
	} else {
		log.Println("No new repositories were added")
	}

	return nil
}

// updateRepositoriesAndFetchTokens updates repository info and fetches GitHub tokens
func updateRepositoriesAndFetchTokens(cfg *types.Config, githubToken string) error {
	for i := range cfg.Repositories {
		repo := &cfg.Repositories[i]
		utils.UpdateRepositoryInfo(repo)

		log.Printf("Fetching token for %s...", repo.FullName)
		token, err := github.GetRunnerToken(*repo, githubToken)
		if err != nil {
			return fmt.Errorf("getting token for %s: %w", repo.FullName, err)
		}
		repo.Token = token
		log.Printf("✓ Token obtained for %s (expires in ~1 hour)", repo.FullName)
	}
	return nil
}

// generateAndDeployDockerCompose generates docker-compose.yml and optionally deploys it
func generateAndDeployDockerCompose(cfg types.Config, skipDocker bool) error {
	// Backup existing docker-compose.yml
	if err := docker.BackupFile("docker-compose.yml"); err != nil {
		log.Printf("Warning: Failed to backup docker-compose.yml: %v", err)
	}

	// Generate new docker-compose.yml
	log.Println("Generating docker-compose.yml...")
	composeContent, err := docker.GenerateDockerCompose(cfg)
	if err != nil {
		return fmt.Errorf("generating docker-compose.yml: %w", err)
	}

	if err := os.WriteFile("docker-compose.yml", composeContent, 0644); err != nil {
		return fmt.Errorf("writing docker-compose.yml: %w", err)
	}

	// Validate the generated file
	log.Println("Validating docker-compose.yml...")
	if err := docker.ValidateDockerCompose(); err != nil {
		return err
	}
	log.Println("✅ docker-compose.yml validated successfully!")

	if !skipDocker {
		return deployDockerContainers(cfg)
	}

	return nil
}

// deployDockerContainers stops, updates, and starts Docker containers
func deployDockerContainers(cfg types.Config) error {
	// Stop existing containers
	log.Println("Stopping current runners...")
	if err := docker.RunDockerCompose("down"); err != nil {
		log.Printf("Warning: Error stopping containers: %v", err)
	}

	// Pull latest image
	log.Println("Pulling latest runner image...")
	if err := docker.RunDockerCompose("pull"); err != nil {
		log.Printf("Warning: Error pulling image: %v", err)
	}

	// Start new containers
	log.Println("Starting updated runners...")
	if err := docker.RunDockerCompose("up", "-d", "--remove-orphans"); err != nil {
		return fmt.Errorf("starting containers: %w", err)
	}

	// Verify containers are running
	log.Println("Waiting for containers to start...")
	time.Sleep(10 * time.Second)

	runningCount, totalCount, err := docker.GetRunningContainerCount()
	if err != nil {
		log.Printf("Warning: Error checking container status: %v", err)
		return nil
	}

	expectedTotal := utils.GetTotalRunnerCount(cfg)
	log.Printf("Running containers: %d/%d (expected: %d)", runningCount, totalCount, expectedTotal)

	if runningCount == expectedTotal {
		log.Println("✅ All runners updated and running successfully!")
	} else {
		log.Printf("⚠️  Warning: Only %d out of %d runners are running.", runningCount, expectedTotal)
		log.Println("Check logs with: docker-compose logs")
	}

	return nil
}
