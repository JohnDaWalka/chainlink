package ccip

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/stretchr/testify/require"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	suiutil "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
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

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
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
			TokenPackageId: suiState[sourceChain].LinkTokenAddress,
			TreasuryCapId:  suiState[sourceChain].LinkTokenTreasuryCapId,
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

	// Get Sui FQ
	ctx := testcontext.Get(t)

	// compile packages
	compiledPkg, err := suiBind.CompilePackage(contracts.CCIPOnramp, map[string]string{
		"ccip":        suiState[sourceChain].CCIPAddress,
		"ccip_onramp": suiState[sourceChain].OnRampAddress,
		"mcms":        suiState[sourceChain].MCMsAddress,
		"mcms_owner":  "0x1",
	})
	require.NoError(t, err)

	fmt.Println("ONRAMP COMPILED")
	// decode modules from base64 → [][]byte
	modules := make([][]byte, len(compiledPkg.Modules))
	for i, m := range compiledPkg.Modules {
		b, err := base64.StdEncoding.DecodeString(m)
		require.NoError(t, err)
		modules[i] = b
	}

	deps := make([]models.SuiAddressBytes, len(compiledPkg.Dependencies))
	for i, dep := range compiledPkg.Dependencies {
		// remove "0x" if present
		cleanDep := strings.TrimPrefix(dep, "0x")

		// decode the hex string to bytes
		depBytes, err := hex.DecodeString(cleanDep)
		require.NoError(t, err)

		// assign as SuiAddressBytes (which is just a []byte wrapper)
		deps[i] = models.SuiAddressBytes(depBytes)
	}

	// policy (upgrade type): 0 = compatible, 1 = add, 2 = dep-only
	policy := byte(0) // usually 0 unless you’re changing public APIs

	// digest: convert hex → []byte

	fmt.Println("ABOUT TO UPGRADE")
	// upgrade the onRamp
	b := uint64(700_000_000)
	resp, err := testhelpers.UpgradeContractDirect(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	},
		e.Env.BlockChains.SuiChains()[sourceChain].Client,
		suiState[sourceChain].OnRampUpgradeCapId,
		modules,
		deps,
		policy,
		compiledPkg.Digest,
	)
	require.NoError(t, err)

	fmt.Println("Upgrade tx digest:", resp.Digest)

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
