#!/bin/bash
set -euox pipefail
export TEST_LOG_LEVEL=debug
export LOG_LEVEL=debug
export SETH_LOG_LEVEL=debug
export TEST_SETH_LOG_LEVEL=debug
export CL_DATABASE_URL="postgresql://postgres:postgres@localhost:5432/chainlink_test?sslmode=disable"

cd integration-tests/smoke/ccip
exec go test -v -tags=integration -count=1 -run Test_CCIPMessaging_EVM2Ton ./...