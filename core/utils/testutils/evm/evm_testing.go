package evm

import (
	"crypto/ecdsa"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

const network = "evm"

type TransactOptsAccounts struct {
	Deployer     *bind.TransactOpts
	Sender       *bind.TransactOpts
	TransactOpts []*bind.TransactOpts
}

type SimulatedBackend struct {
	Backend *simulated.Backend
	ChainID big.Int
	Network string
}

type TestRelayer struct {
	SimulatedBackend     SimulatedBackend
	Relayer              loop.Relayer
	TransactOptsAccounts TransactOptsAccounts
}

type TestAppWithSimulatedBackend struct {
	App              *cltest.TestApplication
	SimulatedBackend SimulatedBackend
	Accounts         TransactOptsAccounts
}

// TestAppConfig allows to customize the chainlink App configuration. It provides access to the simulated backend being used and also the to sender and deployer keys through transactionOptsAccounts.
type TestAppConfig struct {
	CustomizeConfigFunc func(config *chainlink.GeneralConfig,
		backend SimulatedBackend,
		transactOptsAccoutns TransactOptsAccounts)
}

// NewTestAppWithSimulatedBackend creates a new chainlink testing App with an EVM relayer configured backed up by a simulated backend. It also configures the EVM Relayer with a deplyer and sender key ready to be used. The creation of the App can be customized with the testAppConfig parameter.
func NewTestAppWithSimulatedBackend(t *testing.T, testAppConfig TestAppConfig) TestAppWithSimulatedBackend {
	accountKeys := setupKeys(t)
	transactOpts := createTransactOptsWithAccounts(t, accountKeys)
	chainID := testutils.SimulatedChainID
	backend := CreateSimulatedBackend(t, transactOpts)
	simulatedBackend := SimulatedBackend{
		Backend: backend,
		ChainID: *chainID,
		Network: network,
	}
	client := client.NewSimulatedBackendClient(t, backend, chainID)
	db := pgtest.NewSqlxDB(t)

	clconfig := configtest.NewGeneralConfigSimulated(t, nil)
	testAppConfig.CustomizeConfigFunc(&clconfig, simulatedBackend, transactOpts)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, clconfig, backend, db, client)
	signingKey := ethkey.FromPrivateKey(accountKeys.SenderKey)
	app.KeyStore.Eth().XXXTestingOnlyAdd(t.Context(), signingKey)
	err := app.KeyStore.Eth().Enable(t.Context(), signingKey.Address, chainID)
	require.NoError(t, err)
	return TestAppWithSimulatedBackend{
		App:              app,
		SimulatedBackend: simulatedBackend,
		Accounts:         transactOpts,
	}
}

func Client(t *testing.T, backend *simulated.Backend) client.Client {
	return client.NewSimulatedBackendClient(t, backend, big.NewInt(1337))
}

func createTransactOptsWithAccounts(t *testing.T, keys Keys) TransactOptsAccounts {
	deployer, err := bind.NewKeyedTransactorWithChainID(keys.DeployerKey, big.NewInt(1337))
	require.NoError(t, err)

	sender, err := bind.NewKeyedTransactorWithChainID(keys.SenderKey, big.NewInt(1337))
	require.NoError(t, err)

	return TransactOptsAccounts{
		Deployer:     deployer,
		Sender:       sender,
		TransactOpts: []*bind.TransactOpts{deployer, sender},
	}
}

type Keys struct {
	DeployerKey *ecdsa.PrivateKey
	SenderKey   *ecdsa.PrivateKey
}

func CreateSimulatedBackend(t *testing.T, transactOptsAccounts TransactOptsAccounts) *simulated.Backend {
	const commonGasLimitOnEvms = uint64(4712388)
	backend := simulated.NewBackend(
		evmtypes.GenesisAlloc{transactOptsAccounts.Deployer.From: {Balance: big.NewInt(math.MaxInt64)}, transactOptsAccounts.Sender.From: {Balance: big.NewInt(math.MaxInt64)}}, simulated.WithBlockGasLimit(commonGasLimitOnEvms*5000))
	cltest.Mine(backend, 10*time.Millisecond)
	return backend
}

func setupKeys(t *testing.T) Keys {
	deployerPkey, err := crypto.GenerateKey()
	require.NoError(t, err)

	senderPkey, err := crypto.GenerateKey()
	require.NoError(t, err)
	return Keys{
		DeployerKey: deployerPkey,
		SenderKey:   senderPkey,
	}
}
