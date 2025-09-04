package solana_test

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	solToken "github.com/gagliardetto/solana-go/programs/token"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_common"
	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/fee_quoter"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/test_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/timelock"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

func TestTestnetBillingTokenAccounts(t *testing.T) {
	for _, token := range []string{"GAn1RmBY76BjdqPAFnqM4jhqSSeeo9Zf7ASqHoWLdwZb", solana.WrappedSol.String()} {
		tokenPubKey := solana.MustPublicKeyFromBase58(token)
		feeQuoterAddress := solana.MustPublicKeyFromBase58("FeeQPGkKDeRV1MgoYfMH6L8o3KeuYjwUZrgn4LRKfjHi")
		tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, feeQuoterAddress)
		t.Logf("tokenBillingPDA: %s", tokenBillingPDA.String())
		client := solRpc.New(solRpc.DevNet_RPC)
		var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
		err := solCommonUtil.GetAccountDataBorshInto(context.Background(), client, tokenBillingPDA, cldf_solana.SolDefaultCommitment, &token0ConfigAccount)
		if err != nil {
			t.Fatalf("error getting account data: %s", err)
		}
		t.Logf("token0ConfigAccount: %+v", token0ConfigAccount)
	}
}

func TestStagingBillingTokenAccounts(t *testing.T) {
	tokenPubKey := solana.MustPublicKeyFromBase58("D3HCrigxfvScYyokPC1YGpNgqyheVMVwbgP7XPywvEdc")
	feeQuoterAddress := solana.MustPublicKeyFromBase58("FeeQhewH1cd6ZyHqhfMiKAQntgzPT6bWwK26cJ5qSFo6")
	tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, feeQuoterAddress)
	t.Logf("tokenBillingPDA: %s", tokenBillingPDA.String())
	client := solRpc.New(solRpc.DevNet_RPC)
	var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
	err := solCommonUtil.GetAccountDataBorshInto(context.Background(), client, tokenBillingPDA, cldf_solana.SolDefaultCommitment, &token0ConfigAccount)
	if err != nil {
		t.Fatalf("error getting account data: %s", err)
	}
	t.Logf("token0ConfigAccount: %+v", token0ConfigAccount)
}

func TestMainnetBillingTokenAccounts(t *testing.T) {
	for _, token := range []string{"AF94N2kSQEq6t6tBtPSLAGeKtn4QSJvayVyKL8m2DTfS", solana.WrappedSol.String()} {
		tokenPubKey := solana.MustPublicKeyFromBase58(token)
		feeQuoterAddress := solana.MustPublicKeyFromBase58("FeeQPGkKDeRV1MgoYfMH6L8o3KeuYjwUZrgn4LRKfjHi")
		tokenBillingPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenPubKey, feeQuoterAddress)
		t.Logf("tokenBillingPDA: %s", tokenBillingPDA.String())
		client := solRpc.New(solRpc.MainNetBeta_RPC)
		var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
		err := solCommonUtil.GetAccountDataBorshInto(context.Background(), client, tokenBillingPDA, cldf_solana.SolDefaultCommitment, &token0ConfigAccount)
		if err != nil {
			t.Fatalf("error getting account data: %s", err)
		}
		t.Logf("token0ConfigAccount: %+v", token0ConfigAccount)
	}
}

func TestMainnetTimelock(t *testing.T) {
	programID, seed, err := state.DecodeAddressWithSeed("DoajfR5tK24xVw51fWcawUZWhAXD8yrBJVacc13neVQA.BG7wilBWT4mc6p9yFnmfcu3yX7r9dazl")
	require.NoError(t, err)
	client := solRpc.New(solRpc.MainNetBeta_RPC)
	timelockConfigPDA := state.GetTimelockConfigPDA(programID, seed)
	var timelockData timelock.Config
	err = solCommonUtil.GetAccountDataBorshInto(context.Background(), client, timelockConfigPDA, cldf_solana.SolDefaultCommitment, &timelockData)
	if err != nil {
		t.Fatalf("error getting account data: %s", err)
	}
	t.Log("timelockData: ", timelockData)
}

