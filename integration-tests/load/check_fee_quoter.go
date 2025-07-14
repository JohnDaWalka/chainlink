package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	fee_quoter "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
)

func main() {
	client, err := ethclient.Dial("wss://crib-damjan-geth-1337-ws.main.stage.cldev.sh")
	if err != nil {
		log.Fatalf("Failed to connect to the client: %v", err)
	}

	feeQuoterAddr := common.HexToAddress("0xc3e53F4d16Ae77Db1c982e75a937B9f60FE63690")
	feeQuoter, err := fee_quoter.NewFeeQuoter(feeQuoterAddr, client)
	if err != nil {
		log.Fatalf("Failed to create fee quoter instance: %v", err)
	}

	fmt.Println("=== Fee Quoter Configuration Check ===")
	fmt.Printf("Fee Quoter Address: %s\n", feeQuoterAddr.Hex())
	fmt.Println()

	// Check fee tokens
	fmt.Println("=== Fee Tokens ===")
	feeTokens, err := feeQuoter.GetFeeTokens(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		fmt.Printf("Error getting fee tokens: %v\n", err)
	} else {
		fmt.Printf("Number of fee tokens: %d\n", len(feeTokens))
		for i, token := range feeTokens {
			fmt.Printf("  %d. Token: %s\n", i+1, token.Hex())
		}
	}
	fmt.Println()

	// Check token price feed config for known tokens
	fmt.Println("=== Token Price Feed Config ===")
	tokens := []common.Address{
		common.HexToAddress("0x0165878A594ca255338adfa4d48449f69242Eb8F"), // WETH9
		common.HexToAddress("0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"), // LINK
		common.HexToAddress("0x9A676e781A523b5d0C0e43731313A708CB607508"), // WETH9 (main chain)
	}
	for _, token := range tokens {
		fmt.Printf("Checking token: %s\n", token.Hex())
		cfg, err := feeQuoter.GetTokenPriceFeedConfig(&bind.CallOpts{Context: context.Background()}, token)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  DataFeedAddress: %s\n", cfg.DataFeedAddress.Hex())
			fmt.Printf("  TokenDecimals: %d\n", cfg.TokenDecimals)
		}
	}
	fmt.Println()

	// Check destination chain config for known selectors
	fmt.Println("=== Destination Chain Config ===")
	selectors := []uint64{
		3379446385462418246,
		12922642891491394802,
		4457093679053095497,
	}
	for _, sel := range selectors {
		fmt.Printf("Checking chain selector: %d\n", sel)
		cfg, err := feeQuoter.GetDestChainConfig(&bind.CallOpts{Context: context.Background()}, sel)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  IsEnabled: %t\n", cfg.IsEnabled)
			fmt.Printf("  MaxNumberOfTokensPerMsg: %d\n", cfg.MaxNumberOfTokensPerMsg)
			fmt.Printf("  MaxDataBytes: %d\n", cfg.MaxDataBytes)
			fmt.Printf("  MaxPerMsgGasLimit: %d\n", cfg.MaxPerMsgGasLimit)
			fmt.Printf("  DestGasOverhead: %d\n", cfg.DestGasOverhead)
		}
	}
}
