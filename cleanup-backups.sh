#!/bin/bash

# Keep only the 5 most recent backup files
RUNNER_DIR="/home/christian/actions-runner"
cd "$RUNNER_DIR" || exit 1

# Find docker-compose backup files older than 7 days and delete them
find "$RUNNER_DIR" -name "docker-compose.yml.backup.*" -type f -mtime +7 -delete

# If we still have more than 5 backup files, keep only the 5 most recent
BACKUP_COUNT=$(find "$RUNNER_DIR" -name "docker-compose.yml.backup.*" | wc -l)
if [ "$BACKUP_COUNT" -gt 5 ]; then
  find "$RUNNER_DIR" -name "docker-compose.yml.backup.*" -type f | sort | head -n -5 | xargs rm -f
fi

# Log the cleanup
echo "$(date): Cleaned up old backup files. Remaining: $(find "$RUNNER_DIR" -name "docker-compose.yml.backup.*" | wc -l)" >> "$RUNNER_DIR/cleanup.log"
