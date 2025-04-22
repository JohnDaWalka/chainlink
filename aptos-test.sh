#!/bin/bash
set -euox pipefail
export TEST_LOG_LEVEL=debug
export LOG_LEVEL=debug
export SETH_LOG_LEVEL=debug
export TEST_SETH_LOG_LEVEL=debug
export CL_DATABASE_URL="postgresql://postgres:postgres@localhost:5432/chainlink_test?sslmode=disable"
export CL_APTOS_CMD="$HOME/c/chainlink-aptos/chainlink-aptos"

rm -vf $HOME/ram/aptos.txt $HOME/ram/aptos-log.txt $HOME/ram/loop_*

pushd /home/jkl/c/chainlink/integration-tests/smoke/ccip/logs
rm -vf *.log
popd

pushd $HOME/c/chainlink-aptos
rm -vf chainlink-aptos
go build -o chainlink-aptos ./cmd/chainlink-aptos/main.go
popd

cd integration-tests/smoke/ccip
exec go test -v -tags=integration -count=1 -run Test_CCIPMessaging_EVM2Aptos ./...
