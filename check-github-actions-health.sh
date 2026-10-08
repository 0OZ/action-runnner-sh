#!/bin/bash

# Check if GitHub Actions job service is healthy and log state transitions.
# Exit codes: 0 = service healthy, 1 = service degraded (503), 2 = probe error
# Detection: a POST to the acquirejob endpoint returns 503 only while GitHub's
# upstream is degraded. Any other status (401, 200, ...) means the service is up.

ENDPOINT="https://run-actions-2-azure-eastus.actions.githubusercontent.com/109/acquirejob"
STATE_FILE="/tmp/gh-actions-health.state"
LOG_FILE="/home/christian/actions-runner/gh-actions-health.log"

HTTP_CODE=$(curl -sS -o /dev/null -w "%{http_code}" -X POST --max-time 15 "$ENDPOINT" 2>/dev/null)

if [ -z "$HTTP_CODE" ]; then
  echo "$(date -u): probe error (curl failed)" >> "$LOG_FILE"
  exit 2
fi

if [ "$HTTP_CODE" = "503" ]; then
  echo "$(date -u): GitHub Actions service DOWN (HTTP 503) - runners waiting" >> "$LOG_FILE"
  echo "down" > "$STATE_FILE"
  exit 1
fi

PREV_STATE=$(cat "$STATE_FILE" 2>/dev/null)
if [ "$PREV_STATE" = "down" ] || [ -z "$PREV_STATE" ]; then
  echo "$(date -u): GitHub Actions service BACK UP (HTTP $HTTP_CODE) - runners should auto-reconnect" >> "$LOG_FILE"
fi
echo "up" > "$STATE_FILE"
exit 0
