#!/bin/bash

# Parse the test logs and extract content up to the "Condition never satisfied" error

# Create a timestamped directory
current_timestamp=$(date +"%Y%m%d_%H%M%S")
test_run_dir="test-run-${current_timestamp}"
mkdir -p "$test_run_dir"

# Read sui-test.log and stop at the first occurrence of "Error:      	Condition never satisfied"
sed '/Error:.*Condition never satisfied/q' sui-test.log > "$test_run_dir/parsed-sui-test.log"

pushd $test_run_dir

mkdir node_errors

cat ../integration-tests/smoke/ccip/logs/*.log | grep -i 'error' | grep -v -e "Solana" -e "solana" > node_errors/errors.log
cat ../sui-test.log | grep -i 'error' | grep -v -e "Solana" -e "solana" > sui-test-errors.log

popd



