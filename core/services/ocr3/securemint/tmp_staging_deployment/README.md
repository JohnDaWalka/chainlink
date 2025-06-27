# OCR3 Configurator Deployment Script

This script deploys and configures the OCR3 Configurator contract on Sepolia testnet for use with the Secure Mint trigger capability.

## Prerequisites

1. **Go 1.24.1 or later**
2. **Funded Sepolia wallet** with ETH for deployment
3. **Sepolia RPC endpoint** (e.g., from Infura, Alchemy, or your own node)
4. **Oracle node information** including:
   - Oracle transmitter addresses (Ethereum addresses)
   - Oracle public keys (ed25519 public keys)

## Installation

1. Clone or download the script files
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Basic Usage

```bash
go run deploy_ocr3_configurator_sepolia.go \
  --rpc-url="https://sepolia.infura.io/v3/YOUR_PROJECT_ID" \
  --private-key="YOUR_PRIVATE_KEY" \
  --oracle-addresses="0x1234...,0x5678...,0x9abc..." \
  --oracle-public-keys="0xabcd...,0xefgh...,0xijkl..."
```

### All Available Options

```bash
go run deploy_ocr3_configurator_sepolia.go \
  --rpc-url="https://sepolia.infura.io/v3/YOUR_PROJECT_ID" \
  --private-key="YOUR_PRIVATE_KEY" \
  --chain-id=11155111 \
  --oracle-addresses="0x1234...,0x5678...,0x9abc..." \
  --oracle-public-keys="0xabcd...,0xefgh...,0xijkl..." \
  --f=1 \
  --config-id="0x0000000000000000000000000000000000000000000000000000000000000001" \
  --gas-price="20000000000" \
  --gas-limit=5000000
```

### Parameter Descriptions

| Parameter | Required | Description | Default |
|-----------|----------|-------------|---------|
| `--rpc-url` | Yes | Sepolia RPC endpoint URL | - |
| `--private-key` | Yes | Private key for deployment (with or without 0x prefix) | - |
| `--chain-id` | No | Chain ID (should be 11155111 for Sepolia) | 11155111 |
| `--oracle-addresses` | Yes | Comma-separated list of oracle transmitter addresses | - |
| `--oracle-public-keys` | Yes | Comma-separated list of oracle ed25519 public keys | - |
| `--f` | No | Number of faulty nodes the system can tolerate | 1 |
| `--config-id` | No | 32-byte config ID (hex string) | `0x0000000000000000000000000000000000000000000000000000000000000001` |
| `--gas-price` | No | Gas price in wei (optional, uses network default if not specified) | - |
| `--gas-limit` | No | Gas limit for deployment | 5000000 |

## Example with Real Values

```bash
go run deploy_ocr3_configurator_sepolia.go \
  --rpc-url="https://sepolia.infura.io/v3/your-project-id" \
  --private-key="0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef" \
  --oracle-addresses="0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6,0x8ba1f109551bD432803012645Hac136c772c3e90,0x147B8eb97fD247D06C4006D269c90C1908Fb5D54" \
  --oracle-public-keys="0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef,0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890,0x567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234" \
  --f=1
```

## What the Script Does

1. **Connects to Sepolia** using the provided RPC endpoint
2. **Validates configuration** including oracle addresses and public keys
3. **Deploys the Configurator contract** to Sepolia
4. **Configures the contract** with:
   - Oracle signers and transmitters
   - Fault tolerance (f value)
   - Secure mint plugin configuration
   - OCR3 timing parameters
5. **Verifies the deployment** by reading back configuration details
6. **Saves deployment information** to a timestamped file

## Output

The script will output:
- Deployment transaction hash
- Contract address
- Configuration transaction hash
- Latest config digest
- Deployment summary

A deployment info file will be created with the format: `ocr3_configurator_deployment_YYYYMMDD_HHMMSS.txt`

## Using the Deployed Contract

After deployment, you can use the contract address in your Chainlink node's OCR3 job specification:

```yaml
type: "offchainreporting3"
schemaVersion: 1
name: "secure-mint-ocr3"
contractAddress: "DEPLOYED_CONTRACT_ADDRESS"
pluginType: "secure-mint"
relay: "evm"
chainID: "11155111"
fromAddress: "YOUR_TRANSMITTER_ADDRESS"
p2pv2Bootstrappers: ["YOUR_BOOTSTRAP_NODE_ADDRESS"]
keyBundleID: "YOUR_KEY_BUNDLE_ID"
transmitterID: "YOUR_TRANSMITTER_ID"

pluginConfig:
  maxChains: 5
```

## Security Considerations

1. **Never commit private keys** to version control
2. **Use environment variables** for sensitive data in production
3. **Verify contract addresses** on Sepolia block explorers
4. **Test thoroughly** on testnets before mainnet deployment

## Troubleshooting

### Common Issues

1. **Insufficient balance**: Ensure your wallet has enough ETH for deployment
2. **Invalid oracle addresses**: Make sure all addresses are valid Ethereum addresses
3. **Invalid public keys**: Ensure all public keys are 32-byte ed25519 keys
4. **Network issues**: Verify your RPC endpoint is working and accessible

### Error Messages

- `"Number of oracles must be greater than 3*f"`: Increase the number of oracles or decrease the f value
- `"Invalid oracle address"`: Check the format of your oracle addresses
- `"Invalid oracle public key"`: Ensure public keys are 32 bytes and properly formatted

## Support

For issues related to:
- **Chainlink OCR3**: Check the [Chainlink documentation](https://docs.chain.link/)
- **Secure Mint**: Refer to the Chainlink core repository
- **Script issues**: Check the error messages and ensure all parameters are correct