func TestMainnetDestChainFeeQuoter(t *testing.T) {
	var destChainFqAccount solFeeQuoter.DestChain
	feeQuoterAddress := solana.MustPublicKeyFromBase58("FeeQPGkKDeRV1MgoYfMH6L8o3KeuYjwUZrgn4LRKfjHi")
	client := solRpc.New(solRpc.MainNetBeta_RPC)
	for _, destChain := range []uint64{11344663589394136015, 5009297550715157269, 4949039107694359620} {
		fqEvmDestChainPDA, _, _ := solState.FindFqDestChainPDA(destChain, feeQuoterAddress)
		err := solCommonUtil.GetAccountDataBorshInto(context.Background(), client, fqEvmDestChainPDA, cldf_solana.SolDefaultCommitment, &destChainFqAccount)
		if err != nil {
			t.Fatalf("failed to get account info: %s", err)
		}
		require.NoError(t, err, "failed to get account info")
		t.Logf("destChainFqAccount: %+v", destChainFqAccount)
		t.Logf("destChainFqAccount.Config: %+v", destChainFqAccount.Config)
		t.Logf("destChainFqAccount.State: %+v", destChainFqAccount.State)
	}
}

func TestMainnetTokenPools(t *testing.T) {
	maplePool := solana.MustPublicKeyFromBase58("787uwTCd8b2ikQP6g9AapMky36PWDv9x1XpC5ZUAfDYc")
	pepePool := solana.MustPublicKeyFromBase58("Br4t49vpqk7V9JTDjnea7eNmfjZXqoWpuRgX7NhUAEV")
	zeusPool := solana.MustPublicKeyFromBase58("BzJ846z8pxXrXf8xazvrSnMUbYYtpkzN5o8a5rH9exLx")
	solvPool := solana.MustPublicKeyFromBase58("ECvqYduigrFHeAU1kFCkehiiQz9eaeddUz6gH7BfD7AL")
	ohmPool := solana.MustPublicKeyFromBase58("EKyBTYcVaMGYocE7gTa1awNrnfxmzfa1LVjA6CfevsNd")
	uniPool := solana.MustPublicKeyFromBase58("2s5SB1UGcQusXNHUwMmscBaTanDBAWoBrEzNuzuyjfzh")

	SolvBTCToken := solana.MustPublicKeyFromBase58("SoLvHDFVstC74Jr9eNLTDoG4goSUsn1RENmjNtFKZvW")
	xSolvBTCToken := solana.MustPublicKeyFromBase58("SoLvAiHLF7LGEaiTN5KGZt1bNnraoWTi5mjcvRoDAX4")
	SolvBTCJUP := solana.MustPublicKeyFromBase58("SoLvzL3ZVjofmNB5LYFrf94QtNhMUSea4DawFhnAau8")
	syrupToken := solana.MustPublicKeyFromBase58("AvZZF1YaZDziPY2RCK4oJrRVrbN3mTD9NL24hPeaZeUj")
	ZBTCToken := solana.MustPublicKeyFromBase58("zBTCug3er3tLyffELcvDNrKkCymbPWysGcWihESYfLg")
	pepeToken := solana.MustPublicKeyFromBase58("zBTCug3er3tLyffELcvDNrKkCymbPWysGcWihESYfLg")
	ohmToken := solana.MustPublicKeyFromBase58("2Xva1NeLRuBFdK41gEuXqgeWtnKKDve9PKeCnMEpNG6K")
	uniToken := solana.MustPublicKeyFromBase58("uniBKsEV37qLRFZD7v3Z9drX6voyiCM8WcaePqeSSLc")
	client := solRpc.New(solRpc.MainNetBeta_RPC)

	for _, token := range []solana.PublicKey{uniToken, SolvBTCToken, xSolvBTCToken, SolvBTCJUP, syrupToken, ZBTCToken, ohmToken, pepeToken} {
		var tokenPool solana.PublicKey
		var poolConfigAccount solTestTokenPool.State
		switch token {
		case SolvBTCToken, xSolvBTCToken, SolvBTCJUP:
			tokenPool = solvPool
		case syrupToken:
			tokenPool = maplePool
		case ZBTCToken:
			tokenPool = zeusPool
		case pepeToken:
			tokenPool = pepePool
		case ohmToken:
			tokenPool = ohmPool
		case uniToken:
			tokenPool = uniPool
		default:
		}
		poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(token, tokenPool)
		require.NoError(t, err)
		err = solCommonUtil.GetAccountDataBorshInto(context.Background(), client, poolConfigPDA, cldf_solana.SolDefaultCommitment, &poolConfigAccount)
		require.NoError(t, err, "failed to get account info")
		t.Log("poolConfigPDA: ", poolConfigPDA)
		t.Log("poolConfigAccount: ", poolConfigAccount)
		t.Log("poolConfigAccount Owner: ", poolConfigAccount.Config.Owner)
		t.Log("token: \n", token)
	}
}

