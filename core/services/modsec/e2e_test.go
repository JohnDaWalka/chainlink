package modsec_test

import (
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/params"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/weth9"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	v2toml "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	configv2 "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecstorage"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var (
	// Always 1337 for the simulated chain
	simChainID = big.NewInt(1337)
	// prefundAmount is the amount of Ether to pre-fund the deployer account with.
	prefundAmountEth = big.NewInt(1_000_000)
	// prefundAmountWei is the prefund amount in wei.
	prefundAmountWei = prefundAmountEth.Mul(prefundAmountEth, big.NewInt(params.Ether))
)

// evmTestChainSelectors returns the selectors for the test EVM chains. We arbitrarily
// start this from the EVM test selector TEST_90000001 and limit the number of chains you can load
// to 10. This avoid conflicts with other selectors.
var evmTestChainSelectors = []uint64{
	chain_selectors.TEST_90000001.Selector,
	chain_selectors.TEST_90000002.Selector,
	chain_selectors.TEST_90000003.Selector,
	chain_selectors.TEST_90000004.Selector,
	chain_selectors.TEST_90000005.Selector,
	chain_selectors.TEST_90000006.Selector,
	chain_selectors.TEST_90000007.Selector,
	chain_selectors.TEST_90000008.Selector,
	chain_selectors.TEST_90000009.Selector,
	chain_selectors.TEST_90000010.Selector,
}

// startAutoMine triggers the simulated backend to create a new block at intervals defined by
// `blockTime`. After the test is done, it stops the mining goroutine.
func startAutoMine(t *testing.T, backend *simulated.Backend, blockTime time.Duration) {
	t.Helper()

	ctx := t.Context() // Available since Go 1.20
	ticker := time.NewTicker(blockTime)
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				backend.Commit()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// GenerateChainsEVM generates a number of simulated EVM chains for testing purposes.
func generateChainsEVM(t *testing.T, numChains int, numUsers int, blockTime time.Duration) []EVMChain {
	if numChains > len(evmTestChainSelectors) {
		require.Failf(t, "not enough test EVM chain selectors available", "max is %d",
			len(evmTestChainSelectors),
		)
	}

	chains := make([]EVMChain, 0, numChains)
	for i := range numChains {
		selector := evmTestChainSelectors[i]
		evmChainID, err := chain_selectors.ChainIdFromSelector(selector)
		if err != nil {
			t.Fatal(err)
		}

		// Generate a deployer account
		adminKey, err := crypto.GenerateKey()
		require.NoError(t, err, "failed to generate deployer key")

		adminTransactor, err := bind.NewKeyedTransactorWithChainID(adminKey, simChainID)
		require.NoError(t, err)

		// Prefund the admin account
		genesis := types.GenesisAlloc{
			adminTransactor.From: {Balance: prefundAmountWei},
		}

		additionalTransactors := make([]*bind.TransactOpts, 0, numUsers)
		for i := 0; i < numUsers; i++ {
			userKey, err := crypto.GenerateKey()
			require.NoError(t, err, "failed to generate user key")
			userTransactor, err := bind.NewKeyedTransactorWithChainID(userKey, simChainID)
			require.NoError(t, err)
			additionalTransactors = append(additionalTransactors, userTransactor)

			genesis[userTransactor.From] = types.Account{Balance: prefundAmountWei}
		}

		// Initialize the simulated backend with the genesis state
		backend := simulated.NewBackend(genesis, simulated.WithBlockGasLimit(50000000))
		backend.Commit() // Commit the genesis block

		// Start mining blocks if a block time is configured
		if blockTime > 0 {
			startAutoMine(t, backend, blockTime)
		}

		chains = append(chains, EVMChain{
			EVMChainID:  new(big.Int).SetUint64(evmChainID),
			Backend:     backend,
			DeployerKey: adminTransactor,
			Users:       additionalTransactors,
		})
	}

	return chains
}

type EVMChain struct {
	EVMChainID  *big.Int
	Backend     *simulated.Backend
	DeployerKey *bind.TransactOpts
	Users       []*bind.TransactOpts
}

func setupApp(t *testing.T) (chainlink.Application, ccipDeployments) {
	evmChains := generateChainsEVM(t, 2, 1, 500*time.Millisecond)
	source, dest := evmChains[0], evmChains[1]
	deployments := deployContracts(t, source, dest)

	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.Feature.LogPoller = ptr(true)

		c.Log.Level = ptr(configv2.LogLevel(zapcore.DebugLevel))

		var evmConfigs v2toml.EVMConfigs
		for _, chain := range evmChains {
			evmConfigs = append(evmConfigs, createConfigV2Chain(chain.EVMChainID.Uint64()))
		}
		c.EVM = evmConfigs
	})

	// Create clients for the core node backed by sim.
	clients := make(map[uint64]client.Client)
	for _, chain := range evmChains {
		if chain.Backend != nil {
			clients[chain.EVMChainID.Uint64()] = client.NewSimulatedBackendClient(t, chain.Backend, chain.EVMChainID)
		}
	}

	// Set logging.
	lggr := logger.TestLogger(t)

	master := keystore.New(db, utils.FastScryptParams, lggr)
	ctx := t.Context()
	require.NoError(t, master.Unlock(ctx, "password"))
	require.NoError(t, master.CSA().EnsureKey(ctx))
	require.NoError(t, master.Workflow().EnsureKey(ctx))

	// OCR signing key for evm chains
	require.NoError(t, master.OCR2().EnsureKeys(ctx, chaintype.EVM))

	// transmitter key for each chain
	for _, chain := range evmChains {
		require.NoError(t, master.Eth().EnsureKeys(ctx, chain.EVMChainID))
	}

	app, err := chainlink.NewApplication(ctx, chainlink.ApplicationOpts{
		Config:   cfg,
		DS:       db,
		KeyStore: master,
		// TODO BCF-2513 Stop injecting ethClient via override, instead use httptest.
		EVMFactoryConfigFn: func(fc *chainlink.EVMFactoryConfig) {
			// Create ChainStores that always sign with 1337
			fc.GenChainStore = func(ks core.Keystore, i *big.Int) keys.ChainStore {
				return keys.NewChainStore(ks, big.NewInt(1337))
			}
			fc.GenEthClient = func(i *big.Int) client.Client {
				ethClient, ok := clients[i.Uint64()]
				if !ok {
					return client.NewNullClient(i, lggr)
				}
				return ethClient
			}
		},
		Logger:                   lggr,
		ExternalInitiatorManager: nil,
		CloseLogger:              lggr.Sync,
		UnrestrictedHTTPClient:   &http.Client{},
		RestrictedHTTPClient:     &http.Client{},
		AuditLogger:              audit.NoopLogger,
		RetirementReportCache:    retirement.NewRetirementReportCache(lggr, db),
	})
	require.NoError(t, err)

	return app, deployments
}

