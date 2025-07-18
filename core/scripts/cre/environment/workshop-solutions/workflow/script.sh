#!/bin/bash

START_DIR=$(pwd)

### Sanity check: stop all running docker containers.
stop_docker_containers() {
  echo "Checking for running Docker containers..."
  if [ "$(docker ps -q)" ]; then
     echo "Stopping all running Docker containers..."
     docker ps -q | xargs -r docker stop
  else
     echo "No running Docker containers to stop."
  fi
}

prepare_capabilities_repo() {
  echo "Preparing capabilities repository..."
  cd ./capabilities || exit 1
  if ! git fetch; then
    echo "Failed to git fetch capabilities repository."
    exit 1
  fi
  echo "Checking out capabilities workshop branch..."
  if ! git checkout dx-1407-workshop-cre; then
    echo "Failed to checkout capabilities workshop branch."
    exit 1
  fi
  echo "Pulling latest changes in capabilities repository..."
  if ! git pull; then
    echo "Failed to pull latest changes in capabilities repository."
    exit 1
  fi
  cd ./cron || exit 1
  echo "Cleaning up Go modules in capabilities repository..."
  if ! go mod tidy; then
    echo "Failed to tidy Go modules."
    exit 1
  fi
  echo "Building cron binary for capabilities repository..."
  if ! GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o cron; then
    echo "Failed to build cron binary."
    exit 1
  fi
  echo "Capabilities repository prepared for workshop successfully."
  cd "$START_DIR" || exit 1
}

prepare_core_repo() {
  echo "Preparing core repository..."
  cd ./chainlink || exit 1
  if ! git fetch; then
    echo "Failed to git fetch core repository."
    exit 1
  fi
  echo "Checking out core workshop branch..."
  if ! git checkout dx-1407-workshop-cre; then
    echo "Failed to checkout core workshop branch."
    exit 1
  fi
  echo "Pulling latest changes in core repository..."
  if ! git pull; then
    echo "Failed to pull latest changes in core repository."
    exit 1
  fi
  echo "Core repository prepared for workshop successfully."
  cd "$START_DIR" || exit 1
}

start_local_cre() {
  echo "Starting local CRE environment..."
  cd ./chainlink/core/scripts/cre/environment || exit 1
  if ! go run . env start; then
    echo "Failed to start local CRE environment."
    exit 1
  fi
  echo "Starting Beholder..."
  if ! go run . env beholder start; then
    echo "Failed to start local CRE Beholder."
    exit 1
  fi
  echo "Local CRE environment started successfully."
  cd "$START_DIR" || exit 1
}

# run_example_workflow will start the local CRE environment and run the example workflow.
run_example_workflow() {
  echo "Running example workflow..."
  cd ./chainlink/core/scripts/cre/environment || exit 1
  if ! go run . env start --with-example --with-beholder; then
    echo "Failed to run the example workflow."
    exit 1
  fi
  echo "Example workflow executed successfully."
  cd "$START_DIR" || exit 1
}

deploy_example_workflow() {
  echo "Deploying example workflow..."
  cd ./chainlink/core/scripts/cre/environment || exit 1
  if ! go run . env workflow deploy -w ./examples/workflows/v2/cron/main.go; then
    echo "Failed to deploy the example workflow."
    exit 1
  fi
  echo "Example workflow deployed successfully."
  cd "$START_DIR" || exit 1
}

shutdown_local_cre() {
  echo "Shutting down local CRE environment..."
  cd ./chainlink/core/scripts/cre/environment || exit 1
  if ! go run . env stop; then
    echo "Failed to stop local CRE environment."
    exit 1
  fi
  echo "Local CRE environment stopped successfully."
  cd "$START_DIR" || exit 1
}

case "$1" in
  stop_docker)
    stop_docker_containers
    ;;
  run_example_workflow)
    run_example_workflow
    ;;
  prepare_capabilities_repo)
    prepare_capabilities_repo
    ;;
  prepare_core_repo)
    prepare_core_repo
    ;;
  start_local_cre)
    start_local_cre
    ;;
  deploy_example_workflow)
    deploy_example_workflow
    ;;
  shutdown_local_cre)
    shutdown_local_cre
    ;;
  *)
    echo "Usage: $0 {stop_docker|run_example_workflow|prepare_capabilities_repo|prepare_core_repo|start_local_cre|deploy_example_workflow|shutdown_local_cre}"
    exit 1
    ;;
esac