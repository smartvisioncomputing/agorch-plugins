#!/usr/bin/env bash
# Regenerates the sealed packages in two phases:
#   1. connectors -> index/: one sealed package per connector, built from the
#      last commit that touched connectors/<name> and pinned to it with
#      --source-commit, --source-tree and --source-path (the package's
#      presentation.source names exactly what it was built from);
#   2. plugins -> aicl.lock.json and starters/: each plugin's lock is resolved
#      against the new index/, and starters/ gets each plugin's sealed
#      agent_system package.
# The working tree must be clean: commit the connector or plugin change first,
# run this script, then commit index/, the locks and starters/ together.
# Usage: AICL=/path/to/aicl scripts/index.sh
set -euo pipefail

AICL="${AICL:-aicl}"
cd "$(dirname "$0")/.."
if [ -n "$(git status --porcelain)" ]; then
  git status --short >&2
  echo "index: the working tree is not clean; commit the change first" >&2
  exit 1
fi
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "== phase 1: connectors -> index/"
rm -f index/*.json
for dir in connectors/*/; do
  path="${dir%/}"
  name="$(basename "$path")"
  commit="$(git rev-list -1 HEAD -- "$path")"
  tree="$(git rev-parse "$commit:$path")"
  mkdir -p "$tmp/src-$name"
  git archive "$commit" "$path" | tar -x -C "$tmp/src-$name"
  "$AICL" plugin build "$tmp/src-$name/$path" --out "$tmp/$name" \
    --source-commit "$commit" --source-tree "$tree" --source-path "$path" >/dev/null
  cp "$tmp/$name"/*.json index/
  echo "$path ${commit:0:12}"
done

echo "== phase 2: plugins -> aicl.lock.json and starters/"
rm -rf starters && mkdir starters
for dir in plugins/*/; do
  name="$(basename "$dir")"
  rm -f "$dir/aicl.lock.json"
  "$AICL" plugin build "$dir" --index index --out "$tmp/p-$name" >"$tmp/p-$name.out"
  cp "$tmp/p-$name/aicl.lock.json" "$dir/aicl.lock.json"
  for f in $(jq -r '.packages[] | select(.kind == "agent_system") | .file' "$tmp/p-$name.out"); do
    cp "$tmp/p-$name/$f" starters/
  done
done
ls index starters
