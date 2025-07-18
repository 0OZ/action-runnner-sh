package docker

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github-runner-manager/internal/types"
)

const dockerComposeTemplate = `x-common-env: &common-env
  environment:
    RUNNER_SCOPE: "repo"
    LABELS: "linux,x64,gpu"
  security_opt:
    - label:disable
  volumes:
    - "/var/run/docker.sock:/var/run/docker.sock"
    - "/tmp/runner:/tmp/runner"
{{range .Repositories}}
x-{{.AnchorName}}: &{{.AnchorName}}
  RUNNER_TOKEN: "{{.Token}}"
{{end}}

services:
{{- range $repo := .Repositories}}
{{- range .Runners}}
  {{.ServiceName}}:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *{{$repo.AnchorName}}
      REPO_URL: {{$repo.URL}}
      RUNNER_NAME: {{.RunnerName}}
      RUNNER_WORKDIR: {{.WorkDir}}
    <<: *common-env
{{end}}
{{- end}}`

// GenerateDockerCompose generates the docker-compose.yml content
func GenerateDockerCompose(config types.Config) ([]byte, error) {
	tmpl, err := template.New("docker-compose").Parse(dockerComposeTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return buf.Bytes(), nil
}

// BackupFile creates a backup of the existing docker-compose.yml
func BackupFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // No file to backup
	}

	backupName := fmt.Sprintf("%s.backup.%d", filename, time.Now().Unix())
	input, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if err := os.WriteFile(backupName, input, 0644); err != nil {
		return fmt.Errorf("writing backup: %w", err)
	}

	log.Printf("Created backup: %s", backupName)
	return nil
}

// ValidateDockerCompose validates the docker-compose.yml file
func ValidateDockerCompose() error {
	cmd := exec.Command("docker-compose", "config")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("validation failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// RunDockerCompose executes docker-compose commands
func RunDockerCompose(args ...string) error {
	cmd := exec.Command("docker-compose", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetRunningContainerCount returns the count of running containers and total containers
func GetRunningContainerCount() (int, int, error) {
	cmd := exec.Command("docker-compose", "ps", "--services", "--filter", "status=running")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("getting running containers: %w", err)
	}

	runningCount := 0
	if len(strings.TrimSpace(string(output))) > 0 {
		runningCount = len(strings.Split(strings.TrimSpace(string(output)), "\n"))
	}

	cmdTotal := exec.Command("docker-compose", "ps", "--services")
	outputTotal, err := cmdTotal.Output()
	if err != nil {
		return runningCount, 0, fmt.Errorf("getting total containers: %w", err)
	}

	totalCount := 0
	if len(strings.TrimSpace(string(outputTotal))) > 0 {
		totalCount = len(strings.Split(strings.TrimSpace(string(outputTotal)), "\n"))
	}

	return runningCount, totalCount, nil
}
