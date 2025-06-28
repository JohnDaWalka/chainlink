package solana

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	soltypes "github.com/smartcontractkit/chainlink-common/pkg/types/solana"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	idl "github.com/smartcontractkit/chainlink-ccip/chains/solana"
	ccipconsts "github.com/smartcontractkit/chainlink-ccip/pkg/consts"
)

var ccipOfframpIDL = idl.FetchCCIPOfframpIDL()
var ccipRouterIDL = idl.FetchCCIPRouterIDL()

// TODO: Remove IDL once V2 execute configs are live
//
//go:embed ccip_common.json
var ccipCommonIDL string

const (
	sourceChainSelectorPath       = "Info.AbstractReports.Messages.Header.SourceChainSelector"
	destTokenAddress              = "Info.AbstractReports.Messages.TokenAmounts.DestTokenAddress"
	tokenReceiverAddress          = "ExtraData.ExtraArgsDecoded.tokenReceiver"
	merkleRootSourceChainSelector = "Info.MerkleRoots.ChainSel"
	merkleRoot                    = "Info.MerkleRoots.MerkleRoot"
)

type ExecuteMethodConfigFunc func(string, string) soltypes.MethodConfig

func getCommitMethodConfig(fromAddress string, offrampProgramAddress string, priceOnly bool) soltypes.MethodConfig {
	chainSpecificName := "commit"
	if priceOnly {
		chainSpecificName = "commitPriceOnly"
	}
	return soltypes.MethodConfig{
		FromAddress: fromAddress,
		InputModifications: []codec.ModifierConfig{
			&codec.RenameModifierConfig{
				Fields: map[string]string{"ReportContextByteWords": "ReportContext"},
			},
			&codec.RenameModifierConfig{
				Fields: map[string]string{"RawReport": "Report"},
			},
		},
		ChainSpecificName: chainSpecificName,
		ArgsTransform:     "CCIPCommit",
		LookupTables: soltypes.LookupTables{
			DerivedLookupTables: []soltypes.DerivedLookupTable{
				getCommonAddressLookupTableConfig(offrampProgramAddress),
			},
		},
		Accounts:        buildCommitAccountsList(fromAddress, offrampProgramAddress, priceOnly),
		DebugIDLocation: "",
	}
}

func buildCommitAccountsList(fromAddress, offrampProgramAddress string, priceOnly bool) []soltypes.Lookup {
	accounts := []soltypes.Lookup{}
	accounts = append(accounts,
		getOfframpAccountConfig(offrampProgramAddress),
		getReferenceAddressesConfig(offrampProgramAddress),
	)
	if !priceOnly {
		accounts = append(accounts,
			soltypes.Lookup{
				PDALookups: &soltypes.PDALookups{
					Name:      "SourceChainState",
					PublicKey: getAddressConstant(offrampProgramAddress),
					Seeds: []soltypes.Seed{
						{Static: []byte("source_chain_state")},
						{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: merkleRootSourceChainSelector}}},
					},
					IsSigner:   false,
					IsWritable: true,
				},
			},
			soltypes.Lookup{
				PDALookups: &soltypes.PDALookups{
					Name:      "CommitReport",
					PublicKey: getAddressConstant(offrampProgramAddress),
					Seeds: []soltypes.Seed{
						{Static: []byte("commit_report")},
						{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: merkleRootSourceChainSelector}}},
						{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: merkleRoot}}},
					},
					IsSigner:   false,
					IsWritable: true,
				},
			},
		)
	}
	accounts = append(accounts,
		getAuthorityAccountConstant(fromAddress),
		getSystemProgramConstant(),
		getSysVarInstructionConstant(),
		getFeeBillingSignerConfig(offrampProgramAddress),
		getFeeQuoterProgramAccount(offrampProgramAddress),
		getFeeQuoterAllowedPriceUpdater(offrampProgramAddress),
		getFeeQuoterConfigLookup(offrampProgramAddress),
		getRMNRemoteProgramAccount(offrampProgramAddress),
		getRMNRemoteCursesLookup(offrampProgramAddress),
		getRMNRemoteConfigLookup(offrampProgramAddress),
		getGlobalStateConfig(offrampProgramAddress),
		getBillingTokenConfig(offrampProgramAddress),
		getChainConfigGasPriceConfig(offrampProgramAddress),
	)
	return accounts
}

