#!/bin/bash

checkout_broken_capability() {
  echo "Checking out broken capability..."
  cd ./capabilities/cron || exit 1
  if ! git fetch && git checkout dx-1343-broken-capability; then
    echo "Failed to checkout broken capability branch."
    exit 1
  fi
  echo "Building broken cron capability binary..."
  if ! go mod tidy && GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o cron; then
    echo "Failed to build broken cron capability binary."
    exit 1
  fi
  echo "Broken cron capability checked out successfully."
  cd - || exit 1
}

start_observability_stack() {
  echo "Starting observability stack..."
  if ! ctf obs u; then
    echo "Failed to start observability stack."
    exit 1
  fi
  echo "Observability stack started successfully."
  cd - || exit 1
}

restart_local_cre() {
  echo "Restarting local CRE environment..."
  cd ./chainlink/core/scripts/cre/environment || exit 1
  if ! go run . env restart; then
    echo "Failed to restart local CRE environment."
    exit 1
  fi
  if ! go run . env beholder start; then
    echo "Failed to restart local CRE Beholder."
    exit 1
  fi
  echo "Local CRE environment restarted successfully."
  cd - || exit 1
}

shutdown_local_cre() {
  echo "Shutting down local CRE environment..."
  cd ./chainlink/core/scripts/cre/environment || exit 1
  if ! go run . env stop; then
    echo "Failed to stop local CRE environment."
    exit 1
  fi
  echo "Local CRE environment stopped successfully."
  cd - || exit 1
}

run_workflow_test() {
  echo "Running workflow test..."
  cd ./chainlink/system-tests/tests/smoke/cre || exit 1
  if ! go test -run Test_V2_Workflow_Workshop -timeout 4m; then
    echo "Failed to run workflow test."
    exit 1
  fi
  echo "Workflow test executed successfully."
  cd - || exit 1
}

compile_cron_capability() {
  echo "Compiling cron capability..."
  cd ./capabilities/cron || exit 1
  if ! GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o cron; then
    echo "Failed to compile cron capability."
    exit 1
  fi
  echo "Cron capability compiled successfully."
  cd - || exit 1
}

case "$1" in
  checkout_broken_capability)
    checkout_broken_capability
    ;;
  start_observability_stack)
    start_observability_stack
    ;;
  restart_local_cre)
    restart_local_cre
    ;;
  shutdown_local_cre)
    shutdown_local_cre
    ;;
  run_workflow_test)
    run_workflow_test
    ;;
  compile_cron_capability)
    compile_cron_capability
    ;;
  *)
    echo "Usage: $0 {checkout_broken_capability|start_observability_stack|restart_local_cre|shutdown_local_cre|run_workflow_test|compile_cron_capability}"
    exit 1
    ;;
esac