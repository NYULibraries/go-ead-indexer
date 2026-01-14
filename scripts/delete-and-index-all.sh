#!/usr/bin/env bash

set -euo pipefail

# check required positional arg: path to the cloned EAD repo (must be a directory).
# `${1:-}` safely handles the case where no argument is passed under `set -u`,
ead_repo_path="${1:-}"
if [ -z "$ead_repo_path" ] || [ ! -d "$ead_repo_path" ]; then
  echo >&2 "Usage: $0 <full path to EAD repository>"
  exit 1
fi

# check that go indexer has its env var
if [ -z "${SOLR_ORIGIN_WITH_PORT:-}" ]; then
  echo >&2 "Must set SOLR_ORIGIN_WITH_PORT; aborting!"
  exit 1
fi

# Solr core/collection used for delete + verify.
SOLR_CORE="${SOLR_CORE:-findingaids}"

# HARD SAFETY GUARDRAIL:
# Refuse to run unless SOLR_ORIGIN_WITH_PORT matches one of the known Solr9 endpoints.
case "$SOLR_ORIGIN_WITH_PORT" in
  http://astra-dev.library.nyu.edu:8983|http://astra.library.nyu.edu:8983)
    ;;
  *)
    echo >&2 "Refusing to run: unexpected SOLR_ORIGIN_WITH_PORT=$SOLR_ORIGIN_WITH_PORT"
    echo >&2 "Expected one of:"
    echo >&2 "  - http://astra-dev.library.nyu.edu:8983"
    echo >&2 "  - http://astra.library.nyu.edu:8983"
    exit 2
    ;;
esac

echo "Target SOLR_ORIGIN_WITH_PORT=$SOLR_ORIGIN_WITH_PORT"
echo "Target SOLR_CORE=$SOLR_CORE"

echo "Preflight: Solr admin/info/system"
curl -fsS "$SOLR_ORIGIN_WITH_PORT/solr/admin/info/system?wt=json" >/dev/null

echo "Step 2: delete all docs (commit=true) via XML body (stream.body is disabled)"
curl -fsS -X POST \
  "$SOLR_ORIGIN_WITH_PORT/solr/$SOLR_CORE/update?commit=true&wt=json" \
  -H 'Content-Type: text/xml' \
  --data-binary '<delete><query>*:*</query></delete>' \
  >/dev/null

echo "Step 3: wait until index is empty (numFound=0)"
max_attempts="${MAX_EMPTY_CHECKS:-60}"
sleep_s="${EMPTY_CHECK_SLEEP_SECONDS:-5}"

attempt=0
while true; do
  attempt=$((attempt + 1))

  resp="$(curl -fsS \
    "$SOLR_ORIGIN_WITH_PORT/solr/$SOLR_CORE/select?q=*:*&rows=0&wt=json")"

  # Extract "numFound" as an integer from Solr JSON using sed:
  # match the "numFound" field, capture the digits, and print only the capture.
  numFound="$(echo "$resp" | sed -n 's/.*"numFound":[ ]*\([0-9][0-9]*\).*/\1/p')"

  echo "Attempt $attempt/$max_attempts: numFound=${numFound:-UNKNOWN}"

  if [ "${numFound:-}" = "0" ]; then
    echo "Confirmed empty."
    # Print a short snippet of the Solr JSON response (first 300 bytes).
    # This is enough to show `numFound: 0` for a Jira ticket without dumping the full response
    echo "$resp" | head -c 300
    echo
    break
  fi

  if [ "$attempt" -ge "$max_attempts" ]; then
    echo >&2 "ERROR: index not empty after $max_attempts attempts; aborting."
    echo >&2 "Last response snippet:"
    echo "$resp" | head -c 500 >&2
    echo >&2
    exit 1
  fi

  sleep "$sleep_s"
done

echo "Step 4: run full re-index"
exec "$(dirname "$0")/index-all.sh" "$ead_repo_path"