func getExecuteMethodConfigV1(fromAddress string, offrampProgramAddress string) soltypes.MethodConfig {
	return soltypes.MethodConfig{
		FromAddress: fromAddress,
		InputModifications: []codec.ModifierConfig{
			&codec.RenameModifierConfig{
				Fields: map[string]string{"ReportContextByteWords": "ReportContext"},
			},
			&codec.RenameModifierConfig{
				Fields: map[string]string{"RawExecutionReport": "Report"},
			},
		},
		ChainSpecificName:        "execute",
		ArgsTransform:            "CCIPExecute",
		ComputeUnitLimitOverhead: 150_000,
		BufferPayloadMethod:      "CCIPExecutionReportBuffer",
		LookupTables: soltypes.LookupTables{
			DerivedLookupTables: []soltypes.DerivedLookupTable{
				{
					Name: "PoolLookupTable",
					Accounts: soltypes.Lookup{
						PDALookups: &soltypes.PDALookups{
							Name:      "TokenAdminRegistry",
							PublicKey: getRouterProgramAccount(offrampProgramAddress),
							Seeds: []soltypes.Seed{
								{Static: []byte("token_admin_registry")},
								{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: destTokenAddress}}},
							},
							IsSigner:   false,
							IsWritable: false,
							InternalField: soltypes.InternalField{
								TypeName: "TokenAdminRegistry",
								Location: "LookupTable",
								// TokenAdminRegistry is in the common program so need to provide the IDL
								IDL: ccipCommonIDL,
							},
						},
					},
					Optional: true, // Lookup table is optional if DestTokenAddress is not present in report
				},
				getCommonAddressLookupTableConfig(offrampProgramAddress),
			},
		},
		ATAs: []soltypes.ATALookup{
			{
				Location:      destTokenAddress,
				WalletAddress: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: tokenReceiverAddress}},
				TokenProgram: soltypes.Lookup{
					AccountsFromLookupTable: &soltypes.AccountsFromLookupTable{
						LookupTableName: "PoolLookupTable",
						IncludeIndexes:  []int{6},
					},
				},
				MintAddress: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: destTokenAddress}},
				Optional:    true, // ATA lookup is optional if DestTokenAddress is not present in report
			},
		},
		Accounts: []soltypes.Lookup{
			getOfframpAccountConfig(offrampProgramAddress),
			getReferenceAddressesConfig(offrampProgramAddress),
			{
				PDALookups: &soltypes.PDALookups{
					Name:      "SourceChainState",
					PublicKey: getAddressConstant(offrampProgramAddress),
					Seeds: []soltypes.Seed{
						{Static: []byte("source_chain_state")},
						{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: sourceChainSelectorPath}}},
					},
					IsSigner:   false,
					IsWritable: false,
				},
			},
			{
				PDALookups: &soltypes.PDALookups{
					Name:      "CommitReport",
					PublicKey: getAddressConstant(offrampProgramAddress),
					Seeds: []soltypes.Seed{
						{Static: []byte("commit_report")},
						{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: sourceChainSelectorPath}}},
						{Dynamic: soltypes.Lookup{
							AccountLookup: &soltypes.AccountLookup{
								// The seed is the merkle root of the report, as passed into the input params.
								Location: merkleRoot,
							}},
						},
					},
					IsSigner:   false,
					IsWritable: true,
				},
			},
			getAddressConstant(offrampProgramAddress),
			{
				PDALookups: &soltypes.PDALookups{
					Name:      "AllowedOfframp",
					PublicKey: getRouterProgramAccount(offrampProgramAddress),
					Seeds: []soltypes.Seed{
						{Static: []byte("allowed_offramp")},
						{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: sourceChainSelectorPath}}},
						{Dynamic: getAddressConstant(offrampProgramAddress)},
					},
					IsSigner:   false,
					IsWritable: false,
				},
			},
			getAuthorityAccountConstant(fromAddress),
			getSystemProgramConstant(),
			getSysVarInstructionConstant(),
			getRMNRemoteProgramAccount(offrampProgramAddress),
			getRMNRemoteCursesLookup(offrampProgramAddress),
			getRMNRemoteConfigLookup(offrampProgramAddress),
			// logic receiver and user defined messaging accounts are appended in the CCIPExecute args transform
			// user token account, token billing config, pool chain config, and pool lookup table accounts
			// are appended to the accounts list in the CCIPExecute args transform for each token transfer
		},
		DebugIDLocation: "Info.AbstractReports.Messages.Header.MessageID",
	}
}

