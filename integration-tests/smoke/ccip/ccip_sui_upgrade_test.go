package ccip

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/stretchr/testify/require"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	suiutil "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIP_Upgrade_Sui2EVM(t *testing.T) {
	ctx := testcontext.Get(t)

	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	fmt.Println("EVM: ", evmChainSelectors[0])
	fmt.Println("Sui: ", suiChainSelectors[0])

	sourceChain := suiChainSelectors[0]
	destChain := evmChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (EVM): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	suiSenderAddr, err := e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)

	normalizedAddr, err := suiutil.ConvertStringToAddressBytes(suiSenderAddr)
	require.NoError(t, err)

	// SUI FeeToken
	// mint link token to use as feeToken
	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.MintLinkToken{}, sui_cs.MintLinkTokenConfig{
			ChainSelector:  sourceChain,
			TokenPackageId: state.SuiChains[sourceChain].LinkTokenAddress,
			TreasuryCapId:  state.SuiChains[sourceChain].LinkTokenTreasuryCapId,
			Amount:         1000000000000, // 1000 Link with 1e9
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]
	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[linkops.MintLinkTokenOutput])
	require.True(t, ok)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(normalizedAddr[:], 32)
		out    messagingtest.TestCaseOutput
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	fmt.Println("Upgrading SUI onRamp...")
	// Upgrade sui onRamp contract
	upgradeSuiOnRamp(ctx, t, e, sourceChain)
	fmt.Println("Upgraded SUI onRamp")

	fmt.Println("Upgrading SUI offRamp...")
	// Upgrade sui offramp contract
	upgradeSuiOffRamp(ctx, t, e, sourceChain)
	fmt.Println("Upgraded SUI offRamp")

	fmt.Println("Upgrading SUI CCIP...")
	// Upgrade sui offramp contract
	upgradeCCIP(ctx, t, e, sourceChain)

	fmt.Println("Upgraded SUI CCIP")

	t.Run("Message to EVM - Should Succeed", func(t *testing.T) {
		out = messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				ExtraArgs:              nil,
				Replayed:               true,
				FeeToken:               outputMap.Objects.MintedLinkTokenObjectId,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	fmt.Printf("out: %v\n", out)
}

func upgradeSuiOnRamp(ctx context.Context, t *testing.T, e testhelpers.DeployedEnv, sourceChain uint64) {
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// compile packages
	compiledPackage, err := suiBind.CompilePackage(contracts.CCIPOnramp, map[string]string{
		"ccip":        state.SuiChains[sourceChain].CCIPAddress,
		"ccip_onramp": "0x0",
		"mcms":        state.SuiChains[sourceChain].MCMsAddress,
		"mcms_owner":  "0x1",
	})
	require.NoError(t, err)

	// decode modules from base64 -> [][]byte
	moduleBytes := make([][]byte, len(compiledPackage.Modules))
	for i, moduleBase64 := range compiledPackage.Modules {
		decoded, err := base64.StdEncoding.DecodeString(moduleBase64)
		require.NoError(t, err)

		moduleBytes[i] = decoded
	}

	depAddresses := make([]models.SuiAddress, len(compiledPackage.Dependencies))
	for i, dep := range compiledPackage.Dependencies {
		depAddresses[i] = models.SuiAddress(dep)
	}

	policy := byte(0)

	// upgrade the onRamp
	b := uint64(700_000_000)
	resp, err := testhelpers.UpgradeContractDirect(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	},

		e.Env.BlockChains.SuiChains()[sourceChain].Client,
		state.SuiChains[sourceChain].OnRampAddress,
		state.SuiChains[sourceChain].OnRampUpgradeCapId,
		moduleBytes,
		depAddresses,
		policy,
		compiledPackage.Digest,
	)
	require.NoError(t, err)

	newOnRampPkgId, err := bind.FindPackageIdFromPublishTx(*resp)
	require.NoError(t, err)

	// Add new PackageID
	onRamp, err := module_onramp.NewOnramp(state.SuiChains[sourceChain].OnRampAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	// add new pkgId to state
	_, err = onRamp.AddPackageId(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	}, suiBind.Object{Id: state.SuiChains[sourceChain].OnRampStateObjectId}, suiBind.Object{Id: state.SuiChains[sourceChain].OnRampOwnerCapObjectId}, newOnRampPkgId)
	require.NoError(t, err)
}

