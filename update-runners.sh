#!/bin/bash

# Script to update GitHub Action runner tokens while preserving YAML anchors
# Usage: ./fixed-script.sh <github_token>

# Check if GitHub Personal Access Token is provided
if [ -z "$1" ]; then
  echo "Error: GitHub Personal Access Token required"
  echo "Usage: ./fixed-script.sh <github_token>"
  exit 1
fi

GITHUB_TOKEN="$1"
COMPOSE_FILE="docker-compose.yml"
FRONTEND_REPO="0OZ/monorepo-frontend"
BACKEND_REPO="0OZ/auftrag-select-backend"
TOOLBOX_REPO="0OZ/Toolbox"

echo "===== Starting Token Update Process ====="

# Function to get new runner token for a repository
get_runner_token() {
  local repo=$1
  echo "Fetching new registration token for $repo..." >&2
  
  TOKEN_RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/$repo/actions/runners/registration-token")
  
  # Extract the token using jq if available
  if command -v jq >/dev/null 2>&1; then
    RUNNER_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r .token)
  else
    # Fallback to a more reliable extraction method using sed
    RUNNER_TOKEN=$(echo "$TOKEN_RESPONSE" | sed -n 's/.*"token": "\([^"]*\)".*/\1/p')
  fi
  
  if [ -z "$RUNNER_TOKEN" ]; then
    echo "Error: Failed to fetch runner token for $repo" >&2
    echo "Response: $TOKEN_RESPONSE" >&2
    exit 1
  fi
  
  echo "Successfully obtained new token for $repo" >&2
  # Output only the token
  echo "$RUNNER_TOKEN"
}

# First, make a backup of the current file with timestamp
BACKUP_FILE="${COMPOSE_FILE}.backup.$(date +%s)"
echo "Making backup of current docker-compose.yml to $BACKUP_FILE..."
cp "$COMPOSE_FILE" "$BACKUP_FILE"

# Get new runner tokens
echo "Getting frontend token..."
F_API_KEY=$(get_runner_token $FRONTEND_REPO)
echo "Getting backend token..."
B_API_KEY=$(get_runner_token $BACKEND_REPO)
echo "Getting toolbox token..."
TOOLBOX_API_KEY=$(get_runner_token $TOOLBOX_REPO)

# Debug - remove these in production
echo "Frontend token (sanitized): ${F_API_KEY:0:5}...${F_API_KEY: -5}"
echo "Backend token (sanitized): ${B_API_KEY:0:5}...${B_API_KEY: -5}"
echo "Toolbox token (sanitized): ${TOOLBOX_API_KEY:0:5}...${TOOLBOX_API_KEY: -5}"

# Create updated docker-compose.yml file while preserving YAML anchors
echo "Creating updated docker-compose.yml file..."
cat > "$COMPOSE_FILE" << EOL
x-common-env: &common-env
  environment:
    RUNNER_SCOPE: "repo"
    LABELS: "linux,x64,gpu"
  security_opt:
    - label:disable
  volumes:
    - "/var/run/docker.sock:/var/run/docker.sock"
    - "/tmp/runner:/tmp/runner"

x-b-api-env: &b-api-env
  RUNNER_TOKEN: "${B_API_KEY}"

x-f-api-env: &f-api-env
  RUNNER_TOKEN: "${F_API_KEY}"

x-toolbox-api-env: &toolbox-api-env
  RUNNER_TOKEN: "${TOOLBOX_API_KEY}"