func createConfigV2Chain(chainID uint64) *v2toml.EVMConfig {
	chainIDBig := evmutils.NewI(int64(chainID))
	chain := v2toml.Defaults(chainIDBig)
	chain.GasEstimator.LimitDefault = ptr(uint64(5e6))
	chain.LogPollInterval = config.MustNewDuration(500 * time.Millisecond)
	chain.Transactions.ForwardersEnabled = ptr(false)
	chain.FinalityDepth = ptr(uint32(2))
	return &v2toml.EVMConfig{
		ChainID: chainIDBig,
		Enabled: ptr(true),
		Chain:   chain,
		Nodes:   v2toml.EVMNodes{&v2toml.Node{}},
	}
}

type ccipChainDeployment struct {
	evmChainID   *big.Int
	routerAddr   common.Address
	mockRMNAddr  common.Address
	armProxyAddr common.Address
	wethAddr     common.Address
}

type ccipDeployments struct {
	source ccipChainDeployment
	dest   ccipChainDeployment
}

func deploySingleChain(t *testing.T, chain EVMChain) ccipChainDeployment {
	// deploy wrapped native
	wethAddr, _, _, err := weth9.DeployWETH9(chain.DeployerKey, chain.Backend.Client())
	require.NoError(t, err)
	chain.Backend.Commit()

	// deploy mock arm
	mockRMNAddr, _, _, err := mock_rmn_contract.DeployMockRMNContract(chain.DeployerKey, chain.Backend.Client())
	require.NoError(t, err)
	chain.Backend.Commit()

	// deploy arm proxy
	armProxyAddr, _, _, err := rmn_proxy_contract.DeployRMNProxy(chain.DeployerKey, chain.Backend.Client(), mockRMNAddr)
	require.NoError(t, err)
	chain.Backend.Commit()

	// deploy router
	routerAddr, _, _, err := router.DeployRouter(chain.DeployerKey, chain.Backend.Client(), wethAddr, armProxyAddr)
	require.NoError(t, err)
	chain.Backend.Commit()

	return ccipChainDeployment{
		evmChainID:   chain.EVMChainID,
		routerAddr:   routerAddr,
		mockRMNAddr:  mockRMNAddr,
		armProxyAddr: armProxyAddr,
		wethAddr:     wethAddr,
	}
}

func deployContracts(t *testing.T, source, dest EVMChain) ccipDeployments {
	sourceChain := deploySingleChain(t, source)
	destChain := deploySingleChain(t, dest)

	return ccipDeployments{
		source: sourceChain,
		dest:   destChain,
	}
}

func sendMessage(t *testing.T, sourceDeployment, destDeployment ccipChainDeployment) {
	// TODO: implement
}

func TestModsec_E2E(t *testing.T) {
	app, deployments := setupApp(t)

	require.NoError(t, app.Start(t.Context()))
	t.Cleanup(func() {
		require.NoError(t, app.Stop())
	})

	storageServer, cleanup := modsecstorage.NewTestServer()
	t.Cleanup(cleanup)

	jb, err := modsec.ValidatedModsecSpec(testspecs.GenerateModsecSpec(
		testspecs.ModsecSpecParams{
			Name:                    "modsec-test-e2e",
			SourceChainID:           deployments.source.evmChainID.String(),
			DestinationChainID:      deployments.dest.evmChainID.String(),
			SourceChainFamily:       string(chaintype.EVM),
			DestinationChainFamily:  string(chaintype.EVM),
			OnRampAddress:           deployments.source.routerAddr.Hex(),
			OffRampAddress:          deployments.dest.routerAddr.Hex(),
			CCIPMessageSentEventSig: common.HexToHash("0x1").String(),
			StorageEndpoint:         storageServer.URL,
			StorageType:             "std",
		},
	).Toml())
	require.NoError(t, err)

	require.NoError(t, app.AddJobV2(t.Context(), &jb))
}

func ptr[T any](v T) *T { return &v }
