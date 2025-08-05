package suiconfig

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	_ "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	chainwriter "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func GetChainWriterConfig(publicKeyStr string) (chainwriter.ChainWriterConfig, error) {
	// fromAddress, err := utils.HexPublicKeyToAddress(publicKeyStr)
	// if err != nil {
	// 	return chainwriter.ChainWriterConfig{}, fmt.Errorf("failed to parse Sui address from public key %s: %w", publicKeyStr, err)
	// }
	 
	isClockMutable := false

	return chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			consts.ContractNameOffRamp: {
				Name: "offramp",
				Functions: map[string]*chainwriter.ChainWriterFunction{
					consts.MethodCommit: {
						Name:      "commit",
						PublicKey: []byte(publicKeyStr),
						// FromAddress: fromAddress.String(),
						Params: []codec.SuiFunctionParam{
							{
								Name:     "object_ref_id",
								Type:     "object_id",
								Required: true,
							},
							{
								Name:     "off_ramp_state_id",
								Type:     "object_id",
								Required: true,
							},
							{
								Name:      "clock",
								Type:      "object_id",
								Required:  true,
								IsMutable: &isClockMutable,
							},
							{
								Name:     "ReportContext",
								Type:     "vector<vector<u8>>",
								Required: true,
							},
							{
								Name:     "Report",
								Type:     "vector<u8>",
								Required: true,
							},
							{
								Name:     "Signatures",
								Type:     "vector<vector<u8>>",
								Required: true,
							},
						},
					},
					consts.MethodExecute: {
						Name:      "execute",
						PublicKey: []byte(publicKeyStr),
						// FromAddress: fromAddress.String(),
						Params: []codec.SuiFunctionParam{
							{
								Name:     "ReportContext",
								Type:     "vector<vector<u8>>",
								Required: true,
							},
							{
								Name:     "Report",
								Type:     "vector<u8>",
								Required: true,
							},
						},
					},
				},
			},
		},
		// TODO: come back to it
		// FeeStrategy: chainwriter.DefaultFeeStrategy,
	}, nil

	// return map[string]any{
	// 	"modules": map[string]any{
	// 		consts.ContractNameOffRamp: map[string]any{
	// 			"name": "offramp",
	// 			"functions": map[string]any{
	// 				consts.MethodCommit: map[string]any{
	// 					"name":         "commit",
	// 					"public_key":   publicKeyStr,
	// 					"from_address": fromAddress.String(),
	// 					"params": []map[string]any{
	// 						{
	// 							"name":     "object_ref_id",
	// 							"type":     "object_id",
	// 							"required": true,
	// 						},
	// 						{
	// 							"name":      "off_ramp_state_id",
	// 							"type":      "object_id",
	// 							"required":  true,
	// 							"isMutable": false,
	// 						},
	// 						{
	// 							"name":     "clock",
	// 							"type":     "object_id",
	// 							"required": true,
	// 						},
	// 						{
	// 							"name":     "ReportContext",
	// 							"type":     "vector<vector<u8>>",
	// 							"required": true,
	// 						},
	// 						{
	// 							"name":     "Report",
	// 							"type":     "vector<u8>",
	// 							"required": true,
	// 						},
	// 						{
	// 							"name":     "Signatures",
	// 							"type":     "vector<vector<u8>>",
	// 							"required": true,
	// 						},
	// 					},
	// 				},
	// 				consts.MethodExecute: map[string]any{
	// 					"name":         "execute",
	// 					"public_key":   publicKeyStr,
	// 					"from_address": fromAddress.String(),
	// 					"params": []map[string]any{
	// 						{
	// 							"name":     "ReportContext",
	// 							"type":     "vector<vector<u8>>",
	// 							"required": true,
	// 						},
	// 						{
	// 							"name":     "Report",
	// 							"type":     "vector<u8>",
	// 							"required": true,
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"fee_strategy": "default", // Assuming chainwriter.DefaultFeeStrategy is a string constant
	// }, nil
}
