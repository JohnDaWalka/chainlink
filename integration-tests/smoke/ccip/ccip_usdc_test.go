package smoke

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

/*
* Chain topology for this test
* 	chainA (USDC, MY_TOKEN)
*			|
*			| ------- chainC (USDC, MY_TOKEN)
*			|
* 	chainB (USDC)
 */
func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	config := &changeset.TestConfigs{
		IsUSDC: true,
	}
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, config)
	//tenv := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, lggr, 3, 4, config)

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	require.Len(t, allChainSelectors, 3, "expected 3 chains for this test")
	chainA := allChainSelectors[0]
	chainC := allChainSelectors[1]
	chainB := allChainSelectors[2]

	aChainUSDC, cChainUSDC, err := changeset.ConfigureUSDCTokenPools(lggr, e.Chains, chainA, chainC, state)
	require.NoError(t, err)

	bChainUSDC, _, err := changeset.ConfigureUSDCTokenPools(lggr, e.Chains, chainB, chainC, state)
	require.NoError(t, err)

	aChainToken, _, cChainToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		chainA,
		chainC,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))

	mintAndAllow(t, e, state, map[uint64][]*burn_mint_erc677.BurnMintERC677{
		chainA: {aChainUSDC, aChainToken},
		chainB: {bChainUSDC},
		chainC: {cChainUSDC, cChainToken},
	})
	updateFeeQtrGrp := errgroup.Group{}
	updateFeeQtrGrp.Go(func() error {
		return changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[chainA], state.Chains[chainA], chainC, aChainUSDC)
	})

	updateFeeQtrGrp.Go(func() error {
		return changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[chainB], state.Chains[chainB], chainC, bChainUSDC)
	})
	updateFeeQtrGrp.Go(func() error {
		return changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[chainC], state.Chains[chainC], chainA, cChainUSDC)
	})
	require.NoError(t, updateFeeQtrGrp.Wait())

	// MockE2EUSDCTransmitter always mint 1, see MockE2EUSDCTransmitter.sol for more details
	tinyOneCoin := new(big.Int).SetUint64(1)
	tcs := []struct {
		id                     int
		name                   string
		receiver               common.Address
		sourceChain            uint64
		destChain              uint64
		tokens                 []router.ClientEVMTokenAmount
		data                   []byte
		expectedTokenBalances  map[common.Address]*big.Int
		expectedExecutionState int
		transferReturnData     transferReturnData
		initialBalances        map[common.Address]*big.Int
	}{
		{
			name:        "single USDC token transfer to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: chainC,
			destChain:   chainA,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				}},
			expectedTokenBalances: map[common.Address]*big.Int{
				aChainUSDC.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "multiple USDC tokens within the same message",
			receiver:    utils.RandomAddress(),
			sourceChain: chainC,
			destChain:   chainA,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				// 2 coins because of the same receiver
				aChainUSDC.Address(): new(big.Int).Add(tinyOneCoin, tinyOneCoin),
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "USDC token together with another token transferred to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: chainA,
			destChain:   chainC,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  aChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  aChainToken.Address(),
					Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				cChainUSDC.Address():  tinyOneCoin,
				cChainToken.Address(): new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "programmable token transfer to valid contract receiver",
			receiver:    state.Chains[chainC].Receiver.Address(),
			sourceChain: chainA,
			destChain:   chainC,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  aChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			data: []byte("hello world"),
			expectedTokenBalances: map[common.Address]*big.Int{
				cChainUSDC.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
	}

	for i, tt := range tcs {
		tcs[i].initialBalances = make(map[common.Address]*big.Int)
		for token := range tt.expectedTokenBalances {
			initialBalance := getTokenBalance(t, token, tt.receiver, e.Chains[tt.destChain])
			tcs[i].initialBalances[token] = initialBalance
		}

		// Send all requests first and update the return data
		tcs[i].transferReturnData = transfer(t, e, state, tt.sourceChain, tt.destChain, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(tt.receiver.Bytes(), 32),
			Data:         tt.data,
			TokenAmounts: tt.tokens,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})
	}

	for _, tt := range tcs {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			waitForSuccess(
				t,
				e,
				state,
				tt.sourceChain,
				tt.destChain,
				tt.transferReturnData,
				tt.expectedExecutionState,
			)

			for token, balance := range tt.expectedTokenBalances {
				t.Log("Checking token balance for token", token, "receiver", tt.receiver, "dest chain", tt.destChain)
				expected := new(big.Int).Add(tt.initialBalances[token], balance)
				waitForTheTokenBalance(t, token, tt.receiver, e.Chains[tt.destChain], expected)
			}
		})
	}

	t.Run("multi-source USDC transfer targeting the same dest receiver", func(t *testing.T) {
		sendSingleTokenTransfer := func(source, dest uint64, token common.Address, receiver common.Address) (*onramp.OnRampCCIPMessageSent, changeset.SourceDestPair) {
			msg := changeset.TestSendRequest(t, e, state, source, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
				Data:         []byte{},
				TokenAmounts: []router.ClientEVMTokenAmount{{Token: token, Amount: tinyOneCoin}},
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			return msg, changeset.SourceDestPair{
				SourceChainSelector: source,
				DestChainSelector:   dest,
			}
		}

		receiver := utils.RandomAddress()

		startBlocks := make(map[uint64]*uint64)
		expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
		expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)

		latesthdr, err := e.Chains[chainC].Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[chainC] = &block

		message1, message1ID := sendSingleTokenTransfer(chainA, chainC, aChainUSDC.Address(), receiver)
		expectedSeqNum[message1ID] = message1.SequenceNumber
		expectedSeqNumExec[message1ID] = []uint64{message1.SequenceNumber}

		message2, message2ID := sendSingleTokenTransfer(chainB, chainC, bChainUSDC.Address(), receiver)
		expectedSeqNum[message2ID] = message2.SequenceNumber
		expectedSeqNumExec[message2ID] = []uint64{message2.SequenceNumber}

		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)
		states := changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		require.Equal(t, changeset.EXECUTION_STATE_SUCCESS, states[message1ID][message1.SequenceNumber])
		require.Equal(t, changeset.EXECUTION_STATE_SUCCESS, states[message2ID][message2.SequenceNumber])

		// We sent 1 coin from each source chain, so we should have 2 coins on the destination chain
		// Receiver is randomly generated so we don't need to get the initial balance first
		expectedBalance := new(big.Int).Add(tinyOneCoin, tinyOneCoin)
		waitForTheTokenBalance(t, cChainUSDC.Address(), receiver, e.Chains[chainC], expectedBalance)
	})
}

