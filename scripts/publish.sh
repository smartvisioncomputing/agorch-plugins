#!/usr/bin/env bash
# Signs the packages in index/ and publishes them to an Agorch registry through
# the ordinary publisher path (account, claim, key, artifact, listing, release,
# checks, review, stable channel), then applies the presentation fields from
# catalog.json. Idempotent: it reads the API's state and only does what is
# missing, so a second run changes nothing.
# Usage: scripts/publish.sh [-index index] [-catalog catalog.json]
# Environment (values are never printed; see README, "Publishing"):
#   AGORCH_API, AGORCH_PUBLISHER_TOKEN, AGORCH_PLATFORM_TOKEN, AGORCH_PUBLISHER_SEED_FILE
#   AGORCH_OWNER_EMAIL (only to create the account), AGORCH_PUBLISHER_KEY_ID (optional)
#   AICL (the aicl binary for the full check; default: aicl on PATH)
set -euo pipefail
cd "$(dirname "$0")/.."

for v in AGORCH_API AGORCH_PUBLISHER_TOKEN AGORCH_PLATFORM_TOKEN AGORCH_PUBLISHER_SEED_FILE; do
  [ -n "${!v:-}" ] || { echo "publish: $v is not set" >&2; exit 2; }
done

# The full check, never only the static one: every package in index/ is
# rebuilt from the commit it names and must match byte for byte, so a sealed
# file edited by hand (its license, scopes or tools) publishes nothing.
# Needs the aicl binary (AICL=/path/to/aicl).
scripts/check.sh
exec go run ./cmd/publish "$@"