services:
  worker:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *f-api-env
      REPO_URL: https://github.com/0OZ/monorepo-frontend
      RUNNER_NAME: tensor-1
      RUNNER_WORKDIR: /tmp/runner/f1
    <<: *common-env

  frontend-1:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *f-api-env
      REPO_URL: https://github.com/0OZ/monorepo-frontend
      RUNNER_NAME: frontend-1
      RUNNER_WORKDIR: /tmp/runner/f2
    <<: *common-env

  frontend-2:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *f-api-env
      REPO_URL: https://github.com/0OZ/monorepo-frontend
      RUNNER_NAME: frontend-2
      RUNNER_WORKDIR: /tmp/runner/f3
    <<: *common-env

  worker-2:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *f-api-env
      REPO_URL: https://github.com/0OZ/monorepo-frontend
      RUNNER_NAME: tensor-2
      RUNNER_WORKDIR: /tmp/runner/f4
    <<: *common-env

  worker-3:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *b-api-env
      REPO_URL: https://github.com/0OZ/auftrag-select-backend
      RUNNER_NAME: tensor-7
      RUNNER_WORKDIR: /tmp/runner/work-3
    <<: *common-env

  worker-4:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *b-api-env
      REPO_URL: https://github.com/0OZ/auftrag-select-backend
      RUNNER_NAME: tensor-3
      RUNNER_WORKDIR: /tmp/runner/work-4
    <<: *common-env

  worker-5:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *b-api-env
      REPO_URL: https://github.com/0OZ/auftrag-select-backend
      RUNNER_NAME: tensor-4
      RUNNER_WORKDIR: /tmp/runner/work-5
    <<: *common-env

  worker-6:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *b-api-env
      REPO_URL: https://github.com/0OZ/auftrag-select-backend
      RUNNER_NAME: tensor-6
      RUNNER_WORKDIR: /tmp/runner/work-7
    <<: *common-env

  worker-9:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *b-api-env
      REPO_URL: https://github.com/0OZ/auftrag-select-backend
      RUNNER_NAME: tensor-9
      RUNNER_WORKDIR: /tmp/runner/work-9
    <<: *common-env

  worker-10:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *b-api-env
      REPO_URL: https://github.com/0OZ/auftrag-select-backend
      RUNNER_NAME: tensor-10
      RUNNER_WORKDIR: /tmp/runner/work-10
    <<: *common-env

  worker-11:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *toolbox-api-env
      REPO_URL: https://github.com/0OZ/Toolbox
      RUNNER_NAME: toolbox-1
      RUNNER_WORKDIR: /tmp/runner/work-10
    <<: *common-env

  worker-12:
    image: myoung34/github-runner:latest
    restart: "always"
    environment:
      <<: *b-api-env
      REPO_URL: https://github.com/0OZ/auftrag-select-backend
      RUNNER_NAME: tensor-5
      RUNNER_WORKDIR: /tmp/runner/work-12
    <<: *common-env
EOL

# Validate the new docker-compose file
echo "Validating the updated docker-compose.yml file..."
if docker-compose config > /dev/null; then
  echo "✅ docker-compose.yml file validated successfully!"
else
  echo "❌ Error validating docker-compose.yml. Restoring from backup..."
  cp "$BACKUP_FILE" "$COMPOSE_FILE"
  echo "Backup restored. Please check docker-compose.yml manually."
  exit 1
fi

# Update containers
echo "Stopping current runners..."
docker-compose down

echo "Pulling latest runner image..."
docker-compose pull

echo "Starting updated runners with new tokens..."
docker-compose up -d

# Check if containers are running
echo "Verifying runners are running..."
sleep 10
RUNNING_CONTAINERS=$(docker-compose ps --services --filter "status=running" | wc -l)
TOTAL_CONTAINERS=$(docker-compose ps --services | wc -l)

echo "Running containers: $RUNNING_CONTAINERS/$TOTAL_CONTAINERS"

if [ "$RUNNING_CONTAINERS" -eq "$TOTAL_CONTAINERS" ]; then
  echo "✅ All runners updated and running successfully!"
else
  echo "⚠️ Warning: Only $RUNNING_CONTAINERS out of $TOTAL_CONTAINERS runners are running."
  echo "Check logs with: docker-compose logs"
fi

echo "===== Token Update Process Completed ====="
