# GitHub Actions Runner Setup

This repository contains a Go-based runner manager and scripts for setting up and managing GitHub Actions self-hosted runners using Docker.

## Overview

The system uses a Docker Compose configuration to run multiple GitHub Actions runners in containers, connecting to different repositories. The setup supports:

- Multiple runners for different repositories
- Automatic token refresh
- YAML anchors for clean configuration
- Configuration management through JSON files
- Dynamic repository addition

## Runner Manager (Go Application)

The main component is now a Go application (`runner-manager`) that provides comprehensive management of GitHub Actions runners.

### Installation

#### Using Make

```bash
# Build the application
make build

# Install dependencies
make deps

# Install to system path (requires sudo)
make install
```

#### Manual Build

```bash
go build -o build/runner-manager ./cmd/runner-manager
```

### Usage

#### Basic Token Update

```bash
./build/runner-manager -token <GITHUB_TOKEN>
```

#### Add New Repository

```bash
./build/runner-manager -token <GITHUB_TOKEN> -add-repo owner/repo -runners 3
```

#### Using Environment Variable

```bash
export GITHUB_TOKEN="your_token_here"
./build/runner-manager
```

#### Configuration File

```bash
./build/runner-manager -token <GITHUB_TOKEN> -config my-config.json
```

#### Skip Docker Operations (Generate config only)

```bash
./build/runner-manager -token <GITHUB_TOKEN> -skip-docker
```

### Features

- **Automatic token management**: Fetches new registration tokens for all repositories
- **Dynamic repository addition**: Add new repositories without manual configuration
- **Configuration persistence**: Stores repository configuration in JSON format
- **Docker Compose generation**: Automatically generates docker-compose.yml with YAML anchors
- **Backup and validation**: Creates backups and validates configurations before deployment
- **Container management**: Stops, updates, and starts Docker containers
- **Status verification**: Verifies that all containers are running correctly after update

### Configuration

The runner manager uses a JSON configuration file (`runners-config.json` by default) to store repository and runner information:

```json
{
  "repositories": [
    {
      "owner": "0OZ",
      "name": "monorepo-frontend",
      "fullName": "0OZ/monorepo-frontend",
      "url": "https://github.com/0OZ/monorepo-frontend",
      "token": "...",
      "anchorName": "f-api-env",
      "runners": [
        {
          "serviceName": "frontend-runner-1",
          "runnerName": "frontend-runner-1",
          "workDir": "/tmp/runner/frontend-1"
        }
      ]
    }
  ]
}
```

## Legacy Scripts

### update-runners.sh

Legacy bash script for updating GitHub Actions runner registration tokens.

#### Legacy Script Usage

```bash
./update-runners.sh <GITHUB_TOKEN>
```

Where `<GITHUB_TOKEN>` is a GitHub Personal Access Token with appropriate permissions to register runners.

## Project Structure

```text
├── cmd/
│   └── runner-manager/        # Main application entry point
│       └── main.go
├── internal/
│   ├── config/               # Configuration management
│   │   ├── config.go
│   │   └── default.go
│   ├── docker/               # Docker Compose generation and management
│   │   └── compose.go
│   ├── github/               # GitHub API client
│   │   └── client.go
│   ├── types/                # Type definitions
│   │   └── types.go
│   └── utils/                # Utility functions
│       └── repository.go
├── build/                    # Build output directory
├── docker-compose.yml        # Generated Docker Compose configuration
├── runners-config.json       # Runner configuration (auto-generated)
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── Makefile                  # Build automation
```

## Development

### Available Make Targets

```bash
make help          # Show available targets
make build         # Build the application
make deps          # Install dependencies
make test          # Run tests
make clean         # Clean build artifacts
make run           # Build and run the application
make fmt           # Format Go code
make lint          # Lint Go code (requires golangci-lint)
make tidy          # Tidy go modules
make install       # Install binary to /usr/local/bin
```

### Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...
```

## Docker Compose Configuration

The `docker-compose.yml` file is automatically generated and uses YAML anchors to avoid repetition:

- `x-common-env`: Contains common environment variables and volumes
- `x-{repo}-env`: Contains repository-specific runner tokens (dynamically generated)

## Repositories

The setup can manage runners for multiple repositories. Default repositories include:

- `0OZ/monorepo-frontend`: Frontend repository
- `0OZ/auftrag-select-backend`: Backend repository  
- `0OZ/Toolbox`: Toolbox repository

Additional repositories can be added using the `-add-repo` flag.

## Requirements

### For Go Application

- Go 1.23.1 or later
- Docker and Docker Compose
- GitHub Personal Access Token with appropriate permissions

### For Legacy Scripts

- Docker and Docker Compose
- curl
- jq (optional, for better token parsing)
- GitHub Personal Access Token with appropriate permissions

## Troubleshooting

### Go Application Issues

1. **Build errors**: Run `make deps` to ensure dependencies are installed
2. **Permission errors**: Check that your GitHub token has `repo` scope and admin access to target repositories
3. **Docker issues**: Ensure Docker service is running and docker-compose is available
4. **Configuration issues**: Delete `runners-config.json` to reset to defaults

### General Issues

1. Check the logs with `docker-compose logs`
2. Ensure your GitHub token has the correct permissions
3. Verify network connectivity to GitHub API
4. If needed, restore from the backup file created during the update process

### Common Commands

```bash
# Check container status
docker-compose ps

# View logs for specific service
docker-compose logs frontend-runner-1

# Restart all services
docker-compose restart

# Rebuild and restart everything
make build && ./build/runner-manager -token $GITHUB_TOKEN
```

## Migration from Bash Scripts

If you're migrating from the previous bash script setup:

1. **Backup your current setup**:

   ```bash
   cp docker-compose.yml docker-compose.yml.backup
   ```

2. **Build the Go application**:

   ```bash
   make build
   ```

3. **Run the migration**:

   ```bash
   # The application will detect existing docker-compose.yml and create a configuration
   ./build/runner-manager -token $GITHUB_TOKEN
   ```

4. **Verify the setup**:

   ```bash
   docker-compose ps
   ```

The Go application will automatically:

- Parse your existing docker-compose.yml structure
- Generate a `runners-config.json` configuration file
- Maintain compatibility with your current runner setup

## Security Notes

- Never commit your GitHub token to version control
- Regularly rotate your GitHub tokens
- Consider using environment variables or a secrets manager instead of passing tokens directly in command line arguments