func getExecuteMethodConfigV2(fromAddress string, _ string) soltypes.MethodConfig {
	return soltypes.MethodConfig{
		FromAddress: fromAddress,
		InputModifications: []codec.ModifierConfig{
			&codec.RenameModifierConfig{
				Fields: map[string]string{"ReportContextByteWords": "ReportContext"},
			},
			&codec.RenameModifierConfig{
				Fields: map[string]string{"RawExecutionReport": "Report"},
			},
		},
		ChainSpecificName:        "execute",
		ArgsTransform:            "CCIPExecuteV2",
		ComputeUnitLimitOverhead: 150_000,
		BufferPayloadMethod:      "CCIPExecutionReportBuffer",
		ATAs: []soltypes.ATALookup{
			{
				Location:      destTokenAddress,
				WalletAddress: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: tokenReceiverAddress}},
				MintAddress:   soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: destTokenAddress}},
				Optional:      true, // ATA lookup is optional if DestTokenAddress is not present in report
			},
		},
		// All accounts and lookup tables including the ones for messaging and token transfers are derived using an on-chain method
		// https://github.com/smartcontractkit/chainlink-ccip/blob/main/chains/solana/contracts/programs/ccip-offramp/src/instructions/v1/execute/derive.rs
		Accounts:        nil,
		DebugIDLocation: "Info.AbstractReports.Messages.Header.MessageID",
	}
}

func GetSolanaChainWriterConfig(offrampProgramAddress string, fromAddress string, configVersion *string) (soltypes.ContractWriterConfig, error) {
	// check fromAddress
	pk, err := solana.PublicKeyFromBase58(fromAddress)
	if err != nil {
		return soltypes.ContractWriterConfig{}, fmt.Errorf("invalid from address %s: %w", fromAddress, err)
	}

	if pk.IsZero() {
		return soltypes.ContractWriterConfig{}, errors.New("from address cannot be empty")
	}

	// validate CCIP Offramp IDL, errors not expected
	var offrampIDL soltypes.IDL
	if err = json.Unmarshal([]byte(ccipOfframpIDL), &offrampIDL); err != nil {
		return soltypes.ContractWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Offramp IDL, error: %w", err)
	}
	// validate CCIP Router IDL, errors not expected
	var routerIDL soltypes.IDL
	if err = json.Unmarshal([]byte(ccipRouterIDL), &routerIDL); err != nil {
		return soltypes.ContractWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}
	executeMethodConfigFunc, err := getExecuteMethodConfigByVersion(configVersion)
	if err != nil {
		return soltypes.ContractWriterConfig{}, fmt.Errorf("failed to get execute method config func for version: %w", err)
	}
	solConfig := soltypes.ContractWriterConfig{
		Programs: map[string]soltypes.ProgramConfig{
			ccipconsts.ContractNameOffRamp: {
				Methods: map[string]soltypes.MethodConfig{
					ccipconsts.MethodExecute:         executeMethodConfigFunc(fromAddress, offrampProgramAddress),
					ccipconsts.MethodCommit:          getCommitMethodConfig(fromAddress, offrampProgramAddress, false),
					ccipconsts.MethodCommitPriceOnly: getCommitMethodConfig(fromAddress, offrampProgramAddress, true),
				},
				IDL: ccipOfframpIDL,
			},
		},
	}

	return solConfig, nil
}

func getExecuteMethodConfigByVersion(version *string) (ExecuteMethodConfigFunc, error) {
	if version == nil {
		return getExecuteMethodConfigV1, nil
	}
	versionStr := *version
	switch versionStr {
	case "", types.SolanaChainWriterExecuteConfigVersionV1:
		return getExecuteMethodConfigV1, nil
	case types.SolanaChainWriterExecuteConfigVersionV2:
		return getExecuteMethodConfigV2, nil
	default:
		return nil, fmt.Errorf("unsupported execute config version: %s", versionStr)
	}
}

func getOfframpAccountConfig(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name: "OfframpAccountConfig",
			PublicKey: soltypes.Lookup{
				AccountConstant: &soltypes.AccountConstant{
					Address: offrampProgramAddress,
				},
			},
			Seeds: []soltypes.Seed{
				{Static: []byte("config")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getAddressConstant(address string) soltypes.Lookup {
	return soltypes.Lookup{
		AccountConstant: &soltypes.AccountConstant{
			Address:    address,
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getFeeQuoterProgramAccount(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      ccipconsts.ContractNameFeeQuoter,
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("reference_addresses")},
			},
			IsSigner:   false,
			IsWritable: false,
			// Reads the address from the reference addresses account
			InternalField: soltypes.InternalField{
				TypeName: "ReferenceAddresses",
				Location: "FeeQuoter",
				IDL:      ccipOfframpIDL,
			},
		},
	}
}

func getRouterProgramAccount(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      ccipconsts.ContractNameRouter,
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("reference_addresses")},
			},
			IsSigner:   false,
			IsWritable: false,
			// Reads the address from the reference addresses account
			InternalField: soltypes.InternalField{
				TypeName: "ReferenceAddresses",
				Location: "Router",
				IDL:      ccipOfframpIDL,
			},
		},
	}
}

