#!/bin/bash
# filepath: ./run_tests.sh
# Usage: ./run_tests.sh <number_of_iterations>
if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <number_of_iterations>"
    exit 1
fi

ITERATIONS=$1
RESULT_LOG="test_results.log"
rm -f "$RESULT_LOG"

for ((i=1; i<=ITERATIONS; i++)); do
    echo "=== Iteration: $i ===" | tee -a "$RESULT_LOG"
    OUTPUT=$(go test -timeout 20m -run '^TestRMN_TwoMessagesOnTwoLanesIncludingBatching$' -v ./smoke/ccip/ccip_rmn_test.go 2>&1)
    echo "$OUTPUT" | tee -a "$RESULT_LOG"
    echo -e "\n------------------------------------\n" | tee -a "$RESULT_LOG"
done