#!/usr/bin/env bash
# Remove stale GitHub Actions runner temp dirs from /tmp/runner.
# Must run as root: the per-job dirs are root-owned and /tmp/runner is not
# world-writable. Anything modified within MAX_AGE_MIN is kept, so in-flight
# and very recent jobs are never touched.
set -euo pipefail

RUNNER_TMP="/tmp/runner"
MAX_AGE_MIN="${MAX_AGE_MIN:-1440}"   # 24h; override via env
LOG="/home/christian/actions-runner/cleanup-tmp-runner.log"

[ -d "$RUNNER_TMP" ] || exit 0

before=$(du -sh "$RUNNER_TMP" 2>/dev/null | cut -f1)
count=$(find "$RUNNER_TMP" -maxdepth 1 -mindepth 1 -type d -mmin +"$MAX_AGE_MIN" 2>/dev/null | wc -l)

# -depth/+rm -rf via -exec {} + ; ignore races where a dir vanishes mid-run
find "$RUNNER_TMP" -maxdepth 1 -mindepth 1 -type d -mmin +"$MAX_AGE_MIN" \
  -exec rm -rf {} + 2>/dev/null || true

after=$(du -sh "$RUNNER_TMP" 2>/dev/null | cut -f1)
printf '%s  removed %s stale dirs  (%s -> %s)\n' \
  "$(date -Is)" "$count" "$before" "$after" >> "$LOG"
