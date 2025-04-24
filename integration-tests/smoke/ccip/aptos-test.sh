#!/bin/bash
set -euox pipefail
# docker run -d --rm --name chainlink-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_HOST_AUTH_METHOD=trust -v $HOME/chainlink-pg-data/:/var/lib/postgresql/data -p 5432:5432 postgres:14 postgres -N 500 -B 1024MB
export CL_DATABASE_URL="postgresql://postgres:@localhost:5432/chainlink_test?sslmode=disable"
# go build -o chainlink-aptos ./cmd/chainlink-aptos/main.go
export CL_APTOS_CMD="/Users/friedemannf/git/chainlink-aptos/chainlink-aptos"

exec go test -v -run "Test_CCIP_Messaging_EVM2Aptos"