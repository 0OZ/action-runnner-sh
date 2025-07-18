#!/bin/bash

# Script to set up crontab for daily runner token refresh and cleanup
# Uses the Go runner-manager application instead of shell scripts
# Usage: ./crontab-setup.sh <github_token>

if [ -z "$1" ]; then
  echo "Error: GitHub Personal Access Token required"
  echo "Usage: ./crontab-setup.sh <github_token>"
  exit 1
fi

GITHUB_TOKEN="$1"
RUNNER_DIR="/home/christian/actions-runner"
TOKEN_SECRET_FILE="${RUNNER_DIR}/.github_token"
RUNNER_BINARY="${RUNNER_DIR}/runner-manager"

# Check if runner-manager binary exists
if [ ! -f "$RUNNER_BINARY" ]; then
  echo "Error: runner-manager binary not found at $RUNNER_BINARY"
  echo "Please run ./deploy-runner-manager.sh first to build and deploy the binary"
  exit 1
fi

# Store the token in a secure file
echo "$GITHUB_TOKEN" > "$TOKEN_SECRET_FILE"
chmod 600 "$TOKEN_SECRET_FILE"

# Create cleanup script
CLEANUP_SCRIPT="${RUNNER_DIR}/cleanup-backups.sh"
cat > "$CLEANUP_SCRIPT" << 'EOF'
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
EOF

chmod +x "$CLEANUP_SCRIPT"

# Create crontab entries
(crontab -l 2>/dev/null || echo "") | grep -v "$RUNNER_BINARY\|$CLEANUP_SCRIPT" > temp_crontab

# Add new entries
cat >> temp_crontab << EOF
# Update GitHub Runner tokens daily at 2:00 AM
0 2 * * * cd $RUNNER_DIR && GITHUB_TOKEN=\$(cat $TOKEN_SECRET_FILE) $RUNNER_BINARY >> $RUNNER_DIR/runner-update.log 2>&1

# Clean up old docker-compose backup files daily at 3:00 AM
0 3 * * * $CLEANUP_SCRIPT >> $RUNNER_DIR/cleanup.log 2>&1
EOF

# Install new crontab
crontab temp_crontab
rm temp_crontab

echo "✅ Crontab setup complete!"
echo "• GitHub token securely stored in $TOKEN_SECRET_FILE"
echo "• Runners will be updated daily at 2:00 AM using runner-manager"
echo "• Old backup files will be cleaned up daily at 3:00 AM"
echo "• Logs will be stored in $RUNNER_DIR/runner-update.log and $RUNNER_DIR/cleanup.log"
echo ""
echo "You can view the current crontab with: crontab -l"
echo "Make sure the runner-manager binary is present in $RUNNER_DIR"
