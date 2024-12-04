#!/bin/bash

set -e

# Specify the version of zksolc you want to install
VERSION="1.5.6"
# ASSET_NAME="zksolc-linux-amd64-gnu-${VERSION}"
# ASSET_NAME="zksolc-macosx-arm64-v1.5.6"

# # Define the GitHub repository and URL for the release
REPO="matter-labs/era-compiler-solidity"
GITHUB_URL="https://api.github.com/repos/$REPO/releases/tags/$VERSION"

# ASSET_NAME2="zksolc-linux-amd64-gnu-v${VERSION}"
# ASSET_NAME2="zksolc-macosx-arm64-v1.5.6"
ASSET_NAME="zksolc-macosx-arm64-v${VERSION}"
# ASSET_NAME="zksolc-linux-amd64-gnu-v${VERSION}"
# Fetch the release info using GitHub API and get the download URL for the asset
ASSET_URL=$(curl --silent "$GITHUB_URL" | jq -r ".assets[] | select(.name == \"$ASSET_NAME\") | .browser_download_url")

if [ -z "$ASSET_URL" ]; then
  echo "Error: Could not find the asset $ASSET_NAME in release $VERSION."
  exit 1
fi


echo "Removing existing link if any"
sudo rm -f /usr/local/bin/zksolc

echo "Downloading $ASSET_NAME... from $ASSET_URL"
sudo curl -L "$ASSET_URL" -o /usr/local/bin/zksolc

echo "Giving permission to execute"
sudo chmod +x /usr/local/bin/zksolc

echo "Testing"
zksolc --version

# # Make the downloaded file executable
# chmod +x ./$ASSET_NAME

# # Move it to a directory in your PATH (e.g., ../tools/bin)
# sudo mv ./$ASSET_NAME /usr/local/bin/zksolc

# echo "zksolc has been installed successfully"

# ZKSOLC_VERSION="v1.5.6"
# ZKSOLC_URL="https://github.com/matter-labs/era-compiler-solidity/releases/download/${ZKSOLC_VERSION}/zksolc-linux-amd64-gnu-${ZKSOLC_VERSION}"

# echo "Downloading zksolc from ${ZKSOLC_URL}..."
# curl -L "${ZKSOLC_URL}" -o "zksolc"
# echo "Extracting zksolc..."
# tar -xvzf zksolc.tar.gz
# mv zksolc /usr/local/bin/zksolc

# # Clean up
# rm zksolc.tar.gz

# echo "zksolc installation complete."