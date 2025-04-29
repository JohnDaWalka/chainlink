package v0_5_0

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/stretchr/testify/require"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestCallSetProductionConfig(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true, 0)

	chainSelector := e.AllChainSelectors()[0]

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployConfiguratorChangeset,
			DeployConfiguratorConfig{
				ChainsToDeploy: []uint64{chainSelector},
			},
		),
	)

	require.NoError(t, err)

	configuratorAddrHex, err := dsutil.GetContractAddress(e.DataStore.Addresses(), types.Configurator)
	require.NoError(t, err)

	configuratorAddr := common.HexToAddress(configuratorAddrHex)

	onchainConfigHex := "0000000000000000000000000000000000000000000000000000000000000001" +
		"0000000000000000000000000000000000000000000000000000000000000000"
	onchainConfig, err := hex.DecodeString(onchainConfigHex)
	require.NoError(t, err)
	require.Len(t, onchainConfig, 64)

	prodCfg := SetProductionConfig{
		ConfiguratorAddress: configuratorAddr,
		ConfigID:            [32]byte{},
		Signers: [][]byte{
			{0x01}, {0x02}, {0x03}, {0x04},
		},
		OffchainTransmitters: [][32]byte{
			{}, {}, {}, {},
		},
		F:                     1,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain config"),
	}

	callConf := SetProductionConfigConfig{
		ConfigurationsByChain: map[uint64][]SetProductionConfig{
			chainSelector: {prodCfg},
		},
		MCMSConfig: nil,
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetProductionConfigChangeset,
			callConf,
		),
	)

	require.NoError(t, err)
}
