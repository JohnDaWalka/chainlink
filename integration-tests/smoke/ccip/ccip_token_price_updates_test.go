package ccip

import (
	"context"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPTokenPriceUpdates(t *testing.T) {
	ctx := testhelpers.Context(t)
	callOpts := &bind.CallOpts{Context: ctx}

	var tokenPriceExpiry = 5 * time.Second
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			if params.CommitOffChainConfig != nil {
				params.CommitOffChainConfig.TokenPriceBatchWriteFrequency = *config.MustNewDuration(tokenPriceExpiry)
			}
			return params
		}))
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	testhelpers.AddLanesForAll(t, &e, state)

	allChainSelectors := maps.Keys(e.Env.Chains)
	assert.GreaterOrEqual(t, len(allChainSelectors), 2, "test requires at least 2 chains")

	sourceChain1 := allChainSelectors[0]

	feeQuoter1 := state.Chains[sourceChain1].FeeQuoter

	feeTokensChain1, err := feeQuoter1.GetFeeTokens(callOpts)
	require.NoError(t, err)
	t.Logf("feeTokens: %v", feeTokensChain1)

	tokenPricesBefore, err := feeQuoter1.GetTokenPrices(callOpts, feeTokensChain1)
	require.NoError(t, err)
	t.Logf("tokenPrices: %v", tokenPricesBefore)

	// assert token prices updated due to time expiration
	assert.Eventually(t, func() bool {
		tokenPricesNow, err := feeQuoter1.GetTokenPrices(callOpts, feeTokensChain1)
		require.NoError(t, err)
		t.Logf("tokenPrices: %v", tokenPricesNow)

		// both tokens should have same price but different timestamp since there was an update due to time deviation
		for i, price := range tokenPricesNow {
			if tokenPricesBefore[i].Timestamp == price.Timestamp {
				tokenPricesBefore = tokenPricesNow
				return false // timestamp is the same
			}
			if tokenPricesBefore[i].Value.Cmp(price.Value) != 0 {
				tokenPricesBefore = tokenPricesNow
				return false // price was updated
			}
		}
		t.Log("time expiration assertions complete")
		return true
	}, tests.WaitTimeout(t), 500*time.Millisecond)

	// disable oracles to prevent price updates while we manually edit token prices
	disabledOracleIDs := disableOracles(ctx, t, e.Env.Offchain)

	assert.Eventually(t, func() bool {
		// manually update token prices by setting values to maxUint64 and 0
		tx, err := feeQuoter1.UpdatePrices(e.Env.Chains[sourceChain1].DeployerKey, fee_quoter.InternalPriceUpdates{
			TokenPriceUpdates: []fee_quoter.InternalTokenPriceUpdate{
				{SourceToken: feeTokensChain1[0], UsdPerToken: big.NewInt(0).SetUint64(math.MaxUint64)},
				{SourceToken: feeTokensChain1[1], UsdPerToken: big.NewInt(0)},
			},
		})
		require.NoError(t, err)

		_, err = deployment.ConfirmIfNoError(e.Env.Chains[sourceChain1], tx, err)
		require.NoError(t, err)
		t.Logf("manually editing token prices")

		tokenPricesNow, err := feeQuoter1.GetTokenPrices(callOpts, feeTokensChain1)
		require.NoError(t, err)
		t.Logf("tokenPrices straight after: %v", tokenPricesNow)

		if uint64(math.MaxUint64) != tokenPricesNow[0].Value.Uint64() {
			return false
		}
		if uint64(0) != tokenPricesNow[1].Value.Uint64() {
			return false
		}
		return true

		// retry because there might've been a commit report inflight
	}, tests.WaitTimeout(t), 200*time.Millisecond)

	enableOracles(ctx, t, e.Env.Offchain, disabledOracleIDs)

	// wait until price goes back to the original
	assert.Eventually(t, func() bool {
		tokenPricesNow, err := feeQuoter1.GetTokenPrices(callOpts, feeTokensChain1)
		require.NoError(t, err)
		t.Logf("tokenPrices: %v tokenPricesBefore: %v", tokenPricesNow, tokenPricesBefore)

		if tokenPricesNow[0].Value.Cmp(tokenPricesBefore[0].Value) != 0 {
			return false
		}
		if tokenPricesNow[1].Value.Cmp(tokenPricesBefore[1].Value) != 0 {
			return false
		}
		return true
	}, tests.WaitTimeout(t), 500*time.Millisecond)
}