// mintAndAllow mints tokens for deployers and allow router to spend them
func mintAndAllow(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	tkMap map[uint64][]*burn_mint_erc677.BurnMintERC677,
) {
	for chain, tokens := range tkMap {
		for _, token := range tokens {
			twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))

			tx, err := token.Mint(
				e.Chains[chain].DeployerKey,
				e.Chains[chain].DeployerKey.From,
				new(big.Int).Mul(twoCoins, big.NewInt(10)),
			)
			require.NoError(t, err)
			_, err = e.Chains[chain].Confirm(tx)
			require.NoError(t, err)

			tx, err = token.Approve(e.Chains[chain].DeployerKey, state.Chains[chain].Router.Address(), twoCoins)
			require.NoError(t, err)
			_, err = e.Chains[chain].Confirm(tx)
			require.NoError(t, err)
		}
	}
}

type transferReturnData struct {
	startBlocks        map[uint64]*uint64
	expectedSeqNum     map[changeset.SourceDestPair]uint64
	expectedSeqNumExec map[changeset.SourceDestPair][]uint64
}

func transfer(
	t *testing.T,
	env deployment.Environment,
	state changeset.CCIPOnChainState,
	sourceChain, destChain uint64,
	evm2AnyMessage router.ClientEVM2AnyMessage,
) transferReturnData {
	identifier := changeset.SourceDestPair{
		SourceChainSelector: sourceChain,
		DestChainSelector:   destChain,
	}

	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)

	latesthdr, err := env.Chains[destChain].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[destChain] = &block

	msgSentEvent := changeset.TestSendRequest(t, env, state, sourceChain, destChain, false, evm2AnyMessage)
	expectedSeqNum[identifier] = msgSentEvent.SequenceNumber
	expectedSeqNumExec[identifier] = []uint64{msgSentEvent.SequenceNumber}
	return transferReturnData{
		startBlocks:        startBlocks,
		expectedSeqNum:     expectedSeqNum,
		expectedSeqNumExec: expectedSeqNumExec,
	}
}

// transferAndWaitForSuccess sends a message from sourceChain to destChain and waits for it to be executed
func waitForSuccess(
	t *testing.T,
	env deployment.Environment,
	state changeset.CCIPOnChainState,
	sourceChain, destChain uint64,
	transferData transferReturnData,
	expectedStatus int,
) {
	identifier := changeset.SourceDestPair{
		SourceChainSelector: sourceChain,
		DestChainSelector:   destChain,
	}

	startBlocks := transferData.startBlocks
	expectedSeqNum := transferData.expectedSeqNum
	expectedSeqNumExec := transferData.expectedSeqNumExec

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, env, state, expectedSeqNum, startBlocks)

	// Wait for all exec reports to land
	states := changeset.ConfirmExecWithSeqNrsForAll(t, env, state, expectedSeqNumExec, startBlocks)
	require.Equal(t, expectedStatus, states[identifier][expectedSeqNum[identifier]])
}

func waitForTheTokenBalance(
	t *testing.T,
	token common.Address,
	receiver common.Address,
	chain deployment.Chain,
	expected *big.Int,
) {
	tokenContract, err := burn_mint_erc677.NewBurnMintERC677(token, chain.Client)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		actualBalance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, receiver)
		require.NoError(t, err)

		t.Log("Waiting for the token balance",
			"expected", expected,
			"actual", actualBalance,
			"token", token,
			"receiver", receiver,
		)

		return actualBalance.Cmp(expected) == 0
	}, tests.WaitTimeout(t), 100*time.Millisecond)
}

func getTokenBalance(
	t *testing.T,
	token common.Address,
	receiver common.Address,
	chain deployment.Chain,
) *big.Int {
	tokenContract, err := burn_mint_erc677.NewBurnMintERC677(token, chain.Client)
	require.NoError(t, err)

	balance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, receiver)
	require.NoError(t, err)

	t.Log("Getting token balance",
		"actual", balance,
		"token", token,
		"receiver", receiver,
	)

	return balance
}