func upgradeSuiOffRamp(ctx context.Context, t *testing.T, e testhelpers.DeployedEnv, sourceChain uint64) {
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// compile packages
	compiledPackage, err := suiBind.CompilePackage(contracts.CCIPOfframp, map[string]string{
		"ccip":         state.SuiChains[sourceChain].CCIPAddress,
		"ccip_offramp": "0x0",
		"mcms":         state.SuiChains[sourceChain].MCMsAddress,
		"mcms_owner":   "0x1",
	})
	require.NoError(t, err)

	// decode modules from base64 -> [][]byte
	moduleBytes := make([][]byte, len(compiledPackage.Modules))
	for i, moduleBase64 := range compiledPackage.Modules {
		decoded, err := base64.StdEncoding.DecodeString(moduleBase64)
		require.NoError(t, err)

		moduleBytes[i] = decoded
	}

	depAddresses := make([]models.SuiAddress, len(compiledPackage.Dependencies))
	for i, dep := range compiledPackage.Dependencies {
		depAddresses[i] = models.SuiAddress(dep)
	}

	policy := byte(0)

	// upgrade the offramp
	b := uint64(700_000_000)
	resp, err := testhelpers.UpgradeContractDirect(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	},

		e.Env.BlockChains.SuiChains()[sourceChain].Client,
		state.SuiChains[sourceChain].OffRampAddress,
		state.SuiChains[sourceChain].OffRampUpgradeCapId,
		moduleBytes,
		depAddresses,
		policy,
		compiledPackage.Digest,
	)
	require.NoError(t, err)

	newOffRampPkgId, err := bind.FindPackageIdFromPublishTx(*resp)
	require.NoError(t, err)

	// Add new PackageID
	offRamp, err := module_offramp.NewOfframp(state.SuiChains[sourceChain].OffRampAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	// add new pkgId to state
	_, err = offRamp.AddPackageId(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	}, suiBind.Object{Id: state.SuiChains[sourceChain].OffRampStateObjectId}, suiBind.Object{Id: state.SuiChains[sourceChain].OffRampOwnerCapId}, newOffRampPkgId)
	require.NoError(t, err)
}

func upgradeCCIP(ctx context.Context, t *testing.T, e testhelpers.DeployedEnv, sourceChain uint64) {
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// compile packages
	compiledPackage, err := suiBind.CompilePackage(contracts.CCIP, map[string]string{
		"ccip":       "0x0",
		"mcms":       state.SuiChains[sourceChain].MCMsAddress,
		"mcms_owner": "0x1",
	})
	require.NoError(t, err)

	// decode modules from base64 -> [][]byte
	moduleBytes := make([][]byte, len(compiledPackage.Modules))
	for i, moduleBase64 := range compiledPackage.Modules {
		decoded, err := base64.StdEncoding.DecodeString(moduleBase64)
		require.NoError(t, err)

		moduleBytes[i] = decoded
	}

	depAddresses := make([]models.SuiAddress, len(compiledPackage.Dependencies))
	for i, dep := range compiledPackage.Dependencies {
		depAddresses[i] = models.SuiAddress(dep)
	}

	policy := byte(0)

	// upgrade the ccipPkg
	b := uint64(700_000_000)
	resp, err := testhelpers.UpgradeContractDirect(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	},

		e.Env.BlockChains.SuiChains()[sourceChain].Client,
		state.SuiChains[sourceChain].CCIPAddress,
		state.SuiChains[sourceChain].CCIPUpgradeCapObjectId,
		moduleBytes,
		depAddresses,
		policy,
		compiledPackage.Digest,
	)
	require.NoError(t, err)

	newCCIPPkgId, err := bind.FindPackageIdFromPublishTx(*resp)
	require.NoError(t, err)

	// Add new PackageID
	ccipStateObject, err := module_state_object.NewStateObject(state.SuiChains[sourceChain].CCIPAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	// add new pkgId to state
	_, err = ccipStateObject.AddPackageId(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	}, suiBind.Object{Id: state.SuiChains[sourceChain].CCIPObjectRef}, suiBind.Object{Id: state.SuiChains[sourceChain].CCIPOwnerCapObjectId}, newCCIPPkgId)
	require.NoError(t, err)
}
