package solana

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var (
	// DeployUSDCTokenPoolContractsChangeset is a changeset that deploys the USDC token pool contracts.
	DeployUSDCTokenPoolContractsChangeset = cldf.CreateChangeSet(deployUSDCTokenPoolLogic, deployUSDCTokenPoolPrecondition)
)

type DeployUSDCTokenPoolContractsConfig struct {
}

func deployUSDCTokenPoolLogic(e cldf.Environment, c DeployUSDCTokenPoolContractsConfig) (cldf.ChangesetOutput, error) {
	return cldf.ChangesetOutput{}, nil
}

func deployUSDCTokenPoolPrecondition(e cldf.Environment, c DeployUSDCTokenPoolContractsConfig) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	fmt.Println(state)
	return nil
}
