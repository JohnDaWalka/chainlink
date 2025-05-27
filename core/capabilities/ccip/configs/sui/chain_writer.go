package suiconfig

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-aptos/relayer/utils"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
)

func GetChainWriterConfig(publicKeyStr string) (map[string]any, error) {
	fromAddress, err := utils.HexPublicKeyToAddress(publicKeyStr)
	if err != nil {
		return map[string]any{}, fmt.Errorf("failed to parse Sui address from public key %s: %w", publicKeyStr, err)
	}

	fmt.Printf("DEBUG: Aptos GetChainWriterConfig: fromAddressStr=%s, pubKeyStr=%s\n", fromAddress.String(), publicKeyStr)

	return map[string]any{
		"modules": map[string]any{
			consts.ContractNameOffRamp: map[string]any{
				"name": "offramp",
				"functions": map[string]any{
					consts.MethodCommit: map[string]any{
						"name":         "commit",
						"public_key":   publicKeyStr,
						"from_address": fromAddress.String(),
						"params": []map[string]any{
							{
								"name":     "ReportContext",
								"type":     "vector<vector<u8>>",
								"required": true,
							},
							{
								"name":     "Report",
								"type":     "vector<u8>",
								"required": true,
							},
							{
								"name":     "Signatures",
								"type":     "vector<vector<u8>>",
								"required": true,
							},
						},
					},
					consts.MethodExecute: map[string]any{
						"name":         "execute",
						"public_key":   publicKeyStr,
						"from_address": fromAddress.String(),
						"params": []map[string]any{
							{
								"name":     "ReportContext",
								"type":     "vector<vector<u8>>",
								"required": true,
							},
							{
								"name":     "Report",
								"type":     "vector<u8>",
								"required": true,
							},
						},
					},
				},
			},
		},
		"fee_strategy": "default", // Assuming chainwriter.DefaultFeeStrategy is a string constant
	}, nil
}
