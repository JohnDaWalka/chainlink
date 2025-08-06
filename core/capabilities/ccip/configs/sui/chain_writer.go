package suiconfig

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	_ "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	chainwriter "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func GetChainWriterConfig(publicKeyStr string) (chainwriter.ChainWriterConfig, error) {
	// rawPubKey, err := hex.DecodeString(publicKeyStr)
	// if err != nil {
	// 	return chainwriter.C, fmt.Errorf("invalid public key hex %q: %w", publicKeyStr, err)
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

}
