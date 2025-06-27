package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
)

const (
	// Sepolia chain ID
	SepoliaChainID = 11155111

	// Default OCR3 timing parameters (adjust as needed)
	DeltaProgress                           = 2 * time.Second
	DeltaResend                             = 20 * time.Second
	DeltaInitial                            = 400 * time.Millisecond
	DeltaRound                              = 500 * time.Millisecond
	DeltaGrace                              = 250 * time.Millisecond
	DeltaCertifiedCommitRequest             = 300 * time.Millisecond
	DeltaStage                              = 1 * time.Minute
	RMax                                    = 100
	MaxDurationInitialization               = 250 * time.Millisecond
	MaxDurationQuery                        = 1 * time.Second
	MaxDurationObservation                  = 1 * time.Second
	MaxDurationShouldAcceptAttestedReport   = 1 * time.Second
	MaxDurationShouldTransmitAcceptedReport = 1 * time.Second

	// Default secure mint plugin config
	MaxChains = 5
)

type Config struct {
	// Network configuration
	RPCURL     string
	PrivateKey string
	ChainID    int64

	// Oracle configuration
	OracleAddresses  []string // List of oracle addresses (transmitters)
	OraclePublicKeys []string // List of oracle public keys (signers)
	FaultTolerance   int      // Number of faulty nodes the system can tolerate

	// Contract configuration
	ConfigID string // Config ID for the configurator (32-byte hex string)

	// Gas configuration
	GasPrice *big.Int
	GasLimit uint64
}

func main() {
	config := parseFlags()

	// Validate configuration
	if err := validateConfig(config); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Deploy and configure the OCR3 Configurator
	if err := deployAndConfigureOCR3Configurator(config); err != nil {
		log.Fatalf("Failed to deploy and configure OCR3 Configurator: %v", err)
	}

	log.Println("✅ OCR3 Configurator deployment and configuration completed successfully!")
}

func parseFlags() *Config {
	config := &Config{}
	config.RPCURL = "https://rpcs.cldev.sh/ethereum/sepolia"

	// Network flags
	flag.StringVar(&config.PrivateKey, "private-key", "", "Private key for deployment (required)")
	flag.Int64Var(&config.ChainID, "chain-id", SepoliaChainID, "Chain ID (default: 11155111 for Sepolia)")

	// Oracle configuration flags
	oracleAddrs := flag.String("oracle-addresses", "", "Comma-separated list of oracle transmitter addresses (required)")
	oraclePubKeys := flag.String("oracle-public-keys", "", "Comma-separated list of oracle public keys (required)")
	flag.IntVar(&config.FaultTolerance, "f", 1, "Number of faulty nodes the system can tolerate")

	// Contract configuration
	flag.StringVar(&config.ConfigID, "config-id", "0x0000000000000000000000000000000000000000000000000000000000000001", "Config ID (32-byte hex string)")

	// Gas configuration
	gasPriceStr := flag.String("gas-price", "", "Gas price in wei (optional, will use network default if not specified)")
	flag.Uint64Var(&config.GasLimit, "gas-limit", 5000000, "Gas limit for deployment")

	flag.Parse()

	// Parse oracle addresses
	if *oracleAddrs != "" {
		config.OracleAddresses = strings.Split(*oracleAddrs, ",")
		for i, addr := range config.OracleAddresses {
			config.OracleAddresses[i] = strings.TrimSpace(addr)
		}
	}

	// Parse oracle public keys
	if *oraclePubKeys != "" {
		config.OraclePublicKeys = strings.Split(*oraclePubKeys, ",")
		for i, key := range config.OraclePublicKeys {
			config.OraclePublicKeys[i] = strings.TrimSpace(key)
		}
	}

	// Parse gas price
	if *gasPriceStr != "" {
		gasPrice, ok := new(big.Int).SetString(*gasPriceStr, 10)
		if !ok {
			log.Fatal("Invalid gas price format")
		}
		config.GasPrice = gasPrice
	}

	return config
}

func validateConfig(config *Config) error {
	if config.RPCURL == "" {
		return fmt.Errorf("RPC URL is required")
	}
	if config.PrivateKey == "" {
		return fmt.Errorf("Private key is required")
	}
	if len(config.OracleAddresses) == 0 {
		return fmt.Errorf("At least one oracle address is required")
	}
	if len(config.OraclePublicKeys) == 0 {
		return fmt.Errorf("At least one oracle public key is required")
	}
	if len(config.OracleAddresses) != len(config.OraclePublicKeys) {
		return fmt.Errorf("Number of oracle addresses must match number of oracle public keys")
	}
	if config.FaultTolerance <= 0 {
		return fmt.Errorf("Fault tolerance must be positive")
	}
	if len(config.OracleAddresses) <= 3*config.FaultTolerance {
		return fmt.Errorf("Number of oracles (%d) must be greater than 3*f (%d)", len(config.OracleAddresses), 3*config.FaultTolerance)
	}

	// Validate config ID format
	if !strings.HasPrefix(config.ConfigID, "0x") {
		return fmt.Errorf("Config ID must be a hex string starting with 0x")
	}
	if len(config.ConfigID) != 66 { // 0x + 64 hex chars
		return fmt.Errorf("Config ID must be 32 bytes (64 hex characters)")
	}

	return nil
}

