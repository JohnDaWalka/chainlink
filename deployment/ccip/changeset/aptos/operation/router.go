package operation

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

// UpdateRouterOp...
var UpdateRouterOp = operations.NewOperation(
	"update-router-op",
	Version1_0_0,
	"Updates Router destination chain configurations",
	updateRouterDests,
)

func updateRouterDests(b operations.Bundle, deps AptosDeps, in UpdateFeeQuoterDestsInput) ([]types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	routerBind := ccip_router.Bind(ccipAddress, deps.AptosChain.Client)

	// Process each destination chain config update
	// TODO: make this real
	var txs []types.Transaction

	moduleInfo, function, _, args, err := routerBind.Router().Encoder().SetOnRampVersions([]uint64{14767482510784806043, 16015286601757825753}, [][]byte{{1, 6, 0}, {1, 6, 0}})
	if err != nil {
		return []types.Transaction{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates for chains %d: %w", uint64(14767482510784806043), err)
	}

	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return []types.Transaction{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	txs = append(txs, types.Transaction{
		To:               ccipAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	})

	b.Logger.Infow("Adding Router destination config update operation")

	return txs, nil
}
