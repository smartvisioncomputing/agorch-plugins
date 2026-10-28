#!/usr/bin/env bash
# Checks the whole repository, the same way locally and in CI.
# Static part (needs only git and jq):
#   1. catalog.json x index/ x plugins/;
#   2. secret scan of the tracked files (aicl also scans manifests and packages).
# With the aicl binary:
#   3. fmt -check on every .aicl;
#   4. each sealed connector in index/ is rebuilt from the source it names
#      (presentation.source: commit, tree, path), twice, and must match the
#      file byte for byte; that source must be the current content of the
#      connector at HEAD, a version is never sealed with two digests, and
#      every connector has exactly one package;
#   5. plugin validate on every connector and plugin (plugins against index/);
#   6. each plugin is built twice (same digests and bytes) and starters/ holds
#      exactly its agent_system packages.
# Usage: AICL=/path/to/aicl scripts/check.sh
#        scripts/check.sh static      # items 1 and 2 only, without aicl
set -uo pipefail

AICL="${AICL:-aicl}"
cd "$(dirname "$0")/.." || exit 2
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
status=0
fail() { echo "FAIL: $*"; status=1; }

echo "== catalog.json x index/ x plugins/"
# Closed category enum: the same list the API validates.
enum='["sales_marketing","payments_commerce","work_docs","communication","engineering"]'
# auth_kind from the auth_schemes of the sealed connection_definition.
in_index="$(jq -r '.package | "\(.namespace)/\(.package_id)"' index/*.json | sort)"
in_catalog="$(jq -r '.connectors[].id' catalog.json | sort)"
[ "$in_index" = "$in_catalog" ] || { diff <(echo "$in_index") <(echo "$in_catalog"); fail "catalog.json ids != index/ (< index, > catalog)"; }
[ "$(echo "$in_catalog" | uniq -d)" = "" ] || fail "repeated id in catalog.json"
jq -e --argjson enum "$enum" '
  ([.categories[].id] | sort) == ($enum | sort)
  and ([.connectors[].rank] | (sort == (unique)) and all(type == "number" and . == floor))
  and all(.connectors[];
        (.category as $c | $enum | index($c)) != null
        and (.vendor_official | type) == "boolean" and (.beta | type) == "boolean" and (.verified | type) == "boolean"
        and (.auth_kind | IN("oauth", "api_key", "none"))
        and (.display_name | length > 0) and (.summary | length > 0)
        and (.logo | test("^[a-z0-9-]+\\.svg$")))' catalog.json >/dev/null ||
  fail "catalog.json: category outside the enum, repeated rank, wrong type or empty field"
