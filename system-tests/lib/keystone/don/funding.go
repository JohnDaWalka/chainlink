package don

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func FundNodes(t *testing.T, keystoneEnv *keystonetypes.KeystoneEnvironment) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.SethClient, "seth client must be set")
	require.NotNil(t, keystoneEnv.Dons, "dons must be set")

	for _, don := range keystoneEnv.Dons {
		for _, node := range don.Nodes {
			_, err := libfunding.SendFunds(zerolog.Logger{}, keystoneEnv.SethClient, libtypes.FundsToSend{
				ToAddress:  common.HexToAddress(node.AccountAddr[keystoneEnv.SethClient.Cfg.Network.ChainID]),
				Amount:     big.NewInt(5000000000000000000),
				PrivateKey: keystoneEnv.SethClient.MustGetRootPrivateKey(),
			})
			require.NoError(t, err, "failed to send funds to node %s", node.AccountAddr[keystoneEnv.SethClient.Cfg.Network.ChainID])
		}
	}
}
