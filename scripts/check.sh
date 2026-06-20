#!/usr/bin/env bash
# Format and vet the biblio-cli Go module. Run before committing.
set -euo pipefail
cd "$(dirname "$0")/.."

echo "==> gofmt"
unformatted=$(gofmt -l .)
if [[ -n "$unformatted" ]]; then
  echo "The following files are not gofmt-formatted:"
  echo "$unformatted"
  exit 1
fi

echo "==> go vet"
go vet ./...