func Test_CCIPTokenPriceUpdatesMultiReport(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-5769")
	// Delta round is set to greater than token price batch write frequency
	// to ensure that token prices are written on every round.
	//
	// If we have two tokens and max prices per report are set to 1 then we should see two different timestamps
	// for the updates since two transactions / two reports must've been produced.
	//
	// After detecting the first price update the next one should not take long to appear since they must
	// both be coming from the same OCR round. If the second price update is not received within deltaRound time
	// it means that it was ignored by the contract and will need a few more rounds to be delivered.
	//
	// Test assumes that TX delivery time is less than deltaRound time. In an in-mem scenario with delta round set to
	// 45s that should be way more than enough time for TX delivery.

	ctx := testhelpers.Context(t)
	callOpts := &bind.CallOpts{Context: ctx}

	deltaRound := 45 * time.Second
	tokenPriceExpiration := time.Second

	// you can set this to a bigger value to verify that the test passes and only fails on multiple reports
	maxPricesPerReport := uint64(1)

	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			if params.CommitOffChainConfig != nil {
				params.CommitOffChainConfig.MaxPricesPerReport = maxPricesPerReport
				params.CommitOffChainConfig.MultipleReportsEnabled = true
				params.CommitOffChainConfig.TokenPriceBatchWriteFrequency = *config.MustNewDuration(tokenPriceExpiration)
				params.OCRParameters.DeltaRound = deltaRound
				params.CommitOffChainConfig.TransmissionDelayMultiplier = time.Second
			}
			return params
		}))

	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	testhelpers.AddLanesForAll(t, &e, state)

	allChainSelectors := maps.Keys(e.Env.Chains)
	assert.GreaterOrEqual(t, len(allChainSelectors), 2, "test requires at least 2 chains")

	sourceChain1 := allChainSelectors[0]
	feeQuoter1 := state.Chains[sourceChain1].FeeQuoter

	feeTokensChain1, err := feeQuoter1.GetFeeTokens(callOpts)
	require.NoError(t, err)
	t.Logf("feeTokens: %v", feeTokensChain1)

	tokenPricesBefore, err := feeQuoter1.GetTokenPrices(callOpts, feeTokensChain1)
	require.NoError(t, err)
	t.Logf("tokenPrices: %v", tokenPricesBefore)

	tokenPricesLastUpdatedAt := make(map[common.Address]uint32)
	for i, tk := range feeTokensChain1 {
		tokenPricesLastUpdatedAt[tk] = tokenPricesBefore[i].Timestamp
	}

	firstPriceUpdateDetectedAt := time.Time{}
	firstPriceUpdateToken := common.Address{}
	t.Logf("waiting for the first price update...")
	assert.Eventually(t, func() bool {
		if !firstPriceUpdateDetectedAt.IsZero() && time.Since(firstPriceUpdateDetectedAt) > deltaRound {
			t.Fatal("second price update not received, it required an extra ocr round")
		}

		tokenPricesNow, err := feeQuoter1.GetTokenPrices(callOpts, feeTokensChain1)
		require.NoError(t, err)

		for i, tk := range feeTokensChain1 {
			prevUpdatedAt, ok := tokenPricesLastUpdatedAt[tk]
			require.True(t, ok)
			if prevUpdatedAt != tokenPricesNow[i].Timestamp {
				tokenPricesLastUpdatedAt[tk] = tokenPricesNow[i].Timestamp

				t.Logf("token price update detected for token %s at %d", tk.Hex(), tokenPricesNow[i].Timestamp)

				if firstPriceUpdateDetectedAt.IsZero() {
					t.Logf("this is the first detected price update")
					firstPriceUpdateDetectedAt = time.Now()
					firstPriceUpdateToken = tk
					continue
				}

				t.Logf("this is the second detected price update, time since first update: %s", time.Since(firstPriceUpdateDetectedAt))
				if feeTokensChain1[i] == firstPriceUpdateToken {
					t.Fatal("second price update was for the same token as the first one, test failed")
				}
				t.Log("test passed")
				return true
			}
		}
		return false
	}, tests.WaitTimeout(t), 300*time.Millisecond)
}

func disableOracles(ctx context.Context, t *testing.T, client deployment.OffchainClient) []string {
	var disabledOracleIDs []string
	listNodesResp, err := client.ListNodes(ctx, &node.ListNodesRequest{})
	require.NoError(t, err)

	for _, n := range listNodesResp.Nodes {
		if strings.HasPrefix(n.Name, "bootstrap") {
			continue
		}
		_, err := client.DisableNode(ctx, &node.DisableNodeRequest{Id: n.Id})
		require.NoError(t, err)
		disabledOracleIDs = append(disabledOracleIDs, n.Id)
		t.Logf("node %s disabled", n.Id)
	}

	return disabledOracleIDs
}

func enableOracles(ctx context.Context, t *testing.T, client deployment.OffchainClient, oracleIDs []string) {
	for _, n := range oracleIDs {
		_, err := client.EnableNode(ctx, &node.EnableNodeRequest{Id: n})
		require.NoError(t, err)
		t.Logf("node %s enabled", n)
	}
}