func TestMainnetRemoteAccounts(t *testing.T) {
	offRampID := solana.MustPublicKeyFromBase58("offqSMQWgQud6WJz694LRzkeN5kMYpCHTpXQr3Rkcjm")
	ccipRouterID := solana.MustPublicKeyFromBase58("Ccip842gzYHhvdDkSyi2YVCoAWPbYJoApMFzSxQroE9C")
	for _, remote := range []uint64{1673871237479749969} {
		client := solRpc.New(solRpc.MainNetBeta_RPC)
		offRampRemoteStatePDA, _, _ := solState.FindOfframpSourceChainPDA(remote, offRampID)
		var destChainStateAccount solOffRamp.SourceChain
		err := solCommonUtil.GetAccountDataBorshInto(context.Background(), client, offRampRemoteStatePDA, cldf_solana.SolDefaultCommitment, &destChainStateAccount)
		require.NoError(t, err, "failed to get account info")
		t.Logf("remote dest: %+v", offRampRemoteStatePDA)

		routerRemoteStatePDA, _ := solState.FindDestChainStatePDA(remote, ccipRouterID)
		var destChainStateAccount2 solRouter.DestChain
		err = solCommonUtil.GetAccountDataBorshInto(context.Background(), client, routerRemoteStatePDA, cldf_solana.SolDefaultCommitment, &destChainStateAccount2)
		require.NoError(t, err, "failed to get account info")
		t.Logf("remote source: %+v", routerRemoteStatePDA)
	}
}

func TestMainnet2TokePoolRemoteChainConfig(t *testing.T) {
	tokenAddress := solana.MustPublicKeyFromBase58("2Xva1NeLRuBFdK41gEuXqgeWtnKKDve9PKeCnMEpNG6K")
	burnMintTokenPool := solana.MustPublicKeyFromBase58("EKyBTYcVaMGYocE7gTa1awNrnfxmzfa1LVjA6CfevsNd")
	for _, remoteChainSelector := range []uint64{5009297550715157269} {
		remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(remoteChainSelector, tokenAddress, burnMintTokenPool)
		var remoteChainConfigAccount solTestTokenPool.ChainConfig
		client := solRpc.New(solRpc.MainNetBeta_RPC)
		err := solCommonUtil.GetAccountDataBorshInto(context.Background(), client, remoteChainConfigPDA, solRpc.CommitmentConfirmed, &remoteChainConfigAccount)
		require.NoError(t, err)
		// require.Equal(t, uint8(9), remoteChainConfigAccount.Base.Remote.Decimals)
		t.Logf("Pool addresses: %v", remoteChainConfigAccount.Base.Remote.PoolAddresses)
		require.Len(t, remoteChainConfigAccount.Base.Remote.PoolAddresses, 1)
		t.Logf("remoteChainConfigAccount: %+v", remoteChainConfigAccount)
	}
}

func TestMainnetTokenAdminRegistry(t *testing.T) {
	tokenAddress := solana.MustPublicKeyFromBase58("Dz9mQ9NzkBcCsuGPFJ3r1bS4wgqKMHBPiVuniW8Mbonk")
	routerAddress := solana.MustPublicKeyFromBase58("Ccip842gzYHhvdDkSyi2YVCoAWPbYJoApMFzSxQroE9C")
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenAddress, routerAddress)
	require.NoError(t, err, "failed to find token admin registry PDA")
	var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
	client := solRpc.New(solRpc.MainNetBeta_RPC)
	err = solCommonUtil.GetAccountDataBorshInto(context.Background(), client, tokenAdminRegistryPDA, solRpc.CommitmentConfirmed, &tokenAdminRegistryAccount)
	require.NoError(t, err, "failed to get token admin registry account data")
	t.Logf("tokenAdminRegistryPDA: %s", tokenAdminRegistryPDA.String())
	t.Logf("tokenAdminRegistryAccount: %+v", tokenAdminRegistryAccount)
}

func TestTokenData(t *testing.T) {
	tokenAddress := solana.MustPublicKeyFromBase58("LinkhB3afbBKb2EQQu7s7umdZceV3wcvAUJhQAfQ23L")
	var tokenMint solToken.Mint
	client := solRpc.New(solRpc.MainNetBeta_RPC)
	err := solCommonUtil.GetAccountDataBorshInto(context.Background(), client, tokenAddress, solRpc.CommitmentConfirmed, &tokenMint)
	require.NoError(t, err, "failed to get token mint account data")
	t.Logf("tokenMint: %+v", tokenMint)
}
