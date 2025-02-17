package solana

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	idl "github.com/smartcontractkit/chainlink-ccip/chains/solana"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainconfig"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter"
	solanacodec "github.com/smartcontractkit/chainlink-solana/pkg/solana/codec"
)

var ccipRouterIDL = idl.FetchCCIPRouterIDL()

const (
	destChainSelectorPath   = "Info.AbstractReports.Messages.Header.DestChainSelector"
	destTokenAddress        = "Info.AbstractReports.Messages.TokenAmounts.DestTokenAddress"
	merkleRootChainSelector = "Info.MerkleRoots.ChainSel"
)

type SolanaCommitPluginReport struct {
	ReportContext int8 //from ReportContextByteWords
	Report        []byte
}

type ReferenceAccountData struct {
	SomeField int
	Router    []byte
}

type TokenPoolAddressData struct {
	LookupTable []byte
}

// This type would exists in CCIP code base and would be just Report or ExecuteReport, not something custom for Solana only.
type ReportPreTransform struct {
	ReportContext  [2][32]byte
	Report         []byte
	Info           ccipocr3.ExecuteReportInfo
	AbstractReport ccip_offramp.ExecutionReportSingleChain
}

// Matches the specific type in Solana. Replaces the ReportPostTransform in chainlink-solana.
type SolanaReport struct {
	ReportPreTransform
	//ReportContext  [2][32]byte
	//Report         []byte
	//Info           ccipocr3.ExecuteReportInfo
	//AbstractReport ccip_offramp.ExecutionReportSingleChain
	TokenIndexes []byte
}