func getReferenceAddressesConfig(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      "ReferenceAddresses",
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("reference_addresses")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getFeeBillingSignerConfig(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      "FeeBillingSigner",
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("fee_billing_signer")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getFeeQuoterAllowedPriceUpdater(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name: "FeeQuoterAllowedPriceUpdater",
			// Fetch fee quoter public key to use as program ID for PDA
			PublicKey: getFeeQuoterProgramAccount(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("allowed_price_updater")},
				{Dynamic: getFeeBillingSignerConfig(offrampProgramAddress)},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getFeeQuoterConfigLookup(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name: "FeeQuoterConfig",
			// Fetch fee quoter public key to use as program ID for PDA
			PublicKey: getFeeQuoterProgramAccount(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("config")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getRMNRemoteProgramAccount(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      ccipconsts.ContractNameRMNRemote,
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("reference_addresses")},
			},
			IsSigner:   false,
			IsWritable: false,
			// Reads the address from the reference addresses account
			InternalField: soltypes.InternalField{
				TypeName: "ReferenceAddresses",
				Location: "RmnRemote",
				IDL:      ccipOfframpIDL,
			},
		},
	}
}

func getRMNRemoteCursesLookup(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      "RMNRemoteCurses",
			PublicKey: getRMNRemoteProgramAccount(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("curses")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getRMNRemoteConfigLookup(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      "RMNRemoteConfig",
			PublicKey: getRMNRemoteProgramAccount(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("config")},
			},
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getGlobalStateConfig(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      "GlobalState",
			PublicKey: getAddressConstant(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("state")},
			},
			IsSigner:   false,
			IsWritable: true,
		},
		Optional: true,
	}
}

func getBillingTokenConfig(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      "BillingTokenConfig",
			PublicKey: getFeeQuoterProgramAccount(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("fee_billing_token_config")},
				{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: "Info.TokenPriceUpdates.TokenID"}}},
			},
			IsSigner:   false,
			IsWritable: true,
		},
		Optional: true,
	}
}

func getChainConfigGasPriceConfig(offrampProgramAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		PDALookups: &soltypes.PDALookups{
			Name:      "ChainConfigGasPrice",
			PublicKey: getFeeQuoterProgramAccount(offrampProgramAddress),
			Seeds: []soltypes.Seed{
				{Static: []byte("dest_chain")},
				{Dynamic: soltypes.Lookup{AccountLookup: &soltypes.AccountLookup{Location: "Info.GasPriceUpdates.ChainSel"}}},
			},
			IsSigner:   false,
			IsWritable: true,
		},
		Optional: true,
	}
}

// getCommonAddressLookupTableConfig returns the lookup table config that fetches the lookup table address from a PDA on-chain
// The offramp contract contains a PDA with a ReferenceAddresses struct that stores the lookup table address in the OfframpLookupTable field
func getCommonAddressLookupTableConfig(offrampProgramAddress string) soltypes.DerivedLookupTable {
	return soltypes.DerivedLookupTable{
		Name: "CommonAddressLookupTable",
		Accounts: soltypes.Lookup{
			PDALookups: &soltypes.PDALookups{
				Name:      "OfframpLookupTable",
				PublicKey: getAddressConstant(offrampProgramAddress),
				Seeds: []soltypes.Seed{
					{Static: []byte("reference_addresses")},
				},
				InternalField: soltypes.InternalField{
					TypeName: "ReferenceAddresses",
					Location: "OfframpLookupTable",
					IDL:      ccipOfframpIDL,
				},
			},
		},
	}
}

func getAuthorityAccountConstant(fromAddress string) soltypes.Lookup {
	return soltypes.Lookup{
		AccountConstant: &soltypes.AccountConstant{
			Name:       "Authority",
			Address:    fromAddress,
			IsSigner:   true,
			IsWritable: true,
		},
	}
}

func getSystemProgramConstant() soltypes.Lookup {
	return soltypes.Lookup{
		AccountConstant: &soltypes.AccountConstant{
			Name:       "SystemProgram",
			Address:    solana.SystemProgramID.String(),
			IsSigner:   false,
			IsWritable: false,
		},
	}
}

func getSysVarInstructionConstant() soltypes.Lookup {
	return soltypes.Lookup{
		AccountConstant: &soltypes.AccountConstant{
			Name:       "SysvarInstructions",
			Address:    solana.SysVarInstructionsPubkey.String(),
			IsSigner:   false,
			IsWritable: false,
		},
	}
}