func deployAndConfigureOCR3Configurator(config *Config) error {
	ctx := context.Background()

	// Connect to Ethereum client
	client, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	if chainID.Int64() != config.ChainID {
		return fmt.Errorf("chain ID mismatch: expected %d, got %d", config.ChainID, chainID.Int64())
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(config.PrivateKey, "0x"))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas price if specified
	if config.GasPrice != nil {
		auth.GasPrice = config.GasPrice
	}

	// Set gas limit
	auth.GasLimit = config.GasLimit

	// Get deployer address
	deployerAddress := auth.From
	log.Printf("Deploying from address: %s", deployerAddress.Hex())

	// Check deployer balance
	balance, err := client.BalanceAt(ctx, deployerAddress, nil)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	// balance is in wei, convert to ETH
	balanceETH := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
	log.Printf("Deployer balance: %f ETH", balanceETH)

	// Step 1: Deploy Configurator contract
	log.Println("🚀 Deploying OCR3 Configurator contract...")
	configuratorAddress, tx, _, err := configurator.DeployConfigurator(auth, client)
	if err != nil {
		return fmt.Errorf("failed to deploy configurator contract: %w", err)
	}

	log.Printf("Configurator deployment transaction: %s", tx.Hash().Hex())
	log.Printf("Configurator contract address: %s", configuratorAddress.Hex())

	// Wait for deployment transaction to be mined
	log.Println("⏳ Waiting for deployment transaction to be mined...")
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for deployment transaction: %w", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("deployment transaction failed")
	}

	log.Printf("✅ Configurator contract deployed successfully at block %d", receipt.BlockNumber)

	// // Step 2: Prepare oracle configuration
	// log.Println("🔧 Preparing oracle configuration...")

	// // Parse oracle addresses
	// oracleAddresses := make([]common.Address, len(config.OracleAddresses))
	// for i, addrStr := range config.OracleAddresses {
	// 	if !common.IsHexAddress(addrStr) {
	// 		return fmt.Errorf("invalid oracle address: %s", addrStr)
	// 	}
	// 	oracleAddresses[i] = common.HexToAddress(addrStr)
	// }

	// // Parse oracle public keys
	// oraclePublicKeys := make([][]byte, len(config.OraclePublicKeys))
	// for i, keyStr := range config.OraclePublicKeys {
	// 	keyBytes, err := hex.DecodeString(strings.TrimPrefix(keyStr, "0x"))
	// 	if err != nil {
	// 		return fmt.Errorf("invalid oracle public key %s: %w", keyStr, err)
	// 	}
	// 	if len(keyBytes) != ed25519.PublicKeySize {
	// 		return fmt.Errorf("oracle public key %s must be %d bytes", keyStr, ed25519.PublicKeySize)
	// 	}
	// 	oraclePublicKeys[i] = keyBytes
	// }

	// // Create oracle identities for OCR3 config helper
	// oracles := make([]confighelper.OracleIdentityExtra, len(oracleAddresses))
	// for i := range oracleAddresses {
	// 	oracles[i] = confighelper.OracleIdentityExtra{
	// 		OracleIdentity: confighelper.OracleIdentity{
	// 			OnchainPublicKey:  oraclePublicKeys[i],
	// 			TransmitAccount:   confighelper.Account(oracleAddresses[i].Hex()),
	// 			OffchainPublicKey: oraclePublicKeys[i], // Using same key for simplicity
	// 			PeerID:            fmt.Sprintf("oracle-%d", i), // Placeholder peer ID
	// 		},
	// 		ConfigEncryptionPublicKey: oraclePublicKeys[i], // Using same key for simplicity
	// 	}
	// }

	// // Step 3: Prepare secure mint plugin configuration
	// log.Println("🔧 Preparing secure mint plugin configuration...")

	// smPluginConfig := por.PorOffchainConfig{MaxChains: MaxChains}
	// smPluginConfigBytes, err := smPluginConfig.Serialize()
	// if err != nil {
	// 	return fmt.Errorf("failed to serialize secure mint plugin config: %w", err)
	// }

	// // Step 4: Prepare onchain configuration
	// log.Println("🔧 Preparing onchain configuration...")

	// onchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
	// 	Version:                 1,
	// 	PredecessorConfigDigest: nil,
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to encode onchain config: %w", err)
	// }

	// // Step 5: Generate OCR3 configuration
	// log.Println("🔧 Generating OCR3 configuration...")

	// signers, _, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
	// 	DeltaProgress,
	// 	DeltaResend,
	// 	DeltaInitial,
	// 	DeltaRound,
	// 	DeltaGrace,
	// 	DeltaCertifiedCommitRequest,
	// 	DeltaStage,
	// 	RMax,
	// 	[]int{len(oracles)},
	// 	oracles,
	// 	smPluginConfigBytes,
	// 	nil, // maxDurationInitialization
	// 	MaxDurationQuery,
	// 	MaxDurationObservation,
	// 	MaxDurationShouldAcceptAttestedReport,
	// 	MaxDurationShouldTransmitAcceptedReport,
	// 	uint8(config.FaultTolerance),
	// 	onchainConfig,
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to generate OCR3 configuration: %w", err)
	// }

	// // Step 6: Prepare signer keys and transmitters
	// log.Println("🔧 Preparing signer keys and transmitters...")

	// signerKeys := make([][]byte, len(signers))
	// for i, signer := range signers {
	// 	signerKeys[i] = signer
	// }

	// transmitters := make([][32]byte, len(oracleAddresses))
	// for i := range oracleAddresses {
	// 	copy(transmitters[i][:], oracleAddresses[i].Bytes())
	// }

	// // Step 7: Parse config ID
	// configIDBytes := common.FromHex(config.ConfigID)
	// var configID [32]byte
	// copy(configID[:], configIDBytes)

	// // Step 8: Set production configuration
	// log.Println("🚀 Setting production configuration...")

	// tx, err = configuratorContract.SetProductionConfig(
	// 	auth,
	// 	configID,
	// 	signerKeys,
	// 	transmitters,
	// 	f,
	// 	outOnchainConfig,
	// 	offchainConfigVersion,
	// 	offchainConfig,
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to set production config: %w", err)
	// }

	// log.Printf("Configuration transaction: %s", tx.Hash().Hex())

	// // Wait for configuration transaction to be mined
	// log.Println("⏳ Waiting for configuration transaction to be mined...")
	// receipt, err = bind.WaitMined(ctx, client, tx)
	// if err != nil {
	// 	return fmt.Errorf("failed to wait for configuration transaction: %w", err)
	// }

	// if receipt.Status == 0 {
	// 	return fmt.Errorf("configuration transaction failed")
	// }

	// log.Printf("✅ Configuration set successfully at block %d", receipt.BlockNumber)

	// // Step 9: Verify configuration
	// log.Println("🔍 Verifying configuration...")

	// // Get the latest config details
	// latestConfigDetails, err := configuratorContract.LatestConfigDetails(nil)
	// if err != nil {
	// 	return fmt.Errorf("failed to get latest config details: %w", err)
	// }

	// log.Printf("Latest config digest: 0x%x", latestConfigDetails.ConfigDigest)
	// log.Printf("Config count: %d", latestConfigDetails.ConfigCount)
	// log.Printf("Block number: %d", latestConfigDetails.BlockNumber)

	// // Print summary
	// log.Println("\n" + strings.Repeat("=", 80))
	// log.Println("🎉 OCR3 Configurator Deployment Summary")
	// log.Println(strings.Repeat("=", 80))
	// log.Printf("Contract Address: %s", configuratorAddress.Hex())
	// log.Printf("Config ID: %s", config.ConfigID)
	// log.Printf("Number of Oracles: %d", len(oracleAddresses))
	// log.Printf("Fault Tolerance (f): %d", config.FaultTolerance)
	// log.Printf("Deployment TX: %s", tx.Hash().Hex())
	// log.Printf("Configuration TX: %s", tx.Hash().Hex())
	// log.Printf("Latest Config Digest: 0x%x", latestConfigDetails.ConfigDigest)
	// log.Println(strings.Repeat("=", 80))

	// // Save deployment info to file
	// if err := saveDeploymentInfo(configuratorAddress, config.ConfigID, latestConfigDetails.ConfigDigest); err != nil {
	// 	log.Printf("Warning: failed to save deployment info: %v", err)
	// }

	return nil
}

