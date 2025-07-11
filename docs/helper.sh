#!/bin/bash
set -e  # Exit immediately if a command fails

# Colors for prettier output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}===== Configuration Documentation Test Script =====${NC}"

# Change to the root of the Chainlink repository
ROOT_DIR=$(git rev-parse --show-toplevel)
cd "$ROOT_DIR"

echo -e "${YELLOW}1. Verifying core.toml documentation...${NC}"
# Check core documentation files
DOC_FILES=(
  "core/config/docs/core.toml"
  "core/config/docs/chain.toml"
  "core/config/docs/cre_config.go"
)

for doc_file in "${DOC_FILES[@]}"; do
  if [ -f "$doc_file" ]; then
    echo "Checking $doc_file for proper documentation..."
    # Check for #Default annotations in Go files
    if [[ "$doc_file" == *.go ]]; then
      if ! grep -q "#Default" "$doc_file"; then
        echo -e "${RED}✗ $doc_file is missing #Default annotations${NC}"
      else
        echo -e "${GREEN}✓ $doc_file has #Default annotations${NC}"
      fi
    fi
    # For TOML files, just confirm they exist and aren't empty
    if [[ "$doc_file" == *.toml ]]; then
      if [ ! -s "$doc_file" ]; then
        echo -e "${RED}✗ $doc_file is empty${NC}"
      else
        echo -e "${GREEN}✓ $doc_file exists and has content${NC}"
      fi
    fi
  else
    echo -e "${RED}✗ Documentation file $doc_file not found${NC}"
  fi
done

echo -e "${YELLOW}2. Generating and verifying config documentation...${NC}"
make config-docs
if [ ! -f "docs/config.md" ]; then
  echo -e "${RED}✗ Config documentation not generated${NC}"
  exit 1
else
  echo -e "${GREEN}✓ Config documentation generated successfully${NC}"
fi

echo -e "${YELLOW}3. Running all config package tests...${NC}"
go test -v ./core/config/... 
go test -v ./core/services/chainlink -run "Test.*Config"
go test -v ./core/config/docs/...

echo -e "${YELLOW}4. Running TOML validation tests...${NC}"
# Run the script tests that validate configurations
go test -v ./core/scripts -run TestScripts

echo -e "${YELLOW}5. Checking test fixtures...${NC}"
# Test TOML files that should contain our new configuration
TOML_FIXTURES=(
  "core/services/chainlink/testdata/config-full.toml"
  "core/services/chainlink/testdata/config-empty-effective.toml"
  "core/services/chainlink/testdata/config-multi-chain-effective.toml"
  "core/web/resolver/testdata/config-full.toml"
  "core/web/resolver/testdata/config-empty-effective.toml" 
  "core/web/resolver/testdata/config-multi-chain-effective.toml"
)

for toml_file in "${TOML_FIXTURES[@]}"; do
  if [ -f "$toml_file" ]; then
    echo "Checking $toml_file..."
    echo -e "${GREEN}✓ $toml_file exists${NC}"
  else
    echo -e "${RED}✗ Test fixture $toml_file not found${NC}"
  fi
done

echo -e "${GREEN}===== Configuration documentation tests complete =====${NC}"
echo -e "Reminder: When adding new configuration options:"
echo -e "1. Add #Default annotation in Go files"
echo -e "2. Update core.toml and other documentation files"
echo -e "3. Add the field to all test fixtures"
echo -e "4. Run 'make config-docs' to regenerate documentation"