package suiconfig

import (
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
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

	return map[string]any{
		"IsLoopPlugin": true,
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
			consts.ContractNameOffRamp: map[string]any{
				"Name": "offramp",
				"Functions": map[string]any{
					consts.MethodNameGetExecutionState: map[string]any{
						"Name": "get_execution_state",
						"Params": []map[string]any{
							{
								"Name":     "sourceChainSelector",
								"Type":     "u64",
								"Required": true,
							},
							{
								"Name":     "sequenceNumber",
								"Type":     "u64",
								"Required": true,
							},
						},
					},
					consts.MethodNameGetMerkleRoot: map[string]any{
						"Name": "get_merkle_root",
						"Params": []map[string]any{
							{
								"Name":     "root",
								"Type":     "vector<u8>",
								"Required": true,
							},
						},
					},
					consts.MethodNameOffRampLatestConfigDetails: map[string]any{
						"Name": "latest_config_details",
						"Params": []map[string]any{
							{
								"Name":     "ocrPluginType",
								"Type":     "u8",
								"Required": true,
							},
						},
						// wrap the returned OCR config
						// https://github.com/smartcontractkit/chainlink-ccip/blob/bee7c32c71cf0aec594c051fef328b4a7281a1fc/pkg/reader/ccip.go#L141
						"ResultTupleToStruct": []string{"ocr_config"},
					},
					consts.MethodNameGetLatestPriceSequenceNumber: map[string]any{
						"Name": "get_latest_price_sequence_number",
					},
					consts.MethodNameOffRampGetStaticConfig: map[string]any{
						"Name": "get_static_config",
						// TODO: field renames
					},
					consts.MethodNameOffRampGetDynamicConfig: map[string]any{
						"Name": "get_dynamic_config",
						// TODO: field renames
					},
					consts.MethodNameGetSourceChainConfig: map[string]any{
						"Name": "get_source_chain_config",
						"Params": []map[string]any{
							{
								"Name":     "sourceChainSelector",
								"Type":     "u64",
								"Required": true,
							},
						},
					},
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
			"onramp": map[string]any{
				"Name": "onramp",
				"Functions": map[string]any{
					"get_dynamic_config": map[string]any{
						"Name":          "get_dynamic_config",
						"SignerAddress": fromAddress,
						"Params": []map[string]any{
							{
								"Name":     "onramp_state",
								"Type":     "object_id",
								"Required": false,
							},
						},
					},
					"get_static_config": map[string]any{
						"Name": "get_static_config",
					},
					"get_dest_chain_config": map[string]any{
						"Name":          "get_dest_chain_config",
						"SignerAddress": fromAddress,
						"Params": []map[string]any{
							{
								"Name":     "onramp_state",
								"Type":     "object_id",
								"Required": true,
							},
							{
								"Name":     "dest_chain_selector",
								"Type":     "u64",
								"Required": true,
							},
						},
						"ResultTupleToStruct": []string{"sequenceNumber", "allowListEnabled", "router"},
					},
					"get_expected_next_sequence_number": map[string]any{
						"Name": "get_expected_next_sequence_number",
						"Params": []map[string]any{
							{
								"Name":     "destChainSelector",
								"Type":     "u64",
								"Required": true,
							},
						},
					},
				},
				"Events": map[string]any{
					consts.EventNameCCIPMessageSent: map[string]any{
						"EventHandleStructName": "OnRampState",
						"EventHandleFieldName":  "ccip_message_sent_events",
						"EventAccountAddress":   "onramp::get_state_address",
						"EventFieldRenames": map[string]any{
							"dest_chain_selector": map[string]any{
								"NewName":         "DestChainSelector",
								"SubFieldRenames": nil,
							},
							"sequence_number": map[string]any{
								"NewName":         "SequenceNumber",
								"SubFieldRenames": nil,
							},
							"message": map[string]any{
								"NewName": "Message",
								"SubFieldRenames": map[string]any{
									"header": map[string]any{
										"NewName": "Header",
										"SubFieldRenames": map[string]any{
											"source_chain_selector": map[string]any{
												"NewName": "SourceChainSelector",
											},
											"dest_chain_selector": map[string]any{
												"NewName": "DestChainSelector",
											},
											"sequence_number": map[string]any{
												"NewName": "SequenceNumber",
											},
											"message_id": map[string]any{
												"NewName": "MessageID",
											},
											"nonce": map[string]any{
												"NewName": "Nonce",
											},
										},
									},
									"sender": map[string]any{
										"NewName": "Sender",
									},
									"data": map[string]any{
										"NewName": "Data",
									},
									"receiver": map[string]any{
										"NewName": "Receiver",
									},
									"extra_args": map[string]any{
										"NewName": "ExtraArgs",
									},
									"fee_token": map[string]any{
										"NewName": "FeeToken",
									},
									"fee_token_amount": map[string]any{
										"NewName": "FeeTokenAmount",
									},
									"fee_value_juels": map[string]any{
										"NewName": "FeeValueJuels",
									},
									"token_amounts": map[string]any{
										"NewName": "TokenAmounts",
										"SubFieldRenames": map[string]any{
											"source_pool_address": map[string]any{
												"NewName": "SourcePoolAddress",
											},
											"dest_token_address": map[string]any{
												"NewName": "DestTokenAddress",
											},
											"extra_data": map[string]any{
												"NewName": "ExtraData",
											},
											"amount": map[string]any{
												"NewName": "Amount",
											},
											"dest_exec_data": map[string]any{
												"NewName": "DestExecData",
											},
										},
									},
								},
							},
						},
						"EventFilterRenames": map[string]string{
							"DestChain": "DestChainSelector",
						},
					},
				},
			},
		},
	}, nil
}
