package testhelpers

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_dummy_receiver"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	"github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/api"
)

type AptosTestDeployPrerequisitesChangeSet struct {
	T                   *testing.T
	AptosChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = AptosTestDeployPrerequisitesChangeSet{}

func (c AptosTestDeployPrerequisitesChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {
	t := c.T

	onchainState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: AptosTestDeployPrerequisitesChangeSet: chain selectors: %+v onchainState: %+v\n", c.AptosChainSelectors)
	for _, chainSelector := range c.AptosChainSelectors {
		aptosChainState := onchainState.AptosChains[chainSelector]
		// TODO: use a real token instead of APT
		aptosChainState.LinkTokenAddress = aptos.AccountAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10}
		err = changeset.SaveOnchainStateAptos(chainSelector, aptosChainState, e)
		require.NoError(t, err)
	}
	return deployment.ChangesetOutput{}, nil
}

type AptosTestDeployContractsChangeSet struct {
	T                   *testing.T
	HomeChainSelector   uint64
	AptosChainSelectors []uint64
	AllChainSelectors   []uint64
}

var _ commoncs.ConfiguredChangeSet = AptosTestDeployContractsChangeSet{}

func (c AptosTestDeployContractsChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {

	t := c.T

	onchainState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: AptosTestDeployContractsChangeSet: chain selectors: %+v onchainState: %+v\n", c.AptosChainSelectors, onchainState)

	for _, chainSelector := range c.AptosChainSelectors {
		aptosChain := e.AptosChains[chainSelector]
		aptosChainState := onchainState.AptosChains[chainSelector]
		c.deployAptosContracts(t, e, chainSelector, aptosChain, aptosChainState, onchainState)
	}
	return deployment.ChangesetOutput{}, nil
}

func (c AptosTestDeployContractsChangeSet) deployAptosContracts(t *testing.T, e deployment.Environment, chainSelector uint64, aptosChain deployment.AptosChain, aptosChainState changeset.AptosCCIPChainState, onchainState changeset.CCIPOnChainState) {
	logger := logger.Test(t)
	adminAddress := aptosChain.DeployerSigner.AccountAddress()

	mcmsAddress, mcmsPendingTx, mcmsBindings, err := mcms.DeployToResourceAccount(aptosChain.DeployerSigner, aptosChain.Client)
	require.NoError(t, err)
	logger.Infow("Deployed Aptos MCMS", "address", mcmsAddress.String(), "pendingTx", mcmsPendingTx.TxnHash())
	_ = mcmsBindings
	waitForTx(t, aptosChain.Client, mcmsPendingTx.TxnHash(), time.Minute*1)

	ccipAddress, ccipPendingTx, ccipBindings, err := ccip.DeployToObject(aptosChain.DeployerSigner, aptosChain.Client, mcmsAddress, false)
	require.NoError(t, err)
	logger.Infow("Deployed Aptos CCIP", "address", ccipAddress.String(), "pendingTx", ccipPendingTx.TxnHash())
	_ = ccipBindings
	waitForTx(t, aptosChain.Client, ccipPendingTx.TxnHash(), time.Minute*1)

	aptosChainState.CCIPAddress = ccipAddress

	ccipRouterPendingTx, ccipRouterBindings, err := ccip_router.DeployToExistingObject(aptosChain.DeployerSigner, aptosChain.Client, ccipAddress, ccipAddress)
	require.NoError(t, err)
	logger.Infow("Deployed Aptos CCIP Router", "address", ccipAddress.String(), "pendingTx", ccipRouterPendingTx.TxnHash())
	_ = ccipRouterBindings
	waitForTx(t, aptosChain.Client, ccipRouterPendingTx.TxnHash(), time.Minute*1)

	ccipDummyReceiverAddress, ccipDummyReceiverPendingTx, ccipDummyReceiverBindings, err := ccip_dummy_receiver.DeployToObject(aptosChain.DeployerSigner, aptosChain.Client, ccipAddress)
	require.NoError(t, err)
	logger.Infow("Deployed Aptos CCIP Dummy Receiver", "address", ccipDummyReceiverAddress.String(), "pendingTx", ccipDummyReceiverPendingTx.TxnHash())
	_ = ccipDummyReceiverBindings
	waitForTx(t, aptosChain.Client, ccipDummyReceiverPendingTx.TxnHash(), time.Minute*1)

	aptosChainState.ReceiverAddress = ccipDummyReceiverAddress

	transactOpts := &bind.TransactOpts{
		Signer: aptosChain.DeployerSigner,
	}

	pendingTx, err := ccipBindings.Onramp.Initialize(transactOpts, chainSelector, adminAddress, []uint64{}, []bool{}, []bool{})
	require.NoError(t, err)
	logger.Infow("Initialized Onramp", "pendingTx", pendingTx.TxnHash())
	waitForTx(t, aptosChain.Client, pendingTx.TxnHash(), time.Minute*1)

	// TODO: actually figure out where this value comes from
	//                                 failed to apply changeset at index 1: invalid changeset config: validate plugin info PluginType: CCIPCommit, Chains: [909606746561742123 5548718428018410741 4457093679053095497]: invalid ccip ocr params: invalid execute off-chain config: MessageVisibilityInterval=8h0m0s does not match the permissionlessExecutionThresholdSeconds in dynamic config =0 for chain 4457093679053095497
	permissionlessExecutionThresholdSecs := uint32(60 * 60 * 8)
	pendingTx, err = ccipBindings.Offramp.Initialize(transactOpts, chainSelector, permissionlessExecutionThresholdSecs, []uint64{}, []bool{}, []bool{}, [][]byte{})
	require.NoError(t, err)
	logger.Infow("Initialized Offramp", "pendingTx", pendingTx.TxnHash())
	waitForTx(t, aptosChain.Client, pendingTx.TxnHash(), time.Minute*1)

	pendingTx, err = ccipBindings.FeeQuoter.Initialize(transactOpts, uint64(1000000), aptosChainState.LinkTokenAddress, uint64(1000000), []aptos.AccountAddress{aptosChainState.LinkTokenAddress})
	logger.Infow("Initialized FeeQuoter", "pendingTx", pendingTx.TxnHash())
	require.NoError(t, err)
	waitForTx(t, aptosChain.Client, pendingTx.TxnHash(), time.Minute*1)

	pendingTx, err = ccipBindings.RMNRemote.Initialize(transactOpts, chainSelector)
	logger.Infow("Initialized RMNRemote", "pendingTx", pendingTx.TxnHash())
	require.NoError(t, err)
	waitForTx(t, aptosChain.Client, pendingTx.TxnHash(), time.Minute*1)

	logger.Infow("All Aptos contracts deployed")

	err = changeset.SaveOnchainStateAptos(chainSelector, aptosChainState, e)
	require.NoError(t, err)
}

type AptosTestConfigureContractsChangeSet struct {
	T                   *testing.T
	HomeChainSelector   uint64
	AptosChainSelectors []uint64
	AllChainSelectors   []uint64
}

var _ commoncs.ConfiguredChangeSet = AptosTestConfigureContractsChangeSet{}

func (c AptosTestConfigureContractsChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {

	t := c.T

	onchainState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: AptosTestConfigureContractsChangeSet: chain selectors: %+v onchainState: %+v\n", c.AptosChainSelectors, onchainState)

	for _, chainSelector := range c.AptosChainSelectors {
		aptosChain := e.AptosChains[chainSelector]
		aptosChainState := onchainState.AptosChains[chainSelector]
		c.configureAptosContracts(t, e, chainSelector, aptosChain, aptosChainState, onchainState)
	}
	return deployment.ChangesetOutput{}, nil
}

func (c AptosTestConfigureContractsChangeSet) configureAptosContracts(t *testing.T, e deployment.Environment, chainSelector uint64, aptosChain deployment.AptosChain, aptosChainState changeset.AptosCCIPChainState, onchainState changeset.CCIPOnChainState) {
	logger := logger.Test(t)

	ccipBindings := ccip.Bind(aptosChainState.CCIPAddress, aptosChain.Client)

	transactOpts := &bind.TransactOpts{
		Signer: aptosChain.DeployerSigner,
	}

	donID, err := internal.DonIDForChain(
		onchainState.Chains[c.HomeChainSelector].CapabilityRegistry,
		onchainState.Chains[c.HomeChainSelector].CCIPHome,
		chainSelector,
	)
	require.NoError(t, err)

	ocr3Args, err := internal.BuildSetOCR3ConfigArgsAptos(
		donID, onchainState.Chains[c.HomeChainSelector].CCIPHome, chainSelector, globals.ConfigTypeActive)
	require.NoError(t, err)

	var commitArgs *internal.MultiOCR3BaseOCRConfigArgsAptos = nil
	var execArgs *internal.MultiOCR3BaseOCRConfigArgsAptos = nil
	for _, ocr3Arg := range ocr3Args {
		if ocr3Arg.OcrPluginType == uint8(types.PluginTypeCCIPCommit) {
			commitArgs = &ocr3Arg
		} else if ocr3Arg.OcrPluginType == uint8(types.PluginTypeCCIPExec) {
			execArgs = &ocr3Arg
		} else {
			t.Fatalf("unexpected ocr3 plugin type %s", ocr3Arg.OcrPluginType)
		}
	}
	require.NotNil(t, commitArgs)
	require.NotNil(t, execArgs)

	commitSigners := [][]byte{}
	for _, signer := range commitArgs.Signers {
		commitSigners = append(commitSigners, signer.Bytes())
	}
	commitTransmitters := []aptos.AccountAddress{}
	for _, transmitter := range commitArgs.Transmitters {
		commitTransmitters = append(commitTransmitters, aptos.AccountAddress(transmitter))
	}
	ccipBindings.Offramp.SetOcr3Config(transactOpts, commitArgs.ConfigDigest[:], uint8(types.PluginTypeCCIPCommit), commitArgs.F, commitArgs.IsSignatureVerificationEnabled, commitSigners, commitTransmitters)

	execSigners := [][]byte{}
	for _, signer := range execArgs.Signers {
		execSigners = append(execSigners, signer.Bytes())
	}
	execTransmitters := []aptos.AccountAddress{}
	for _, transmitter := range execArgs.Transmitters {
		execTransmitters = append(execTransmitters, aptos.AccountAddress(transmitter))
	}
	ccipBindings.Offramp.SetOcr3Config(transactOpts, execArgs.ConfigDigest[:], uint8(types.PluginTypeCCIPCommit), execArgs.F, execArgs.IsSignatureVerificationEnabled, execSigners, execTransmitters)

	logger.Infow("Confirmed")
}

func addLaneAptosChangesets(t *testing.T, e *DeployedEnv, from, to uint64, fromFamily, toFamily string) []commoncs.ConfiguredChangeSet {
	return []commoncs.ConfiguredChangeSet{AptosTestAddLaneChangeSet{
		T:                 t,
		fromChainSelector: from,
		toChainSelector:   to,
		fromFamily:        fromFamily,
		toFamily:          toFamily,
	}}
}

type AptosTestAddLaneChangeSet struct {
	T                 *testing.T
	fromChainSelector uint64
	toChainSelector   uint64
	fromFamily        string
	toFamily          string
}

var _ commoncs.ConfiguredChangeSet = AptosTestAddLaneChangeSet{}

func (c AptosTestAddLaneChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {
	t := c.T
	// TODO: support other paths
	require.Equal(t, c.fromFamily, chainsel.FamilyEVM, "must be from EVM")
	require.Equal(t, c.toFamily, chainsel.FamilyAptos, "must be to Aptos")

	aptosSelector := c.toChainSelector
	aptosChain := e.AptosChains[aptosSelector]

	onchainState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	aptosChainState := onchainState.AptosChains[aptosSelector]

	fmt.Printf("DEBUG: AptosTestAddLaneChangeSet: LINK token: %s CCIP: %s Receiver: %s\n", aptosChainState.LinkTokenAddress.String(), aptosChainState.CCIPAddress.String(), aptosChainState.ReceiverAddress.String())

	require.NotEqual(t, aptosChainState.LinkTokenAddress, aptos.AccountZero, "LINK token address must be set")
	require.NotEqual(t, aptosChainState.CCIPAddress, aptos.AccountZero, "CCIP address must be set")
	require.NotEqual(t, aptosChainState.ReceiverAddress, aptos.AccountZero, "Receiver address must be set")

	ccipBindings := ccip.Bind(aptosChainState.CCIPAddress, aptosChain.Client)
	transactOpts := &bind.TransactOpts{
		Signer: aptosChain.DeployerSigner,
	}

	evmChainState := onchainState.Chains[c.fromChainSelector]
	evmOnrampAddress := evmChainState.OnRamp.Address().Bytes()
	fmt.Printf("DEBUG: AptosTestAddLaneChangeSet: EVM onramp address: %s\n", hex.EncodeToString(evmOnrampAddress))

	sourceChainSelectors := []uint64{c.fromChainSelector}
	sourceChainsIsEnabled := []bool{true}
	sourceChainsIsRMNVerificationDisabled := []bool{true}
	sourceChainsOnramps := [][]byte{evmOnrampAddress}
	pendingTx, err := ccipBindings.Offramp.ApplySourceChainConfigUpdates(transactOpts, sourceChainSelectors, sourceChainsIsEnabled, sourceChainsIsRMNVerificationDisabled, sourceChainsOnramps)
	require.NoError(t, err)
	waitForTx(t, aptosChain.Client, pendingTx.TxnHash(), time.Minute*1)

	fmt.Printf("DEBUG: AptosTestAddLaneChangeSet: Configured offramp\n")

	return deployment.ChangesetOutput{}, nil
}

type Aptos2AnyMessage struct {
	Receiver      []byte
	Data          []byte
	TokenAmounts  []Aptos2AnyTokenAmount
	FeeToken      aptos.AccountAddress
	FeeTokenStore aptos.AccountAddress
	ExtraArgs     []byte
}

type Aptos2AnyTokenAmount struct {
	Token      aptos.AccountAddress
	Amount     uint64
	TokenStore aptos.AccountAddress
}

func SendRequestAptos(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	cfg *CCIPSendReqConfig,
) (*onramp.OnRampCCIPMessageSent, error) { // TODO: chain independent return vailue
	//sourceSelector := cfg.SourceChain
	//destSelector := cfg.DestChain
	//msg := cfg.Message.(Aptos2AnyMessage)
	return nil, errors.New("TODO(aptos): SendRequestAptos")
}

func ConfirmCommitWithExpectedSeqNumRangeAptos(
	t *testing.T,
	srcSelector uint64,
	dest deployment.AptosChain,
	ccipChainState changeset.AptosCCIPChainState,
	startBlock *uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (any, error) {
	return nil, errors.New("TODO(aptos): ConfirmCommitWithExpectedSeqNumRangeAptos")
}

func waitForTx(t *testing.T, client aptos.AptosRpcClient, txHash string, duration time.Duration) {
	stopTime := time.Now().Add(duration)
	for time.Now().Before(stopTime) {
		time.Sleep(time.Second * 1)
		txInfo, err := client.TransactionByHash(txHash)
		if err == nil && txInfo.Type != api.TransactionVariantPending {
			return
		}
	}
	t.Fatalf("Failed to wait for transaction %s", txHash)
}
