#!/usr/bin/env bash
set -euo pipefail

# Delete GitHub Actions runs older than a given age
# Usage: ./delete-old-runs.sh OWNER/REPO MINUTES [--dry-run]
#
# Examples:
#   ./delete-old-runs.sh myorg/myrepo 120
#   ./delete-old-runs.sh myorg/myrepo 60 --dry-run

usage() {
  echo "Usage: $0 OWNER/REPO MINUTES [--dry-run]"
  exit 1
}

REPO="${1:-}"; shift || true
AGE_MINUTES="${1:-}"; shift || true
DRY_RUN="${1:-}"

[[ -z "${REPO}" || -z "${AGE_MINUTES}" ]] && usage
if ! [[ "${AGE_MINUTES}" =~ ^[0-9]+$ ]]; then
  echo "MINUTES must be an integer" >&2
  exit 1
fi

echo "Fetching runs for ${REPO} and deleting those older than ${AGE_MINUTES} minutes..."
[[ "${DRY_RUN:-}" == "--dry-run" ]] && echo "Dry run mode: will only print IDs and not delete."

# We fetch a large page (up to 1000). Increase if needed or loop over pages via gh api.
# Fields: databaseId (run ID), createdAt (RFC3339), status
JQ_FILTER="
  .[]
  | select(.status == \"completed\")
  | select( (now - (.createdAt | fromdateiso8601)) > (${AGE_MINUTES} * 60) )
  | .databaseId
"

# If your gh doesn't support `fromdateiso8601`/`now`, see fallback note below.
RUN_IDS="$(gh run list \
  --repo "${REPO}" \
  --limit 1000 \
  --json databaseId,createdAt,status \
  --jq "${JQ_FILTER}" || true)"

if [[ -z "${RUN_IDS}" ]]; then
  echo "No runs to delete (either none found, or none older than ${AGE_MINUTES} minutes)."
  exit 0
fi

COUNT=0
for ID in ${RUN_IDS}; do
  if [[ "${DRY_RUN:-}" == "--dry-run" ]]; then
    echo "[dry-run] Would delete run ID ${ID}"
  else
    echo "Deleting run ${ID}..."
    gh api --silent --method DELETE "repos/${REPO}/actions/runs/${ID}"
  fi
  COUNT=$((COUNT + 1))
done

if [[ "${DRY_RUN:-}" == "--dry-run" ]]; then
  echo "Dry run complete. ${COUNT} runs would have been deleted."
else
  echo "✅ Deleted ${COUNT} runs older than ${AGE_MINUTES} minutes from ${REPO}."
fi
