package don

import (
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs"
	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

const (
	GistIP = "185.199.108.133"
)

func Configure(t *testing.T, testLogger zerolog.Logger, keystoneEnv *types.KeystoneEnvironment, donToJobSpecs types.DonsToJobSpecs, donToConfigOverrides types.DonsToConfigOverrides) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.Environment, "environment must be set")
	require.NotNil(t, keystoneEnv.Blockchain, "blockchain must be set")
	require.NotNil(t, keystoneEnv.WrappedNodeOutput, "wrapped node output must be set")
	require.NotNil(t, keystoneEnv.JD, "job distributor must be set")
	require.NotNil(t, keystoneEnv.SethClient, "seth client must be set")
	require.NotEmpty(t, keystoneEnv.DONTopology, "DON topology must not be empty")
	require.NotNil(t, keystoneEnv.KeystoneContractAddresses, "keystone contract addresses must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.CapabilitiesRegistryAddress, "capabilities registry address must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.OCR3CapabilityAddress, "OCR3 capability address must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.ForwarderAddress, "forwarder address must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.WorkflowRegistryAddress, "workflow registry address must be set")
	require.GreaterOrEqual(t, len(keystoneEnv.DONTopology), 1, "expected at least one DON topology")
	require.NotNil(t, keystoneEnv.GatewayConnectorData, "gateway connector data must be set")

	for i, donTopology := range keystoneEnv.DONTopology {
		if configOverrides, ok := donToConfigOverrides[donTopology.ID]; ok {
			for j, configOverride := range configOverrides {
				require.GreaterOrEqual(t, len(donTopology.NodeInput.NodeSpecs)-1, j, "config override index out of bounds")
				donTopology.NodeInput.NodeSpecs[j].Node.TestConfigOverrides = configOverride
			}
			keystoneEnv.DONTopology[i].NodeOutput = config.Set(t, donTopology.NodeInput, keystoneEnv.Blockchain)
		}
	}

	nodeOutputs := make([]*types.WrappedNodeOutput, 0, len(keystoneEnv.DONTopology))
	for i := range keystoneEnv.DONTopology {
		nodeOutputs = append(nodeOutputs, keystoneEnv.DONTopology[i].NodeOutput)
	}

	// after restarting the nodes, we need to reinitialize the JD clients otherwise
	// communication between JD and nodes will fail due to invalidated session cookie
	keystoneEnv.Environment = jobs.ReinitialiseJDClients(t, keystoneEnv.Environment, keystoneEnv.JD, nodeOutputs...)
	for _, donTopology := range keystoneEnv.DONTopology {
		if jobSpecs, ok := donToJobSpecs[donTopology.ID]; ok {
			jobs.Create(t, keystoneEnv.Environment.Offchain, donTopology.DON, donTopology.Flags, jobSpecs)
		} else {
			testLogger.Error().Msgf("No job specs found for DON %d", donTopology.ID)
			t.FailNow()
		}
	}
}

func BuildDONTopology(t *testing.T, keystoneEnv *types.KeystoneEnvironment) {
	require.NotNil(t, keystoneEnv, "keystone environment must not be nil")
	require.NotNil(t, keystoneEnv.NodeInput, "keystone environment must have node input")
	require.NotNil(t, keystoneEnv.Dons, "keystone environment must have DONs")
	require.NotNil(t, keystoneEnv.WrappedNodeOutput, "keystone environment must have node outputs")

	require.Equal(t, len(keystoneEnv.Dons), len(keystoneEnv.WrappedNodeOutput), "number of DONs and node outputs must match")
	keystoneEnv.DONTopology = make([]*types.DONTopology, len(keystoneEnv.Dons))

	// one DON to do everything
	if len(keystoneEnv.Dons) == 1 {
		flags, err := keystoneflags.NodeSetFlags(keystoneEnv.NodeInput[0])
		require.NoError(t, err, "failed to convert string flags to bitmap for nodeset %s", keystoneEnv.NodeInput[0].Name)

		keystoneEnv.DONTopology[0] = &types.DONTopology{
			DON:        keystoneEnv.Dons[0],
			NodeInput:  keystoneEnv.NodeInput[0],
			NodeOutput: keystoneEnv.WrappedNodeOutput[0],
			ID:         1,
			Flags:      flags,
		}
	} else {
		for i, don := range keystoneEnv.Dons {
			flags, err := keystoneflags.NodeSetFlags(keystoneEnv.NodeInput[i])
			require.NoError(t, err, "failed to convert string flags to bitmap for nodeset %s", keystoneEnv.NodeInput[i].Name)

			keystoneEnv.DONTopology[i] = &types.DONTopology{
				DON:        don,
				NodeInput:  keystoneEnv.NodeInput[i],
				NodeOutput: keystoneEnv.WrappedNodeOutput[i],
				ID:         libc.MustSafeUint32(i + 1),
				Flags:      flags,
			}
		}
	}

	keystoneEnv.WorkflowDONID = keystoneflags.MustOneDONTopologyWithFlag(t, keystoneEnv.DONTopology, types.WorkflowDON).ID
}

// In order to whitelist host IP in the gateway, we need to resolve the host.docker.internal to the host IP,
// and since CL image doesn't have dig or nslookup, we need to use curl.
func ResolveHostDockerInternaIP(testLogger zerolog.Logger, nsOutput *ns.Output) (string, error) {
	containerName := nsOutput.CLNodes[0].Node.ContainerName
	cmd := []string{"curl", "-v", "http://host.docker.internal"}
	output, err := framework.ExecContainer(containerName, cmd)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`.*Trying ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+).*`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		testLogger.Error().Msgf("failed to extract IP address from curl output:\n%s", output)
		return "", errors.New("failed to extract IP address from curl output")
	}

	testLogger.Info().Msgf("Resolved host.docker.internal to %s", matches[1])

	return matches[1], nil
}
