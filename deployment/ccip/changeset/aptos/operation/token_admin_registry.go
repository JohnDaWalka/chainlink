package operation

import (
	"encoding/json"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	ccipbind "github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

type TransferAdminRoleInput struct {
	Token    aptos.AccountAddress
	NewAdmin aptos.AccountAddress
}

// TransferAdminRoleOp
var TransferAdminRoleOp = operations.NewOperation(
	"transfer-admin-role",
	Version1_0_0,
	"Transfer admin role for pool on TokenAdminRegistry",
	transferAdminRole,
)

func transferAdminRole(b operations.Bundle, deps AptosDeps, in TransferAdminRoleInput) (types.Transaction, error) {
	ccipContract := ccipbind.Bind(deps.OnChainState.CCIPAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := ccipContract.TokenAdminRegistry().Encoder().TransferAdminRole(
		in.Token,
		in.NewAdmin,
	)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode TransferAdminRole for chain %d: %w", deps.AptosChain.Selector, err)
	}

	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	return types.Transaction{
		To:               deps.OnChainState.CCIPAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	}, nil
}

// AcceptAdminRoleOp
var AcceptAdminRoleOp = operations.NewOperation(
	"accept-admin-role",
	Version1_0_0,
	"Accept admin role for pool on TokenAdminRegistry",
	acceptAdminRole,
)

func acceptAdminRole(b operations.Bundle, deps AptosDeps, token aptos.AccountAddress) (types.Transaction, error) {
	ccipContract := ccipbind.Bind(deps.OnChainState.CCIPAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := ccipContract.TokenAdminRegistry().Encoder().AcceptAdminRole(token)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode AcceptAdminRole for chains %d: %w", deps.AptosChain.Selector, err)
	}

	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	return types.Transaction{
		To:               deps.OnChainState.CCIPAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	}, nil
}

type SetPoolInput struct {
	TokenAddress aptos.AccountAddress
	PoolAddress  aptos.AccountAddress
}

// SetPoolOp
var SetPoolOp = operations.NewOperation(
	"set-pool",
	Version1_0_0,
	"Set pool on TokenAdminRegistry",
	setPool,
)

func setPool(b operations.Bundle, deps AptosDeps, in SetPoolInput) (types.Transaction, error) {
	ccipContract := ccipbind.Bind(deps.OnChainState.CCIPAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := ccipContract.TokenAdminRegistry().Encoder().SetPool(
		in.TokenAddress,
		in.PoolAddress,
	)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to encode AcceptAdminRole for chains %d: %w", deps.AptosChain.Selector, err)
	}

	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	return types.Transaction{
		To:               deps.OnChainState.CCIPAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	}, nil
}
