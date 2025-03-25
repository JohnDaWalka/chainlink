package aptosconfig

import (
	"encoding/hex"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader"
	"github.com/smartcontractkit/chainlink-aptos/relayer/chainwriter"
	"golang.org/x/crypto/sha3"
)

func GetChainWriterConfig(publicKeyStr string) (chainwriter.ChainWriterConfig, error) {
	pubKeyBytes, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("failed to decode Aptos public key %s: %w", publicKeyStr, err)
	}
	authKey := sha3.Sum256(append([]byte(pubKeyBytes), 0x00))
	fromAddressStr := fmt.Sprintf("%064x", authKey)

	var fromAddress aptos.AccountAddress
	err = fromAddress.ParseStringRelaxed(fromAddressStr)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("failed to parse Aptos from address %s: %w", fromAddressStr, err)
	}

	fmt.Printf("DEBUG: Aptos GetChainWriterConfig: fromAddressStr=%s, pubKeyStr=%s\n", fromAddressStr, publicKeyStr)

	return chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			"forwarder": {
				Functions: map[string]*chainwriter.ChainWriterFunction{
					"report": {
						PublicKey: publicKeyStr,
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "Receiver",
								Type:     "address",
								Required: true,
							},
							{
								Name:     "RawReport",
								Type:     "vector<u8>", // report_context | metadata | report
								Required: true,
							},
							{
								Name:     "Signatures",
								Type:     "vector<vector<u8>>",
								Required: true,
							},
						},
					},
				},
			},
		},
		FeeStrategy: chainwriter.DefaultFeeStrategy,
	}, nil
}