// Sample function that creates and adapter function for submitting a transaction to a Solana program and get the required set of accounts and lookup tables for creating the transaction.
func getWriteAdapterFunc(accountConstant string, routerProgramAddress string, someLookupTableAddress string, poolLookupTableAddress string) chainconfig.WriteAdapterFunc {
	return func(input any, adapterSupport chainconfig.SolanaAdapterSupport, writeContext chainconfig.WriteContext) (chainconfig.WriteAdapterOutput, error) {

		accountsMeta := []solana.AccountMeta{}
		executeReport := input.(ReportPreTransform)

		// Current ChainWriter configuration patterns implementation using the Adapter API.
		// Example for Account constant
		accountsMeta = append(accountsMeta, solana.AccountMeta{
			PublicKey:  adapterSupport.ToPublicKey(accountConstant),
			IsSigner:   true,
			IsWritable: false,
		})

		// Example for Account Lookup -> Get one or many accounts based on the data of a field (which may be an array) in the input args
		receiverAddress := adapterSupport.ToPublicKey(executeReport.Info.AbstractReports[0].Messages[0].Receiver.String())
		accountsMeta = append(accountsMeta, solana.AccountMeta{
			PublicKey:  receiverAddress,
			IsSigner:   true,
			IsWritable: false,
		})

		// Example for Accounts From Lookup Table -> From the lookup tables configured get all the accounts or only some based on indexes defined in the config.
		lookupTableData := adapterSupport.GetLookupTableData(adapterSupport.ToPublicKey(someLookupTableAddress))
		// We can select a few indexes from the lookup table or all of them
		accountsMeta = append(accountsMeta, solana.AccountMeta{
			PublicKey:  lookupTableData[0],
			IsSigner:   true,
			IsWritable: false,
		})

		// Example for PDA Lookups -> Defines one of [Account constant, Account Lookup, Accounts from Lookup Table] to discover one or multiple accounts. Then for all those accounts it looks for PDA accounts associated based on the seeds configured (the same seeds for all accounts). Seeds can be static or dynamic (meaning they come from the inputs)
		// Dynamic seeds can also come from reading into other accounts on chain.
		poolLookupTablePublicKey := adapterSupport.ToPublicKey(poolLookupTableAddress)
		rootAccountAddress := adapterSupport.GetLookupTableData(poolLookupTablePublicKey)[2]
		poolAccounts, err := adapterSupport.GetPDAAddresses(rootAccountAddress, []chainwriter.Seed{
			{Static: []byte("ccip_tokenpool_billing")},
			{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
			{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
		})
		if err != nil {
			return chainconfig.WriteAdapterOutput{}, err
		}
		for _, poolAccount := range poolAccounts {
			accountsMeta = append(accountsMeta, solana.AccountMeta{
				PublicKey:  poolAccount,
				IsWritable: true,
				IsSigner:   false,
			})
		}

		// Example for ArgsTransform implementation for CCIP. Current implementation makes chainlink-solana to depend on CCIP code base. In this case we brake that dependency since the adapter for CCIP-Solana would be then one depending on CCIP code base and chainlink-solana.
		// We also remove the need for the CCIP code base to required changes like ReportPreTransform/ReportPostTransform for Solana.
		offRampAddress := adapterSupport.ToPublicKey(writeContext.ToAddress)
		pdaAddresses, err := adapterSupport.GetPDAAddresses(offRampAddress, []chainwriter.Seed{
			{Static: []byte("reference_addresses")},
		})
		if err != nil {
			return chainconfig.WriteAdapterOutput{}, err
		}
		referenceProgramData := ReferenceAccountData{}
		routers := solana.PublicKeySlice{}
		for _, pdaAddress := range pdaAddresses {
			adapterSupport.GetDataAccount(pdaAddress, &referenceProgramData)
			routers = append(routers, adapterSupport.ToPublicKey(string(referenceProgramData.Router)))
		}
		routerAddress := routers[0]
		destTokenAddresses := []string{}
		for _, abstractReport := range executeReport.Info.AbstractReports {
			for _, message := range abstractReport.Messages {
				for _, tokenAmount := range message.TokenAmounts {
					destTokenAddresses = append(destTokenAddresses, tokenAmount.DestTokenAddress.String())
				}
			}

		}
		staticSeeds := []chainwriter.Seed{{Static: []byte("token_admin_registry")}}
		for _, destTokenAddress := range destTokenAddresses {
			staticSeeds = append(staticSeeds, chainwriter.Seed{
				Static: []byte(destTokenAddress),
			})
		}

		tokenAdminRegistryAndTokenPooolAddresses, err := adapterSupport.GetPDAAddresses(routerAddress, staticSeeds)
		if err != nil {
			return chainconfig.WriteAdapterOutput{}, err
		}

		// First address would have a different kind of data but for simplifying the example let's assume all of the address have the same data type
		addressData := TokenPoolAddressData{}
		tokenPoolAddresses := solana.PublicKeySlice{}
		for _, address := range tokenAdminRegistryAndTokenPooolAddresses {
			adapterSupport.GetDataAccount(address, &addressData)
			tokenPoolAddresses = append(tokenPoolAddresses, adapterSupport.ToPublicKey(string(addressData.LookupTable)))
		}

		tokenIndexes := []uint8{}
		for i, account := range accountsMeta {
			for _, address := range tokenPoolAddresses {
				if account.PublicKey == address {
					if i > 255 {
						return chainconfig.WriteAdapterOutput{}, fmt.Errorf("index %d out of range for uint8", i)
					}
					tokenIndexes = append(tokenIndexes, uint8(i)) //nolint:gosec
				}
			}
		}
		if len(tokenIndexes) != len(tokenPoolAddresses) {
			return chainconfig.WriteAdapterOutput{}, fmt.Errorf("missing token pools in accounts")
		}

		// Data Transformation to chain on-chain expected data structures / inputs
		onChainData := SolanaReport{
			ReportPreTransform: executeReport,
			TokenIndexes:       tokenIndexes,
		}

		someLookupTablePublicKey := adapterSupport.ToPublicKey(someLookupTableAddress)

		// No need to modify the report but we need to provide the collected set of accounts that must be sent with the transaction.
		return chainconfig.WriteAdapterOutput{
			Data: onChainData,
			LookupTables: map[solana.PublicKey]solana.PublicKeySlice{
				//TODO review this with Silas.
				someLookupTablePublicKey: adapterSupport.GetLookupTableData(someLookupTablePublicKey),
			},
			AccountsMeta: accountsMeta,
		}, nil
	}
}

func getCommitMethodConfig(fromAddress string, routerProgramAddress string, commonAddressesLookupTable solana.PublicKey) chainwriter.MethodConfig {
	sysvarInstructionsAddress := solana.SysVarInstructionsPubkey.String()
	return chainwriter.MethodConfig{
		FromAddress: fromAddress,
		InputModifications: []codec.ModifierConfig{
			&codec.RenameModifierConfig{
				Fields: map[string]string{"ReportContextByteWords": "ReportContext"},
			},
			&codec.RenameModifierConfig{
				Fields: map[string]string{"RawReport": "Report"},
			},
		},
		ChainSpecificName: "commit",
		// Adapter function gets configured here. It resolved the accounts required, the lookup tables to use and knows how to transform the input to the a type that can be serialized to send to the solana program.
		WriteAdapter: getWriteAdapterFunc(commonAddressesLookupTable.String(), routerProgramAddress, "SOME_LOOKUP_TABLE_ADDRESS", "POOL_LOOKUP_TABLE_ADDRESS"),
		// This is how configure the different lookup tables - This would not be needed if using the adapter function
		LookupTables: chainwriter.LookupTables{
			StaticLookupTables: []solana.PublicKey{
				commonAddressesLookupTable,
			},
		},
		// This is how configure the different patterns - This would not be needed if using the adapter function
		Accounts: []chainwriter.Lookup{
			getRouterAccountConfig(routerProgramAddress),
			chainwriter.PDALookups{
				Name: "SourceChainState",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("source_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRootChainSelector}},
				},
				IsSigner:   false,
				IsWritable: true,
			},
			chainwriter.PDALookups{
				Name: "RouterReportAccount",
				PublicKey: chainwriter.AccountConstant{
					Address:    routerProgramAddress,
					IsSigner:   false,
					IsWritable: false,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("commit_report")},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRootChainSelector}},
					{Dynamic: chainwriter.AccountLookup{
						Location: "Info.MerkleRoots.MerkleRoot",
					}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			getAuthorityAccountConstant(fromAddress),
			getSystemProgramConstant(),
			chainwriter.AccountConstant{
				Name:       "SysvarInstructions",
				Address:    sysvarInstructionsAddress,
				IsSigner:   true,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "GlobalState",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("state")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "BillingTokenConfig",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("fee_billing_token_config")},
					{Dynamic: chainwriter.AccountLookup{Location: "Info.TokenPrices.TokenID"}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "ChainConfigGasPrice",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("dest_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: merkleRootChainSelector}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
		},
		DebugIDLocation: "",
	}
}

func getExecuteMethodConfig(fromAddress string, routerProgramAddress string, commonAddressesLookupTable solana.PublicKey) chainwriter.MethodConfig {
	sysvarInstructionsAddress := solana.SysVarInstructionsPubkey.String()
	return chainwriter.MethodConfig{
		FromAddress: fromAddress,
		InputModifications: []codec.ModifierConfig{
			&codec.RenameModifierConfig{
				Fields: map[string]string{"ReportContextByteWords": "ReportContext"},
			},
			&codec.RenameModifierConfig{
				Fields: map[string]string{"RawExecutionReport": "Report"},
			},
		},
		ChainSpecificName: "execute",
		ArgsTransform:     "CCIP",
		LookupTables: chainwriter.LookupTables{
			DerivedLookupTables: []chainwriter.DerivedLookupTable{
				{
					Name: "PoolLookupTable",
					Accounts: chainwriter.PDALookups{
						Name: "TokenAdminRegistry",
						PublicKey: chainwriter.AccountConstant{
							Address: routerProgramAddress,
						},
						Seeds: []chainwriter.Seed{
							{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
						},
						IsSigner:   false,
						IsWritable: false,
						InternalField: chainwriter.InternalField{
							TypeName: "TokenAdminRegistry",
							Location: "LookupTable",
						},
					},
				},
			},
			StaticLookupTables: []solana.PublicKey{
				commonAddressesLookupTable,
			},
		},
		Accounts: []chainwriter.Lookup{
			getRouterAccountConfig(routerProgramAddress),
			chainwriter.PDALookups{
				Name: "SourceChainState",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("source_chain_state")},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "CommitReport",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("external_execution_config")},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
					{Dynamic: chainwriter.AccountLookup{
						// The seed is the merkle root of the report, as passed into the input params.
						Location: "Info.MerkleRoots.MerkleRoot",
					}},
				},
				IsSigner:   false,
				IsWritable: true,
			},
			chainwriter.PDALookups{
				Name: "ExternalExecutionConfig",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("external_execution_config")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			getAuthorityAccountConstant(fromAddress),
			getSystemProgramConstant(),
			chainwriter.AccountConstant{
				Name:       "SysvarInstructions",
				Address:    sysvarInstructionsAddress,
				IsSigner:   true,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "ExternalTokenPoolsSigner",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("external_token_pools_signer")},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.AccountLookup{
				Name:       "UserAccounts",
				Location:   "Info.AbstractReports.Message.ExtraArgsDecoded.Accounts",
				IsWritable: chainwriter.MetaBool{BitmapLocation: "Info.AbstractReports.Message.ExtraArgsDecoded.IsWritableBitmap"},
				IsSigner:   chainwriter.MetaBool{Value: false},
			},
			chainwriter.PDALookups{
				Name: "ReceiverAssociatedTokenAccount",
				PublicKey: chainwriter.AccountConstant{
					Address: solana.SPLAssociatedTokenAccountProgramID.String(),
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte(fromAddress)},
					{Dynamic: chainwriter.AccountLookup{Location: "Info.AbstractReports.Messages.Receiver"}},
					{Dynamic: chainwriter.AccountsFromLookupTable{
						LookupTableName: "PoolLookupTable",
						IncludeIndexes:  []int{6},
					}},
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "PerChainTokenConfig",
				PublicKey: chainwriter.AccountConstant{
					Address: routerProgramAddress,
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("ccip_tokenpool_billing")},
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.PDALookups{
				Name: "PoolChainConfig",
				PublicKey: chainwriter.AccountsFromLookupTable{
					LookupTableName: "PoolLookupTable",
					IncludeIndexes:  []int{2},
				},
				Seeds: []chainwriter.Seed{
					{Static: []byte("ccip_tokenpool_billing")},
					{Dynamic: chainwriter.AccountLookup{Location: destTokenAddress}},
					{Dynamic: chainwriter.AccountLookup{Location: destChainSelectorPath}},
				},
				IsSigner:   false,
				IsWritable: false,
			},
			chainwriter.AccountsFromLookupTable{
				LookupTableName: "PoolLookupTable",
				IncludeIndexes:  []int{},
			},
		},
		DebugIDLocation: "AbstractReport.Message.MessageID",
	}
}

func GetSolanaChainWriterConfig(routerProgramAddress string, commonAddressesLookupTable solana.PublicKey, fromAddress string) (chainwriter.ChainWriterConfig, error) {
	// check fromAddress
	pk, err := solana.PublicKeyFromBase58(fromAddress)
	if err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("invalid from address %s: %w", fromAddress, err)
	}

	if pk.IsZero() {
		return chainwriter.ChainWriterConfig{}, errors.New("from address cannot be empty")
	}

	// validate CCIP Router IDL, errors not expected
	var idl solanacodec.IDL
	if err = json.Unmarshal([]byte(ccipRouterIDL), &idl); err != nil {
		return chainwriter.ChainWriterConfig{}, fmt.Errorf("unexpected error: invalid CCIP Router IDL, error: %w", err)
	}

	// solConfig references the ccip_example_config.go from github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter, which is currently subject to change
	solConfig := chainwriter.ChainWriterConfig{
		Programs: map[string]chainwriter.ProgramConfig{
			"ccip-router": {
				Methods: map[string]chainwriter.MethodConfig{
					"execute": getExecuteMethodConfig(fromAddress, routerProgramAddress, commonAddressesLookupTable),
					"commit":  getCommitMethodConfig(fromAddress, routerProgramAddress, commonAddressesLookupTable),
				},
				IDL: ccipRouterIDL,
			},
		},
	}

	return solConfig, nil
}

func getRouterAccountConfig(routerProgramAddress string) chainwriter.PDALookups {
	return chainwriter.PDALookups{
		Name: "RouterAccountConfig",
		PublicKey: chainwriter.AccountConstant{
			Address: routerProgramAddress,
		},
		Seeds: []chainwriter.Seed{
			{Static: []byte("config")},
		},
		IsSigner:   false,
		IsWritable: false,
	}
}

func getAuthorityAccountConstant(fromAddress string) chainwriter.AccountConstant {
	return chainwriter.AccountConstant{
		Name:       "Authority",
		Address:    fromAddress,
		IsSigner:   true,
		IsWritable: true,
	}
}

func getSystemProgramConstant() chainwriter.AccountConstant {
	return chainwriter.AccountConstant{
		Name:       "SystemProgram",
		Address:    solana.SystemProgramID.String(),
		IsSigner:   false,
		IsWritable: false,
	}
}