// func saveDeploymentInfo(contractAddress common.Address, configID string, configDigest [32]byte) error {
// 	content := fmt.Sprintf(`# OCR3 Configurator Deployment Info

// Contract Address: %s
// Config ID: %s
// Config Digest: 0x%x
// Deployment Time: %s

// ## Usage in Chainlink Node Configuration

// Add this to your Chainlink node's OCR3 job specification:

// ```yaml
// type: "offchainreporting3"
// schemaVersion: 1
// name: "secure-mint-ocr3"
// contractAddress: "%s"
// pluginType: "secure-mint"
// relay: "evm"
// chainID: "%d"
// fromAddress: "YOUR_TRANSMITTER_ADDRESS"
// p2pv2Bootstrappers: ["YOUR_BOOTSTRAP_NODE_ADDRESS"]
// keyBundleID: "YOUR_KEY_BUNDLE_ID"
// transmitterID: "YOUR_TRANSMITTER_ID"

// pluginConfig:
//   maxChains: %d
// `,
// 		contractAddress.Hex(),
// 		configID,
// 		configDigest,
// 		time.Now().Format(time.RFC3339),
// 		contractAddress.Hex(),
// 		SepoliaChainID,
// 		MaxChains,
// 	)

// 	filename := fmt.Sprintf("ocr3_configurator_deployment_%s.txt", time.Now().Format("20060102_150405"))
// 	return os.WriteFile(filename, []byte(content), 0644)
// }