for f in index/*.json; do
  id="$(jq -r '.package | "\(.namespace)/\(.package_id)"' "$f")"
  got="$(jq -r '[.components[] | select(.ref.kind == "connection_definition") | .spec.auth_schemes[].kind] | first
      | if . == "oauth2_authorization_code_pkce" then "oauth" elif . == "none" then "none" else "api_key" end' "$f")"
  want="$(jq -r --arg id "$id" '.connectors[] | select(.id == $id) | .auth_kind' catalog.json)"
  [ "$got" = "$want" ] || fail "$id: auth_kind $want in catalog.json, $got in index/"
done
# Starters: one per plugins/ folder, naming only existing connectors, the same ones plugin.aicl requires.
[ "$(ls plugins)" = "$(jq -r '.starters[].id' catalog.json | sort)" ] || fail "catalog.json starters != plugins/ folders"
jq -e 'all(.starters[]; .status | IN("available", "coming_soon"))' catalog.json >/dev/null || fail "catalog.json: starter status outside available/coming_soon"
for dir in plugins/*/; do
  name="$(basename "$dir")"
  have="$(grep -o 'package *= *"[^"]*"' "$dir/plugin.aicl" | sed 's/.*"\(.*\)"/\1/' | sort)"
  want="$(jq -r --arg n "$name" '.starters[] | select(.id == $n) | .requires[]' catalog.json | sort)"
  [ "$have" = "$want" ] || fail "starter $name: catalog.json requires != plugin.aicl"
  for ref in $want; do echo "$in_index" | grep -qx "$ref" || fail "starter $name names a missing connector: $ref"; done
done
# Namespace and publisher: the registry admits the first_party claim only for publisher:ao / ao.
[ "$(jq -r '"\(.publisher_id) \(.namespace)"' catalog.json)" = "publisher:ao ao" ] || fail "catalog.json: publisher_id/namespace != publisher:ao / ao"
[ -z "$(grep -rEl 'publisher:agorch|namespace *= *"agorch"|"agorch/' connectors plugins index catalog.json)" ] || fail "agorch left as a namespace or publisher"
# No clash with the ids ao-packages already publishes under ao/ (AO_PACKAGES=<dir>/catalog/ao; skipped when unset).
if [ -d "${AO_PACKAGES:-}" ]; then
  taken="$(jq -r '"\(.package.namespace // .namespace)/\(.package.package_id // .package_id)"' "$AO_PACKAGES"/*/package.json | sort)"
  mine="$(jq -r '"\(.package.namespace)/\(.package.package_id)"' index/*.json; for f in plugins/*/plugin.aicl; do echo "ao/$(basename "$(dirname "$f")")"; done)"
  [ -z "$(comm -12 <(echo "$taken") <(echo "$mine" | sort))" ] || fail "id clashes with ao-packages: $(comm -12 <(echo "$taken") <(echo "$mine" | sort) | tr '\n' ' ')"
fi

echo "== secret scan"
# Credential shapes of the catalog's vendors plus generic ones. No file in this
# repository holds a secret value: only *_ref names.
pattern='(sk|rk|pk)_(live|test)_[0-9A-Za-z]{8,}|xox[abposr]-[0-9A-Za-z-]{8,}|xapp-[0-9A-Za-z-]{8,}|gh[opsur]_[0-9A-Za-z]{20,}|github_pat_[0-9A-Za-z_]{20,}|sntrys_[0-9A-Za-z_]{20,}|sntryu_[0-9a-f]{20,}|pat-(na|eu)[0-9]-[0-9a-f-]{20,}|AKIA[0-9A-Z]{16}|-----BEGIN [A-Z ]*PRIVATE KEY-----|(client_secret|api_key|access_token|password)[[:space:]]*[:=][[:space:]]*"[^"$]{8,}"'
# file:line only, never the value: the CI log must not repeat a secret.
hits="$(git ls-files -z --cached --others --exclude-standard | xargs -0 grep -nIE "$pattern" -- 2>/dev/null | cut -d: -f1-2)"
[ -n "$hits" ] && { echo "$hits"; fail "possible literal secret on the lines above"; }

if [ "${1:-}" = static ]; then
  [ "$status" = 0 ] && echo "static OK" || echo "static check failed"
  exit "$status"
fi

echo "== fmt -check"
"$AICL" fmt -check -f . || fail "fmt: run '$AICL fmt -f .'"

echo "== index/ rebuilt from the source each package names"
covered=""
for f in index/*.json; do
  read -r commit tree path < <(jq -r '.package.presentation.source // {} | "\(.commit) \(.tree) \(.path)"' "$f")
  case "$path" in connectors/*/* | connectors/.* | connectors/) bad=1 ;; connectors/*) bad= ;; *) bad=1 ;; esac
  [ -z "$bad" ] || { fail "$f: presentation.source.path is not connectors/<name>"; continue; }
  covered="$covered$path"$'\n'
  git merge-base --is-ancestor "$commit" HEAD 2>/dev/null || { fail "$f: source commit $commit is not in HEAD's history"; continue; }
  [ "$(git rev-parse "$commit:$path" 2>/dev/null)" = "$tree" ] || { fail "$f: source tree is not $path at $commit"; continue; }
  [ "$(git rev-parse "HEAD:$path")" = "$tree" ] && [ -z "$(git status --porcelain -- "$path")" ] ||
    { fail "$f: $path changed since it was sealed: commit it, then run scripts/index.sh"; continue; }
  # A version is sealed once: every committed copy of this file has its digest.
  digest="$(jq -r .canonical_digest "$f")"
  for c in $(git rev-list HEAD -- "$f"); do
    old="$(git show "$c:$f" 2>/dev/null | jq -r .canonical_digest)"
    [ -z "$old" ] || [ "$old" = "$digest" ] || { fail "$f: this version was sealed before as $old at ${c:0:12}: raise the version"; break; }
  done
  name="$(basename "$path")"
  mkdir -p "$tmp/src-$name"
  git archive "$commit" "$path" | tar -x -C "$tmp/src-$name"
  for run in 1 2; do
    "$AICL" plugin build "$tmp/src-$name/$path" --out "$tmp/c-$name-$run" \
      --source-commit "$commit" --source-tree "$tree" --source-path "$path" >/dev/null || fail "build $path ($run)"
  done
  diff -r "$tmp/c-$name-1" "$tmp/c-$name-2" >/dev/null || fail "bytes differ between two builds of $path"
  cmp -s "$tmp/c-$name-1/$(basename "$f")" "$f" && echo "$f $(jq -r .canonical_digest "$f")" ||
    fail "$f is not what $path at ${commit:0:12} builds: run scripts/index.sh"
done
[ "$(printf '%s' "$covered" | sort)" = "$(ls -d connectors/*/ | sed 's:/$::' | sort)" ] ||
  fail "index/ must hold exactly one package per connector: run scripts/index.sh"

mkdir -p "$tmp/starters"
for dir in connectors/*/ plugins/*/; do
  name="$(basename "$dir")"
  kind="$(basename "$(dirname "$dir")")"
  echo "== validate $kind/$name"
  if [ "$kind" = connectors ]; then
    "$AICL" plugin validate "$dir" || fail "validate $kind/$name"
    continue
  fi
  "$AICL" plugin validate "$dir" --index index || { fail "validate $kind/$name"; continue; }

  echo "== build 2x $kind/$name"
  for run in 1 2; do
    "$AICL" plugin build "$dir" --index index --out "$tmp/$kind-$name-$run" >"$tmp/$kind-$name-$run.out" ||
      fail "build $kind/$name ($run)"
    grep -o '"digest": "sha256:[0-9a-f]*"' "$tmp/$kind-$name-$run.out" >"$tmp/$kind-$name-$run.digests"
  done
  diff "$tmp/$kind-$name-1.digests" "$tmp/$kind-$name-2.digests" || fail "digests differ in $kind/$name"
  diff -r "$tmp/$kind-$name-1" "$tmp/$kind-$name-2" >/dev/null || fail "bytes differ in $kind/$name"
  cat "$tmp/$kind-$name-1.digests"
  f="$(jq -r '.packages[] | select(.kind == "agent_system") | .file' "$tmp/$kind-$name-1.out")"
  [ -n "$f" ] && cp "$tmp/$kind-$name-1/$f" "$tmp/starters/" || fail "build $kind/$name has no agent_system package"
done

echo "== starters/ equals the plugins' build"
diff -rq "$tmp/starters" starters || fail "starters/ is stale: run scripts/index.sh"

[ "$status" = 0 ] && echo "OK" || echo "check failed"
exit "$status"
