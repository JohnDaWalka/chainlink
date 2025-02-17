package don

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func FundNodes(keystoneEnv *keystonetypes.KeystoneEnvironment) error {
	if keystoneEnv == nil {
		return errors.New("keystone environment must not be nil")
	}

	for _, don := range keystoneEnv.MustDons() {
		for _, node := range don.Nodes {
			_, err := libfunding.SendFunds(zerolog.Logger{}, keystoneEnv.MustSethClient(), libtypes.FundsToSend{
				ToAddress:  common.HexToAddress(node.AccountAddr[keystoneEnv.MustSethClient().Cfg.Network.ChainID]),
				Amount:     big.NewInt(5000000000000000000),
				PrivateKey: keystoneEnv.MustSethClient().MustGetRootPrivateKey(),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to send funds to node %s", node.AccountAddr[keystoneEnv.MustSethClient().Cfg.Network.ChainID])
			}
		}
	}

	return nil
}
