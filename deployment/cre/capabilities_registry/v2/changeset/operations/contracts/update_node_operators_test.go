package contracts_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestUpdateNopsOpWithMCMS(t *testing.T) {
	t.Parallel()
	doUpdateNopsOp(t, true)
}

func TestUpdateNopsOpWithoutMCMS(t *testing.T) {
	t.Parallel()
	doUpdateNopsOp(t, false)
}

func doUpdateNopsOp(t *testing.T, useMcms bool) {
	env := test.SetupEnvV2(t, false)
	b := optest.NewBundle(t)

	// Prepare input map of node operator params
	nopParams := make(map[uint32]capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams)
	nopParams[1] = capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams{
		Admin: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		Name:  "nop-1-updated",
	}
	nopParams[2] = capabilities_registry_v2.CapabilitiesRegistryNodeOperatorParams{
		Admin: common.HexToAddress("0x0000000000000000000000000000000000000002"),
		Name:  "nop-2-updated",
	}

	deps := contracts.UpdateNopsDeps{
		Env:           env.Env,
		MCMSContracts: te.MCMSWithTimelockState,
	}
	input := contracts.UpdateNopsInput{
		Address:       te.CapabilityRegistryAddressRef().Address.Hex(),
		ChainSelector: te.RegistrySelector,
		Nops:          nopParams,
	}
	if useMcms {
		input.MCMSConfig = &ocr3.MCMSConfig{MinDuration: 0}
	}

	opOutput, err := operations.ExecuteOperation(b, contracts.UpdateNops, deps, input)
	require.NoError(t, err)
	if useMcms {
		require.NotEmpty(t, opOutput.Output.Proposals)
	} else {
		require.Empty(t, opOutput.Output.Proposals)
	}
}
