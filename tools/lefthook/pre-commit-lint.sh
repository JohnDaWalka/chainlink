#!/usr/bin/env bash
# This script is executed by Lefthook on pre-commit.
# It runs golangci-lint on all Go packages that contain staged .go files,
# but only fails the commit if issues are found in the staged files themselves.
# This faithfully replicates the logic from the legacy .githooks/pre-commit script.

set -euo pipefail

# Allow devs to bypass the hook for quick commits.
if [[ "${CL_DEV_PRECOMMIT_LINT:-true}" == "false" ]]; then
    echo "ℹ️ CL_DEV_PRECOMMIT_LINT is false, skipping pre-commit linting."
    exit 0
fi

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "❌ golangci-lint not found in your PATH."
    echo "Please install it. See project documentation for the recommended version."
    echo "You can try running this from the repo root: asdf install"
    exit 1
fi

echo "🔍 Running pre-commit linting for staged Go files..."

# Get the list of staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [[ -z "$STAGED_GO_FILES" ]]; then
    echo "✅ No Go files staged, skipping golangci-lint."
    exit 0
fi

echo "Staged Go files found:"
while IFS= read -r line; do echo "  $line"; done <<< "$STAGED_GO_FILES"
echo

# Determine unique package directories from the staged files
PACKAGES_TO_LINT=$(echo "$STAGED_GO_FILES" | xargs -n1 dirname | sort -u)

if [[ -z "$PACKAGES_TO_LINT" ]]; then
    echo "✅ No packages to lint."
    exit 0
fi

echo "Determined packages to lint:"
while IFS= read -r line; do echo "  ./$line"; done <<< "$PACKAGES_TO_LINT"
echo

# Loop through each package, run linter, and filter results
HAS_ISSUES_IN_STAGED_FILES=false
REPO_ROOT=$(git rev-parse --show-toplevel)
ISSUES_OUTPUT=$(mktemp)
trap 'rm -f "$ISSUES_OUTPUT"' EXIT

for PKG in $PACKAGES_TO_LINT; do
    echo "--------------------------------------------------"
    echo "Linting package: ./$PKG/..."

    LINT_OUTPUT_FOR_PKG=$(mktemp)

    golangci-lint run \
        --issues-exit-code=0 \
        --path-mode=abs \
        --output.text.print-issued-lines \
        --max-issues-per-linter 0 \
        --max-same-issues 0 \
        --new-from-rev=HEAD \
        "./$PKG/..." > "$LINT_OUTPUT_FOR_PKG" 2>&1

    # Filter the combined output for issues in staged files
    for FILE in $STAGED_GO_FILES; do
        # Check if the file belongs to the current package being linted
        if [[ "$(dirname "$FILE")" == "$PKG" ]]; then
            ABS_FILE_PATH="$REPO_ROOT/$FILE"
            # Grep for the absolute path of the staged file in the linter output.
            # The final report will sort and unique the issues.
            if grep -E "^$ABS_FILE_PATH:" "$LINT_OUTPUT_FOR_PKG" >> "$ISSUES_OUTPUT"; then
                HAS_ISSUES_IN_STAGED_FILES=true
            fi
        fi
    done

    rm -f "$LINT_OUTPUT_FOR_PKG"
done

# Report final results
echo "--------------------------------------------------"
if [[ "$HAS_ISSUES_IN_STAGED_FILES" == "true" ]]; then
    echo
    echo "❌ golangci-lint found issues in staged files:"
    echo
    # Sort and unique the collected issues before printing
    sort -u "$ISSUES_OUTPUT" | sed 's/^/  /'
    echo
    echo "Please fix the issues above and re-stage the files."
    echo "To bypass this check for this commit, run:"
    echo "  CL_DEV_PRECOMMIT_LINT=false git commit"
    echo
    exit 1
else
    echo "✅ All staged Go files are clean."
    exit 0
fi
