package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestUpdateDON_Apply(t *testing.T) {
	// SetupEnvV2 deploys a cap reg v2 and configures it. So no need to do that here, just leverage the existing one.
	fixture := test.SetupEnvV2(t, false)

	input := changeset.UpdateDONInput{
		RegistryChainSel:  fixture.RegistrySelector,
		RegistryQualifier: test.RegistryQualifier,
		DonName:           test.DONName,
		IsPublic:          true,
		Force:             true,
	}

	// Preconditions
	err := changeset.UpdateDON{}.VerifyPreconditions(*fixture.Env, input)
	require.NoError(t, err)

	// Apply
	_, err = changeset.UpdateDON{}.Apply(*fixture.Env, input)
	require.NoError(t, err)

	// Validate on-chain state
	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		fixture.RegistryAddress,
		fixture.Env.BlockChains.EVMChains()[fixture.RegistrySelector].Client,
	)
	require.NoError(t, err)

	// DON capability configurations should include new capability config
	don, err := capReg.GetDONByName(nil, test.DONName)
	require.NoError(t, err)

	assert.Equal(t, test.DONName, don.Name)
	assert.True(t, don.IsPublic)
}
