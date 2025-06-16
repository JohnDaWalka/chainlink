#!/usr/bin/env bash

max_retries=${1:-5} # Default to 5 if not provided
count=0
backoff=1  # Start with 1 second

until go run main.go download all \
  --output-dir ../ \
  --gh-token-env-var-name GITHUB_API_TOKEN \
  --cre-cli-version v0.2.0 \
  --capability-names cron \
  --capability-version v1.0.2-alpha
do
  ((count++))
  if (( count >= max_retries )); then
    echo "❌ Failed after $max_retries attempts." >&2
    exit 1
  fi
  echo "🔁 Retrying ($count/$max_retries) in ${backoff}s..." >&2
  sleep "$backoff"
  backoff=$((backoff * 2))
  backoff=$((backoff > 8 ? 8 : backoff))  # Cap at 8 seconds
done

echo "✅ Download succeeded." >&2