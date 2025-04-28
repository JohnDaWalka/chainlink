package testhelpers

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	// "github.com/smartcontractkit/chainlink-ton/bindings/bind"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip_dummy_receiver"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip_offramp"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip_onramp"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip_router"
	// "github.com/smartcontractkit/chainlink-ton/bindings/mcms"
	// "github.com/smartcontractkit/chainlink-ton/relayer/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/onramp"
	tonaddress "github.com/xssnick/tonutils-go/address"
)

type TonTestDeployPrerequisitesChangeSet struct {
	T                 *testing.T
	TonChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestDeployPrerequisitesChangeSet{}

func (c TonTestDeployPrerequisitesChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {
	t := c.T

	onchainState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: TonTestDeployPrerequisitesChangeSet: chain selectors: %+v\n", c.TonChainSelectors)
	for _, chainSelector := range c.TonChainSelectors {
		tonChainState := onchainState.TonChains[chainSelector]
		// TODO: replace with a real token address instead of none
		tonChainState.LinkTokenAddress = *tonaddress.NewAddressNone()
		err = changeset.SaveOnchainStateTon(chainSelector, tonChainState, e)
		require.NoError(t, err)
	}
	return deployment.ChangesetOutput{}, nil
}

type TonTestDeployContractsChangeSet struct {
	T                 *testing.T
	HomeChainSelector uint64
	TonChainSelectors []uint64
	AllChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestDeployContractsChangeSet{}

func (c TonTestDeployContractsChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {

	t := c.T

	onchainState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: TonTestDeployContractsChangeSet: chain selectors: %+v\n", c.TonChainSelectors)

	for _, chainSelector := range c.TonChainSelectors {
		tonChain := e.TonChains[chainSelector]
		tonChainState := onchainState.TonChains[chainSelector]
		c.deployTonContracts(t, e, chainSelector, tonChain, tonChainState, onchainState)
	}
	return deployment.ChangesetOutput{}, nil
}

func (c TonTestDeployContractsChangeSet) deployTonContracts(t *testing.T, e deployment.Environment, chainSelector uint64, tonChain deployment.TonChain, tonChainState changeset.TonCCIPChainState, onchainState changeset.CCIPOnChainState) {
	logger := logger.Test(t)
	adminAddress := tonChain.DeployerSeed.AccountAddress()

	mcmsSeed := fmt.Sprintf("%d", time.Now().UnixNano())
	mcmsAddress, mcmsPendingTx, mcmsBindings, err := mcms.DeployToResourceAccount(tonChain.DeployerSigner, tonChain.Client, mcmsSeed)
	require.NoError(t, err)
	logger.Infow("Deployed Ton MCMS", "address", mcmsAddress.String(), "pendingTx", mcmsPendingTx.TxnHash())
	_ = mcmsBindings
	waitForTx(t, tonChain.Client, mcmsPendingTx.TxnHash(), time.Minute*1)

	ccipAddress, ccipPendingTx, ccipBindings, err := ccip.DeployToObject(tonChain.DeployerSigner, tonChain.Client, mcmsAddress, false)
	require.NoError(t, err)
	logger.Infow("Deployed Ton CCIP", "address", ccipAddress.String(), "pendingTx", ccipPendingTx.TxnHash())
	_ = ccipBindings
	waitForTx(t, tonChain.Client, ccipPendingTx.TxnHash(), time.Minute*1)

	offrampPendingTx, offrampBindings, err := ccip_offramp.DeployToExistingObject(tonChain.DeployerSigner, tonChain.Client, ccipAddress, ccipAddress, mcmsAddress)
	require.NoError(t, err)
	logger.Infow("Deployed Ton Offramp", "address", ccipAddress.String(), "pendingTx", ccipPendingTx.TxnHash())
	_ = offrampBindings
	waitForTx(t, tonChain.Client, offrampPendingTx.TxnHash(), time.Minute*1)

	onrampPendingTx, onrampBindings, err := ccip_onramp.DeployToExistingObject(tonChain.DeployerSigner, tonChain.Client, ccipAddress, ccipAddress, mcmsAddress)
	require.NoError(t, err)
	logger.Infow("Deployed Ton Onramp", "address", ccipAddress.String(), "pendingTx", ccipPendingTx.TxnHash())
	_ = onrampBindings
	waitForTx(t, tonChain.Client, onrampPendingTx.TxnHash(), time.Minute*1)

	tonChainState.CCIPAddress = ccipAddress

	ccipRouterPendingTx, ccipRouterBindings, err := ccip_router.DeployToExistingObject(tonChain.DeployerSigner, tonChain.Client, ccipAddress, ccipAddress, mcmsAddress)
	require.NoError(t, err)
	logger.Infow("Deployed Ton CCIP Router", "address", ccipAddress.String(), "pendingTx", ccipRouterPendingTx.TxnHash())
	_ = ccipRouterBindings
	waitForTx(t, tonChain.Client, ccipRouterPendingTx.TxnHash(), time.Minute*1)

	ccipDummyReceiverAddress, ccipDummyReceiverPendingTx, ccipDummyReceiverBindings, err := ccip_dummy_receiver.DeployToObject(tonChain.DeployerSigner, tonChain.Client, ccipAddress, mcmsAddress)
	require.NoError(t, err)
	logger.Infow("Deployed Ton CCIP Dummy Receiver", "address", ccipDummyReceiverAddress.String(), "pendingTx", ccipDummyReceiverPendingTx.TxnHash())
	_ = ccipDummyReceiverBindings
	waitForTx(t, tonChain.Client, ccipDummyReceiverPendingTx.TxnHash(), time.Minute*1)

	tonChainState.ReceiverAddress = ccipDummyReceiverAddress
	//tonChainState.ReceiverAddress = ton.AccountOne

	transactOpts := &bind.TransactOpts{
		Signer: tonChain.DeployerSigner,
	}

	pendingTx, err := onrampBindings.Onramp().Initialize(transactOpts, chainSelector, ton.AccountZero, adminAddress, []uint64{}, []tonaddress.Address{}, []bool{})
	require.NoError(t, err)
	logger.Infow("Initialized Onramp", "pendingTx", pendingTx.TxnHash())
	waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	// TODO: actually figure out where this value comes from
	//                                 failed to apply changeset at index 1: invalid changeset config: validate plugin info PluginType: CCIPCommit, Chains: [909606746561742123 5548718428018410741 4457093679053095497]: invalid ccip ocr params: invalid execute off-chain config: MessageVisibilityInterval=8h0m0s does not match the permissionlessExecutionThresholdSeconds in dynamic config =0 for chain 4457093679053095497
	permissionlessExecutionThresholdSecs := uint32(60 * 60 * 8)
	pendingTx, err = offrampBindings.Offramp().Initialize(transactOpts, chainSelector, permissionlessExecutionThresholdSecs, []uint64{}, []bool{}, []bool{}, [][]byte{})
	require.NoError(t, err)
	logger.Infow("Initialized Offramp", "pendingTx", pendingTx.TxnHash())
	waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	pendingTx, err = ccipBindings.FeeQuoter().Initialize(transactOpts, uint64(1000000), tonChainState.LinkTokenAddress, uint64(1000000), []tonaddress.Address{tonChainState.LinkTokenAddress})
	logger.Infow("Initialized FeeQuoter", "pendingTx", pendingTx.TxnHash())
	require.NoError(t, err)
	waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	pendingTx, err = ccipBindings.RMNRemote().Initialize(transactOpts, chainSelector)
	logger.Infow("Initialized RMNRemote", "pendingTx", pendingTx.TxnHash())
	require.NoError(t, err)
	waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	logger.Infow("All Ton contracts deployed")

	err = changeset.SaveOnchainStateTon(chainSelector, tonChainState, e)
	require.NoError(t, err)
}

type TonTestConfigureContractsChangeSet struct {
	T                 *testing.T
	HomeChainSelector uint64
	TonChainSelectors []uint64
	AllChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestConfigureContractsChangeSet{}

func (c TonTestConfigureContractsChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {

	t := c.T

	onchainState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: TonTestConfigureContractsChangeSet: chain selectors: %+v\n", c.TonChainSelectors)

	for _, chainSelector := range c.TonChainSelectors {
		tonChain := e.TonChains[chainSelector]
		tonChainState := onchainState.TonChains[chainSelector]
		c.configureTonContracts(t, e, chainSelector, tonChain, tonChainState, onchainState)
	}
	return deployment.ChangesetOutput{}, nil
}

func (c TonTestConfigureContractsChangeSet) configureTonContracts(t *testing.T, e deployment.Environment, chainSelector uint64, tonChain deployment.TonChain, tonChainState changeset.TonCCIPChainState, onchainState changeset.CCIPOnChainState) {
	logger := logger.Test(t)

	offrampBindings := ccip_offramp.Bind(tonChainState.CCIPAddress, tonChain.Client)

	transactOpts := &bind.TransactOpts{
		Signer: tonChain.DeployerSigner,
	}

	donID, err := internal.DonIDForChain(
		onchainState.Chains[c.HomeChainSelector].CapabilityRegistry,
		onchainState.Chains[c.HomeChainSelector].CCIPHome,
		chainSelector,
	)
	require.NoError(t, err)
	fmt.Printf("Ton DON ID: %+v\n", donID)

	allCommitConfigs, err := onchainState.Chains[c.HomeChainSelector].CCIPHome.GetAllConfigs(&ethbind.CallOpts{
		Context: context.Background(),
	}, donID, 0)

	allExecConfigs, err := onchainState.Chains[c.HomeChainSelector].CCIPHome.GetAllConfigs(&ethbind.CallOpts{
		Context: context.Background(),
	}, donID, 1)

	fmt.Printf("DEBUG HOME CHAIN CCIPHome: commit configs: %+v\n", allCommitConfigs)
	fmt.Printf("DEBUG HOME CHAIN CCIPHome: exec configs: %+v\n", allExecConfigs)

	ocr3Args, err := internal.BuildSetOCR3ConfigArgsTon(
		donID, onchainState.Chains[c.HomeChainSelector].CCIPHome, chainSelector, globals.ConfigTypeActive)
	require.NoError(t, err)

	var commitArgs *internal.MultiOCR3BaseOCRConfigArgsTon = nil
	var execArgs *internal.MultiOCR3BaseOCRConfigArgsTon = nil
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
		commitSigners = append(commitSigners, signer)
		fmt.Printf("DEBUG: configureTonContracts commit signer %x\n", signer)
	}
	commitTransmitters := []tonaddress.Address{}
	for _, transmitter := range commitArgs.Transmitters {
		address, err := utils.PublicKeyBytesToAddress(transmitter)
		require.NoError(t, err)
		commitTransmitters = append(commitTransmitters, address)
	}
	pendingTx, err := offrampBindings.Offramp().SetOcr3Config(transactOpts, commitArgs.ConfigDigest[:], uint8(types.PluginTypeCCIPCommit), commitArgs.F, commitArgs.IsSignatureVerificationEnabled, commitSigners, commitTransmitters)
	require.NoError(t, err)
	waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	execSigners := [][]byte{}
	for _, signer := range execArgs.Signers {
		execSigners = append(execSigners, signer)
	}
	execTransmitters := []tonaddress.Address{}
	for _, transmitter := range execArgs.Transmitters {
		address, err := utils.PublicKeyBytesToAddress(transmitter)
		require.NoError(t, err)
		execTransmitters = append(execTransmitters, address)
	}
	pendingTx, err = offrampBindings.Offramp().SetOcr3Config(transactOpts, execArgs.ConfigDigest[:], uint8(types.PluginTypeCCIPExec), execArgs.F, execArgs.IsSignatureVerificationEnabled, execSigners, execTransmitters)
	require.NoError(t, err)
	waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	logger.Infow("Ton contracts configured")

	for _, transmitter := range append(commitTransmitters, execTransmitters...) {
		// 10 APT
		entryFunction, err := ton.CoinTransferPayload(nil, transmitter, 1000000000)
		require.NoError(t, err)

		rawTxn, err := tonChain.Client.BuildTransaction(tonChain.DeployerSigner.AccountAddress(), ton.TransactionPayload{Payload: entryFunction})
		require.NoError(t, err)

		signedTxn, err := rawTxn.SignedTransaction(tonChain.DeployerSigner)
		require.NoError(t, err)

		submitResult, err := tonChain.Client.SubmitTransaction(signedTxn)
		require.NoError(t, err)

		waitForTx(t, tonChain.Client, submitResult.Hash, time.Minute*1)

		fmt.Printf("Sent 10 APT to transmitter %s\n", transmitter.String())
	}
}

func addLaneTonChangesets(t *testing.T, e *DeployedEnv, from, to uint64, fromFamily, toFamily string) []commoncs.ConfiguredChangeSet {
	fmt.Printf("DEBUG: addLaneTonChangesets %d to %d / %s to %s", from, to, fromFamily, toFamily)
	return []commoncs.ConfiguredChangeSet{TonTestAddLaneChangeSet{
		T:                 t,
		fromChainSelector: from,
		toChainSelector:   to,
		fromFamily:        fromFamily,
		toFamily:          toFamily,
	}}
}

type TonTestAddLaneChangeSet struct {
	T                 *testing.T
	fromChainSelector uint64
	toChainSelector   uint64
	fromFamily        string
	toFamily          string
}

var _ commoncs.ConfiguredChangeSet = TonTestAddLaneChangeSet{}

func (c TonTestAddLaneChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {
	t := c.T
	// TODO: support other paths
	require.Equal(t, c.fromFamily, chainsel.FamilyEVM, "must be from EVM")
	require.Equal(t, c.toFamily, chainsel.FamilyTon, "must be to Ton")

	tonSelector := c.toChainSelector
	tonChain := e.TonChains[tonSelector]

	onchainState, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	tonChainState := onchainState.TonChains[tonSelector]

	fmt.Printf("DEBUG: TonTestAddLaneChangeSet: LINK token: %s CCIP: %s Receiver: %s\n", tonChainState.LinkTokenAddress.String(), tonChainState.CCIPAddress.String(), tonChainState.ReceiverAddress.String())

	require.NotEqual(t, tonChainState.LinkTokenAddress, ton.AccountZero, "LINK token address must be set")
	require.NotEqual(t, tonChainState.CCIPAddress, ton.AccountZero, "CCIP address must be set")
	require.NotEqual(t, tonChainState.ReceiverAddress, ton.AccountZero, "Receiver address must be set")

	offrampBindings := ccip_offramp.Bind(tonChainState.CCIPAddress, tonChain.Client)
	transactOpts := &bind.TransactOpts{
		Signer: tonChain.DeployerSigner,
	}

	evmChainState := onchainState.Chains[c.fromChainSelector]
	evmOnrampAddress := evmChainState.OnRamp.Address().Bytes()
	fmt.Printf("DEBUG: TonTestAddLaneChangeSet: EVM chain selector: %d - EVM onramp address: %s\n", c.fromChainSelector, hex.EncodeToString(evmOnrampAddress))

	sourceChainSelectors := []uint64{c.fromChainSelector}
	sourceChainsIsEnabled := []bool{true}
	sourceChainsIsRMNVerificationDisabled := []bool{true}
	sourceChainsOnramps := [][]byte{evmOnrampAddress}
	pendingTx, err := offrampBindings.Offramp().ApplySourceChainConfigUpdates(transactOpts, sourceChainSelectors, sourceChainsIsEnabled, sourceChainsIsRMNVerificationDisabled, sourceChainsOnramps)
	require.NoError(t, err)
	waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	fmt.Printf("DEBUG: TonTestAddLaneChangeSet: Configured offramp\n")

	return deployment.ChangesetOutput{}, nil
}

type Ton2AnyMessage struct {
	Receiver      []byte
	Data          []byte
	TokenAmounts  []Ton2AnyTokenAmount
	FeeToken      tonaddress.Address
	FeeTokenStore tonaddress.Address
	ExtraArgs     []byte
}

type Ton2AnyTokenAmount struct {
	Token      tonaddress.Address
	Amount     uint64
	TokenStore tonaddress.Address
}

func SendRequestTon(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	cfg *CCIPSendReqConfig,
) (*onramp.OnRampCCIPMessageSent, error) { // TODO: chain independent return vailue
	//sourceSelector := cfg.SourceChain
	//destSelector := cfg.DestChain
	//msg := cfg.Message.(Ton2AnyMessage)
	return nil, errors.New("TODO(ton): SendRequestTon")
}

func ConfirmCommitWithExpectedSeqNumRangeTon(
	t *testing.T,
	srcSelector uint64,
	dest deployment.TonChain,
	ccipChainState changeset.TonCCIPChainState,
	startBlock *uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (any, error) {
	fmt.Printf("DEBUG: ConfirmCommitWithExpectedSeqNumRangeTon srcSelector: %d, startBlock: %+v, expectedSeqNumRange: %+v, enforceSingleCommit: %+v\n", srcSelector, startBlock, expectedSeqNumRange, enforceSingleCommit)

	time.Sleep(tests.WaitTimeout(t))
	return nil, errors.New("TODO(ton): ConfirmCommitWithExpectedSeqNumRangeTon")
}

func waitForTx(t *testing.T, client ton.TonRpcClient, txHash string, duration time.Duration) {
	userTx, err := client.WaitForTransaction(txHash, ton.PollTimeout(duration))
	require.NoError(t, err)
	require.True(t, userTx.Success, "transaction failed: %s", userTx.VmStatus)
}

func ConfirmExecWithSeqNrsTon(
	t *testing.T,
	sourceChain uint64,
	dest deployment.TonChain,
	offRampAddress tonaddress.Address,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	fmt.Printf("DEBUG: ConfirmExecWithSeqNrsTon srcSelector: %d, dest: %s, startBlock: %+v, expectedSeqNrs: %+v\n", sourceChain, startBlock, expectedSeqNrs)
	time.Sleep(tests.WaitTimeout(t))
	return nil, errors.New("TODO(ton): ConfirmExecWithSeqNrsTon")
	//fmt.Printf("DEBUG: TODO(ton): ConfirmExecWithSeqNrsTon\n")
	//seqNrsToWatch := make(map[uint64]int)
	//for _, seqNr := range expectedSeqNrs {
	//seqNrsToWatch[seqNr] = 0
	//}
	//return seqNrsToWatch, nil
}
