# Getting OCR tests to run against a testnet

A quick start guide.

## Pre-requisites

1. [Installed Go](https://go.dev/) - `asdf install` works well for this.
1. [Installed Docker](https://www.docker.com/). Consider [increasing resources limits needed by Docker](https://stackoverflow.com/questions/44533319/how-to-assign-more-memory-to-docker-container) as most tests require building several containers for a Decentralized Oracle Network (e.g. OCR requires 6 nodes, 6 DBs, and a mock server).
1. Build docker image

   ```bash
   make build_docker_image image=<your-image-name> tag=<your-tag>
   ```

   Example: `make build_docker_image image=chainlink tag=test-tag`

1. Some familiarity with how `go test` works.

## Config Files

There are two main configuration files that you will need to setup.

1. `~/.testsecrets`
1. `$WORKSPACE/chainlink/integration-tests/testconfig/overrides.toml`

### `~/.testsecrets`

```env
E2E_TEST_CHAINLINK_IMAGE="chainlink"

E2E_TEST_$NETWORK_RPC_WS_URL=""
E2E_TEST_$NETWORK_RPC_HTTP_URL=""
E2E_TEST_$NETWORK_WALLET_KEY=""
```

### `overrides.toml`

```toml
[ChainlinkImage]
version = "test-tag"

[Network]
selected_networks = ["SEPOLIA"]

# These are optional
# [OCR2.Contracts]
# link_token = "0x779877A7B0D9E8603169DdbD7836e478b4624789"
# offchain_aggregators = ["0xc1ce3815d6e7f3705265c2577F1342344752A5Eb"]
#
# [OCR2.Contracts.Settings."0xc1ce3815d6e7f3705265c2577F1342344752A5Eb"]
# use = true       # Default: true. Reuse existing OCR contracts?
# configure = true # Default: true. Configure existing OCR contracts?
```

## Running the tests

```bash
cd $WORKSPACE/chainlink/integration-tests
BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -p 1 -timeout 15m -run "TestOCRv2Basic" ./smoke
```
