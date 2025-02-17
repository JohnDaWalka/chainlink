package don

import (
	"regexp"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs"
	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

func Configure(t *testing.T, testLogger zerolog.Logger, keystoneEnv *types.KeystoneEnvironment, donToJobSpecs types.DonsToJobSpecs, donToConfigOverrides types.DonsToConfigOverrides) error {
	if keystoneEnv == nil {
		return errors.New("keystone environment must not be nil")
	}
	if keystoneEnv.Environment == nil {
		return errors.New("environment must be set")
	}
	if keystoneEnv.Blockchain == nil {
		return errors.New("blockchain must be set")
	}
	if keystoneEnv.WrappedNodeOutput == nil {
		return errors.New("wrapped node output must be set")
	}
	if keystoneEnv.JD == nil {
		return errors.New("job distributor must be set")
	}
	if keystoneEnv.SethClient == nil {
		return errors.New("seth client must be set")
	}
	if len(keystoneEnv.DONTopology) == 0 {
		return errors.New("DON topology must not be empty")
	}
	if keystoneEnv.KeystoneContractAddresses == nil {
		return errors.New("keystone contract addresses must be set")
	}
	if keystoneEnv.KeystoneContractAddresses.CapabilitiesRegistryAddress == (common.Address{}) {
		return errors.New("capabilities registry address must be set")
	}
	if keystoneEnv.KeystoneContractAddresses.OCR3CapabilityAddress == (common.Address{}) {
		return errors.New("OCR3 capability address must be set")
	}
	if keystoneEnv.KeystoneContractAddresses.ForwarderAddress == (common.Address{}) {
		return errors.New("forwarder address must be set")
	}
	if keystoneEnv.KeystoneContractAddresses.WorkflowRegistryAddress == (common.Address{}) {
		return errors.New("workflow registry address must be set")
	}
	if len(keystoneEnv.DONTopology) == 0 {
		return errors.New("expected at least one DON topology")
	}
	if keystoneEnv.GatewayConnectorData == nil {
		return errors.New("gateway connector data must be set")
	}

	for i, donTopology := range keystoneEnv.DONTopology {
		if configOverrides, ok := donToConfigOverrides[donTopology.ID]; ok {
			for j, configOverride := range configOverrides {
				if len(donTopology.NodeInput.NodeSpecs)-1 < j {
					return errors.Errorf("config override index out of bounds: %d", j)
				}
				donTopology.NodeInput.NodeSpecs[j].Node.TestConfigOverrides = configOverride
			}
			var setErr error
			keystoneEnv.DONTopology[i].NodeOutput, setErr = config.Set(t, donTopology.NodeInput, keystoneEnv.Blockchain)
			if setErr != nil {
				return errors.Wrap(setErr, "failed to set node output")
			}
		}
	}

	nodeOutputs := make([]*types.WrappedNodeOutput, 0, len(keystoneEnv.DONTopology))
	for i := range keystoneEnv.DONTopology {
		nodeOutputs = append(nodeOutputs, keystoneEnv.DONTopology[i].NodeOutput)
	}

	// after restarting the nodes, we need to reinitialize the JD clients otherwise
	// communication between JD and nodes will fail due to invalidated session cookie
	var jdErr error
	keystoneEnv.Environment, jdErr = jobs.ReinitialiseJDClients(keystoneEnv.Environment, keystoneEnv.JD, nodeOutputs...)
	if jdErr != nil {
		return errors.Wrap(jdErr, "failed to reinitialize JD clients")
	}
	for _, donTopology := range keystoneEnv.DONTopology {
		if jobSpecs, ok := donToJobSpecs[donTopology.ID]; ok {
			createErr := jobs.Create(keystoneEnv.Environment.Offchain, donTopology.DON, donTopology.Flags, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", donTopology.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", donTopology.ID)
		}
	}

	return nil
}

func BuildDONTopology(keystoneEnv *types.KeystoneEnvironment) error {
	if keystoneEnv == nil {
		return errors.New("keystone environment must not be nil")
	}
	if keystoneEnv.NodeInput == nil {
		return errors.New("node input must be set")
	}
	if len(keystoneEnv.Dons) == 0 {
		return errors.New("Dons must be set")
	}
	if len(keystoneEnv.WrappedNodeOutput) == 0 {
		return errors.New("wrapped node output must be set")
	}
	if len(keystoneEnv.Dons) != len(keystoneEnv.WrappedNodeOutput) {
		return errors.New("number of DONs and node outputs must match")
	}

	keystoneEnv.DONTopology = make([]*types.DONTopology, len(keystoneEnv.Dons))

	// one DON to do everything
	if len(keystoneEnv.Dons) == 1 {
		flags, err := keystoneflags.NodeSetFlags(keystoneEnv.NodeInput[0])
		if err != nil {
			return errors.Wrapf(err, "failed to convert string flags to bitmap for nodeset %s", keystoneEnv.NodeInput[0].Name)
		}

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
			if err != nil {
				return errors.Wrapf(err, "failed to convert string flags to bitmap for nodeset %s", keystoneEnv.NodeInput[i].Name)
			}

			keystoneEnv.DONTopology[i] = &types.DONTopology{
				DON:        don,
				NodeInput:  keystoneEnv.NodeInput[i],
				NodeOutput: keystoneEnv.WrappedNodeOutput[i],
				ID:         libc.MustSafeUint32(i + 1),
				Flags:      flags,
			}
		}
	}

	maybeID, err := keystoneflags.OneDONTopologyWithFlag(keystoneEnv.DONTopology, types.WorkflowDON)
	if err != nil {
		return errors.Wrap(err, "failed to get workflow DON ID")
	}
	keystoneEnv.WorkflowDONID = maybeID.ID

	return nil
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
