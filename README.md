# GitHub Actions Runner Setup

This repository contains scripts for setting up and managing GitHub Actions self-hosted runners using Docker.

## Overview

The system uses a Docker Compose configuration to run multiple GitHub Actions runners in containers, connecting to different repositories. The setup supports:

- Multiple runners for different repositories
- Automatic token refresh
- YAML anchors for clean configuration

## Scripts

### update-runners.sh

This script updates the GitHub Actions runner registration tokens and refreshes all running containers.

#### Usage

```bash
./update-runners.sh <GITHUB_TOKEN>
```

Where `<GITHUB_TOKEN>` is a GitHub Personal Access Token with appropriate permissions to register runners.

#### Features

- Automatically fetches new registration tokens for all repositories
- Preserves YAML anchor structure in docker-compose.yml
- Creates a backup of the existing configuration before making changes
- Validates the new configuration before applying changes
- Restarts all containers with updated tokens
- Verifies that all containers are running correctly after update

## Docker Compose Configuration

The `docker-compose.yml` file uses YAML anchors to avoid repetition:

- `x-common-env`: Contains common environment variables and volumes
- `x-b-api-env`: Contains the backend repository runner token
- `x-f-api-env`: Contains the frontend repository runner token
- `x-toolbox-api-env`: Contains the toolbox repository runner token

## Repositories

The setup manages runners for the following repositories:

- `0OZ/monorepo-frontend`: Frontend repository
- `0OZ/auftrag-select-backend`: Backend repository
- `0OZ/Toolbox`: Toolbox repository

## Requirements

- Docker and Docker Compose
- curl
- jq (optional, for better token parsing)
- GitHub Personal Access Token with appropriate permissions

## Troubleshooting

If you encounter issues with the runner update process:

1. Check the logs with `docker-compose logs`
2. Ensure your GitHub token has the correct permissions
3. Verify network connectivity to GitHub API
4. If needed, restore from the backup file created during the update process

## Security Notes

- Never commit your GitHub token to version control
- Regularly rotate your GitHub tokens
- Consider using environment variables or a secrets manager instead of passing tokens directly in command line arguments
