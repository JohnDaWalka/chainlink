#!/bin/bash

set -e

# Specify the version of zksolc you want to install
# defaulting to 1.5.3 as that is used for ccip
ZKSOLC_VERSION=${1:-1.5.3}

# Define the GitHub repository and URL for the release
REPO="matter-labs/era-compiler-solidity"
GITHUB_URL="https://api.github.com/repos/$REPO/releases/tags/$ZKSOLC_VERSION"

# ASSET_NAME="zksolc-macosx-arm64-v${ZKSOLC_VERSION}"
ASSET_NAME="zksolc-linux-amd64-gnu-v${ZKSOLC_VERSION}"
# Fetch the release info using GitHub API and get the download URL for the asset
ASSET_URL=$(curl --silent "$GITHUB_URL" | jq -r ".assets[] | select(.name == \"$ASSET_NAME\") | .browser_download_url")

if [ -z "$ASSET_URL" ]; then
  echo "Error: Could not find the asset $ASSET_NAME in release $ZKSOLC_VERSION."
  exit 1
fi


echo "Removing existing link if any"
sudo rm -f /usr/local/bin/zksolc

echo "Downloading $ASSET_NAME... from $ASSET_URL"
sudo curl -L "$ASSET_URL" -o /usr/local/bin/zksolc

echo "Giving permission to execute"
sudo chmod +x /usr/local/bin/zksolc

echo "Testing zksolc"
zksolc --version
