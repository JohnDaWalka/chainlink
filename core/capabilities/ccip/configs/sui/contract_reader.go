package suiconfig

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	chainreaderConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"
	"golang.org/x/crypto/blake2b"
)

func PublicKeyToAddress(pubKeyHex string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return "", err
	}

	flagged := append([]byte{suikey.Ed25519Scheme}, pubKeyBytes...)

	hash := blake2b.Sum256(flagged)
	return hex.EncodeToString(hash[:]), nil
}

func GetChainReaderConfig(pubKeyStr string) (map[string]any, error) {
	fromAddress, err := PublicKeyToAddress(pubKeyStr)
	if err != nil {
		return map[string]any{}, fmt.Errorf("unable to derive Sui address from public key %s: %w", pubKeyStr, err)
	}
	fromAddress = "0x" + fromAddress

	offRampStatePointer := "_::offramp::OffRampStatePointer::off_ramp_state_id"
	onRampStatePointer := "_::onramp::OnRampStatePointer::on_ramp_state_id"

	return map[string]any{
		"IsLoopPlugin": true,
		"EventsIndexer": map[string]any{
			"PollingInterval": 10 * time.Second,
			"SyncTimeout":     10 * time.Second,
		},
		"TransactionsIndexer": map[string]any{
			"PollingInterval": 10 * time.Second,
			"SyncTimeout":     10 * time.Second,
		},
		"Modules": map[string]any{
			// TODO: more offramp config and other modules
			consts.ContractNameRMNRemote: map[string]any{
				"Name": "rmn_remote",
				"Functions": map[string]any{
					consts.MethodNameGetReportDigestHeader: map[string]any{
						"Name": "get_report_digest_header",
					},
					consts.MethodNameGetVersionedConfig: map[string]any{
						"Name": "get_versioned_config",
						// ref: https://github.com/smartcontractkit/chainlink-ccip/blob/bee7c32c71cf0aec594c051fef328b4a7281a1fc/pkg/reader/ccip.go#L1440
						"ResultTupleToStruct": []string{"version", "config"},
					},
					consts.MethodNameGetCursedSubjects: map[string]any{
						"Name": "get_cursed_subjects",
					},
				},
			},
			consts.ContractNameRMNProxy: map[string]any{
				"Name": "rmn_remote",
				"Functions": map[string]any{
					consts.MethodNameGetARM: map[string]any{
						"Name": "get_arm",
					},
				},
			},
			consts.ContractNameFeeQuoter: map[string]any{
				"Name": "fee_quoter",
				"Functions": map[string]any{
					consts.MethodNameFeeQuoterGetTokenPrice: map[string]any{
						"Name": "get_token_price",
						"Params": []map[string]any{
							{
								"Name":     "token",
								"Type":     "address",
								"Required": true,
							},
						},
					},
					consts.MethodNameFeeQuoterGetTokenPrices: map[string]any{
						"Name": "get_token_prices",
						"Params": []map[string]any{
							{
								"Name":     "tokens",
								"Type":     "vector<address>",
								"Required": true,
							},
						},
					},
					consts.MethodNameFeeQuoterGetStaticConfig: map[string]any{
						"Name": "get_static_config",
					},
					consts.MethodNameGetFeePriceUpdate: map[string]any{
						"Name": "get_dest_chain_gas_price",
						"Params": []map[string]any{
							{
								"Name":     "destChainSelector",
								"Type":     "u64",
								"Required": true,
							},
						},
					},
				},
			},
			"OffRamp": map[string]any{
				"Name": "offramp",
				"Functions": map[string]*chainreaderConfig.ChainReaderFunction{
					consts.MethodNameOffRampLatestConfigDetails: {
						Name:          "latest_config_details",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "off_ramp_state_id",
								PointerTag: &offRampStatePointer,
								Type:       "object_id",
								Required:   true,
							},
							{
								Name:     "ocrPluginType",
								Type:     "u8",
								Required: true,
							},
						},
						ResultTupleToStruct: []string{"ocr_config"},
					},
					consts.MethodNameGetLatestPriceSequenceNumber: {
						Name:          "get_latest_price_sequence_number",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "off_ramp_state_id",
								PointerTag: &offRampStatePointer,
								Type:       "object_id",
								Required:   true,
							},
						},
					},

					consts.MethodNameOffRampGetStaticConfig: {
						Name:          "get_static_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "off_ramp_state_id",
								PointerTag: &offRampStatePointer,
								Type:       "object_id",
								Required:   true,
							},
						},
						ResultTupleToStruct: []string{"static_config"},
					},
					consts.MethodNameOffRampGetDynamicConfig: {
						Name:          "get_dynamic_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "off_ramp_state_id",
								PointerTag: &offRampStatePointer,
								Type:       "object_id",
								Required:   true,
							},
						},
						ResultTupleToStruct: []string{"dynamic_config"},
					},
					consts.MethodNameGetSourceChainConfig: {
						Name:          "get_source_chain_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "off_ramp_state_id",
								PointerTag: &offRampStatePointer,
								Type:       "object_id",
								Required:   true,
							},
							{
								Name:     "sourceChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
						// is_enabled:true is_rmn_verification_disabled:true min_seq_nr:1 on_ramp
						//     "msg": "MARSHAL BYTES: map[is_enabled:true is_rmn_verification_disabled:true min_seq_nr:1 on_ramp:[78 210 119 77 36 141 82 28 29 50 193 87 62 211 250 142 0 43 137 36] router:[71 145 97 186 101 79 170 178 238 182 160 140 157 245 196 235 206 177 221 167 225 231 209 122 247 22 52 218 63 101 87 77]]",
						ResultTupleToStruct: []string{"Router", "IsEnabled", "MinSeqNr", "IsRMNVerificationDisabled", "OnRamp"},
					},
					// consts.MethodNameGetExecutionState:{
					// 	"Name": "get_execution_state",
					// 	"Params": []map[string]any{
					// 		{
					// 			"Name":     "sourceChainSelector",
					// 			"Type":     "u64",
					// 			"Required": true,
					// 		},
					// 		{
					// 			"Name":     "sequenceNumber",
					// 			"Type":     "u64",
					// 			"Required": true,
					// 		},
					// 	},
					// },
					// consts.MethodNameGetMerkleRoot: map[string]any{
					// 	"Name": "get_merkle_root",
					// 	"Params": []map[string]any{
					// 		{
					// 			"Name":     "root",
					// 			"Type":     "vector<u8>",
					// 			"Required": true,
					// 		},
					// 	},
					// },
				},
				"Events": map[string]any{
					consts.EventNameExecutionStateChanged: map[string]any{
						"EventHandleStructName": "OffRampState",
						"EventHandleFieldName":  "execution_state_changed_events",
						"EventAccountAddress":   "offramp::get_state_address",
						"EventFieldRenames": map[string]any{
							"source_chain_selector": map[string]any{
								"NewName": "SourceChainSelector",
							},
							"sequence_number": map[string]any{
								"NewName": "SequenceNumber",
							},
							"message_id": map[string]any{
								"NewName": "MessageId",
							},
							"message_hash": map[string]any{
								"NewName": "MessageHash",
							},
							"state": map[string]any{
								"NewName": "State",
							},
						},
					},
					consts.EventNameCommitReportAccepted: map[string]any{
						"EventHandleStructName": "OffRampState",
						"EventHandleFieldName":  "commit_report_accepted_events",
						"EventAccountAddress":   "offramp::get_state_address",
						"EventFieldRenames": map[string]any{
							"blessed_merkle_roots": map[string]any{
								"NewName": "BlessedMerkleRoots",
								"SubFieldRenames": map[string]any{
									"source_chain_selector": map[string]any{
										"NewName": "SourceChainSelector",
									},
									"on_ramp_address": map[string]any{
										"NewName": "OnRampAddress",
									},
									"min_seq_nr": map[string]any{
										"NewName": "MinSeqNr",
									},
									"max_seq_nr": map[string]any{
										"NewName": "MaxSeqNr",
									},
									"merkle_root": map[string]any{
										"NewName": "MerkleRoot",
									},
								},
							},
							"unblessed_merkle_roots": map[string]any{
								"NewName": "UnblessedMerkleRoots",
								"SubFieldRenames": map[string]any{
									"source_chain_selector": map[string]any{
										"NewName": "SourceChainSelector",
									},
									"on_ramp_address": map[string]any{
										"NewName": "OnRampAddress",
									},
									"min_seq_nr": map[string]any{
										"NewName": "MinSeqNr",
									},
									"max_seq_nr": map[string]any{
										"NewName": "MaxSeqNr",
									},
									"merkle_root": map[string]any{
										"NewName": "MerkleRoot",
									},
								},
							},
							"price_updates": map[string]any{
								"NewName": "PriceUpdates",
								"SubFieldRenames": map[string]any{
									"token_price_updates": map[string]any{
										"NewName": "TokenPriceUpdates",
										"SubFieldRenames": map[string]any{
											"source_token": map[string]any{
												"NewName": "SourceToken",
											},
											"usd_per_token": map[string]any{
												"NewName": "UsdPerToken",
											},
										},
									},
									"gas_price_updates": map[string]any{
										"NewName": "GasPriceUpdates",
										"SubFieldRenames": map[string]any{
											"dest_chain_selector": map[string]any{
												"NewName": "DestChainSelector",
											},
											"usd_per_unit_gas": map[string]any{
												"NewName": "UsdPerUnitGas",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"OnRamp": map[string]any{
				"Name": "onramp",
				"Functions": map[string]*chainreaderConfig.ChainReaderFunction{
					"OnRampGetDynamicConfig": {
						Name:          "get_dynamic_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "on_ramp_state_id",
								Type:       "object_id",
								PointerTag: &onRampStatePointer,
								Required:   true,
							},
						},
					},
					"OnRampGetStaticConfig": {
						Name:          "get_static_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "on_ramp_state_id",
								Type:       "object_id",
								PointerTag: &onRampStatePointer,
								Required:   true,
							},
						},
					},
					"OnRampGetDestChainConfig": {
						Name:          "get_dest_chain_config",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "on_ramp_state_id",
								Type:       "object_id",
								PointerTag: &onRampStatePointer,
								Required:   true,
							},
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
						ResultTupleToStruct: []string{"isEnabled", "sequenceNumber", "allowlistEnabled", "allowedSenders"},
					},
					"GetExpectedNextSequenceNumber": {
						Name:          "get_expected_next_sequence_number",
						SignerAddress: fromAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "on_ramp_state_id",
								Type:       "object_id",
								PointerTag: &onRampStatePointer,
								Required:   true,
							},
							{
								Name:     "destChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
				"Events": map[string]*chainreaderConfig.ChainReaderEvent{
					"CCIPMessageSent": {
						Name:      "CCIPMessageSent",
						EventType: "CCIPMessageSent",
						EventSelector: client.EventFilterByMoveEventModule{
							Module: "onramp",
							Event:  "CCIPMessageSent",
						},
						EventFilterRenames: map[string]string{
							"SequenceNumber": "sequenceNumber",
							"DestChain":      "destChainSelector",
							"SourceChain":    "sourceChainSelector",
						},
					},
				},
			},
		},
		"EventSyncInterval": 12 * time.Second,
		"EventSyncTimeout":  10 * time.Second,
		"TxSyncInterval":    12 * time.Second,
		"TxSyncTimeout":     10 * time.Second,
	}, nil
}
