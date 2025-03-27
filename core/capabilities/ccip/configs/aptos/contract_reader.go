package aptosconfig

import (
	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
)

func GetChainReaderConfig() (chainreader.ChainReaderConfig, error) {
	return chainreader.ChainReaderConfig{
		Modules: map[string]*chainreader.ChainReaderModule{
			// TODO: more offramp config and other modules
			consts.ContractNameOffRamp: {
				Name: "offramp",
				Functions: map[string]*chainreader.ChainReaderFunction{
					consts.MethodNameGetExecutionState: {
						Name: "get_execution_state",
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "sourceChainSelector",
								Type:     "u64",
								Required: true,
							},
							{
								Name:     "sequenceNumber",
								Type:     "u64",
								Required: true,
							},
						},
					},
					consts.MethodNameGetMerkleRoot: {
						Name: "get_merkle_root",
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "root",
								Type:     "vector<u8>",
								Required: true,
							},
						},
					},
					consts.MethodNameOffRampLatestConfigDetails: {
						Name: "latest_config_details",
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "ocrPluginType",
								Type:     "u8",
								Required: true,
							},
						},
						// TODO: change to ocr config struct, field renames
					},
					consts.MethodNameGetLatestPriceSequenceNumber: {
						Name:   "get_latest_price_sequence_number",
						Params: []chainreader.AptosFunctionParam{},
					},
					consts.MethodNameOffRampGetStaticConfig: {
						Name:   "get_static_config",
						Params: []chainreader.AptosFunctionParam{},
						// TODO: field renames
					},
					consts.MethodNameOffRampGetDynamicConfig: {
						Name:   "get_dynamic_config",
						Params: []chainreader.AptosFunctionParam{},
						// TODO: field renames
					},
					consts.MethodNameGetSourceChainConfig: {
						Name: "get_source_chain_config",
						Params: []chainreader.AptosFunctionParam{
							{
								Name:     "sourceChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
				Events: map[string]*chainreader.ChainReaderEvent{
					consts.EventNameExecutionStateChanged: {
						EventHandleStructName: "OffRampState",
						EventHandleFieldName:  "execution_state_changed_events",
						EventAccountAddress:   "state_object::get_object_address",
						EventFieldRenames: map[string]chainreader.RenamedField{
							"source_chain_selector": {
								NewName: "SourceChainSelector",
							},
							"sequence_number": {
								NewName: "SequenceNumber",
							},
							"message_id": {
								NewName: "MessageId",
							},
							"message_hash": {
								NewName: "MessageHash",
							},
							"state": {
								NewName: "State",
							},
						},
					},
					consts.EventNameCommitReportAccepted: {
						EventHandleStructName: "OffRampState",
						EventHandleFieldName:  "commit_report_accepted_events",
						EventAccountAddress:   "state_object::get_object_address",
						EventFieldRenames: map[string]chainreader.RenamedField{
							"blessed_merkle_roots": {
								NewName: "BlessedMerkleRoots",
								SubFieldRenames: map[string]chainreader.RenamedField{
									"source_chain_selector": {
										NewName: "SourceChainSelector",
									},
									"on_ramp_address": {
										NewName: "OnRampAddress",
									},
									"min_seq_nr": {
										NewName: "MinSeqNr",
									},
									"max_seq_nr": {
										NewName: "MaxSeqNr",
									},
									"merkle_root": {
										NewName: "MerkleRoot",
									},
								},
							},
							"unblessed_merkle_roots": {
								NewName: "UnblessedMerkleRoots",
								SubFieldRenames: map[string]chainreader.RenamedField{
									"source_chain_selector": {
										NewName: "SourceChainSelector",
									},
									"on_ramp_address": {
										NewName: "OnRampAddress",
									},
									"min_sequence_number": {
										NewName: "MinSeqNr",
									},
									"max_sequence_number": {
										NewName: "MaxSeqNr",
									},
									"merkle_root": {
										NewName: "MerkleRoot",
									},
								},
							},
							"price_updates": {
								NewName: "PriceUpdates",
								SubFieldRenames: map[string]chainreader.RenamedField{
									"token_price_updates": {
										NewName: "TokenPriceUpdates",
										SubFieldRenames: map[string]chainreader.RenamedField{
											"source_token": {
												NewName: "SourceToken",
											},
											"usd_per_token": {
												NewName: "UsdPerToken",
											},
										},
									},
									"gas_price_updates": {
										NewName: "GasPriceUpdates",
										SubFieldRenames: map[string]chainreader.RenamedField{
											"dest_chain_selector": {
												NewName: "DestChainSelector",
											},
											"usd_per_unit_gas": {
												NewName: "UsdPerUnitGas",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
