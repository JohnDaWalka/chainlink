package don

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

func Configure(t *testing.T, testLogger zerolog.Logger, keystoneEnv *types.KeystoneEnvironment, donToJobSpecs types.DonsToJobSpecs, donToConfigOverrides types.DonsToConfigOverrides) error {
	if keystoneEnv == nil {
		return errors.New("keystone environment must not be nil")
	}

	for i, donTopology := range keystoneEnv.MustDONTopology() {
		if configOverrides, ok := donToConfigOverrides[donTopology.ID]; ok {
			for j, configOverride := range configOverrides {
				if len(donTopology.NodeInput.NodeSpecs)-1 < j {
					return errors.Errorf("config override index out of bounds: %d", j)
				}
				donTopology.NodeInput.NodeSpecs[j].Node.TestConfigOverrides = configOverride
			}
			var setErr error
			keystoneEnv.MustDONTopology()[i].NodeOutput, setErr = config.Set(t, donTopology.NodeInput, keystoneEnv.MustBlockchain())
			if setErr != nil {
				return errors.Wrap(setErr, "failed to set node output")
			}
		}
	}

	nodeOutputs := make([]*types.WrappedNodeOutput, 0, len(keystoneEnv.MustDONTopology()))
	for i := range keystoneEnv.MustDONTopology() {
		nodeOutputs = append(nodeOutputs, keystoneEnv.MustDONTopology()[i].NodeOutput)
	}

	// after restarting the nodes, we need to reinitialize the JD clients otherwise
	// communication between JD and nodes will fail due to invalidated session cookie
	var jdErr error
	keystoneEnv.Environment, jdErr = jobs.ReinitialiseJDClients(keystoneEnv.MustCLDEnvironment(), keystoneEnv.MustJD(), nodeOutputs...)
	if jdErr != nil {
		return errors.Wrap(jdErr, "failed to reinitialize JD clients")
	}
	for _, donTopology := range keystoneEnv.MustDONTopology() {
		if jobSpecs, ok := donToJobSpecs[donTopology.ID]; ok {
			createErr := jobs.Create(keystoneEnv.MustCLDEnvironment().Offchain, donTopology.DON, donTopology.Flags, jobSpecs)
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

	keystoneEnv.DONTopology = make([]*types.DONTopology, len(keystoneEnv.MustDons()))

	// one DON to do everything
	if len(keystoneEnv.MustDons()) == 1 {
		flags, err := flags.NodeSetFlags(keystoneEnv.MustNodeInput()[0])
		if err != nil {
			return errors.Wrapf(err, "failed to convert string flags to bitmap for nodeset %s", keystoneEnv.MustNodeInput()[0].Name)
		}

		keystoneEnv.DONTopology[0] = &types.DONTopology{
			DON:        keystoneEnv.MustDons()[0],
			NodeInput:  keystoneEnv.MustNodeInput()[0],
			NodeOutput: keystoneEnv.MustWrappedNodeOutput()[0],
			ID:         1,
			Flags:      flags,
		}
	} else {
		for i, don := range keystoneEnv.MustDons() {
			flags, err := flags.NodeSetFlags(keystoneEnv.MustNodeInput()[i])
			if err != nil {
				return errors.Wrapf(err, "failed to convert string flags to bitmap for nodeset %s", keystoneEnv.MustNodeInput()[i].Name)
			}

			keystoneEnv.DONTopology[i] = &types.DONTopology{
				DON:        don,
				NodeInput:  keystoneEnv.MustNodeInput()[i],
				NodeOutput: keystoneEnv.MustWrappedNodeOutput()[i],
				ID:         libc.MustSafeUint32(i + 1),
				Flags:      flags,
			}
		}
	}

	maybeID, err := flags.OneDONTopologyWithFlag(keystoneEnv.MustDONTopology(), types.WorkflowDON)
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

func Start(nsInputs []*types.CapabilitiesAwareNodeSet, keystoneEnv *types.KeystoneEnvironment) error {
	if keystoneEnv == nil {
		return errors.New("keystone environment must be set")
	}

	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range nsInputs {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range nsInputs[i].NodeSpecs {
				nsInputs[i].NodeSpecs[j].Node.Image = image
			}
		}
	}

	for _, nsInput := range nsInputs {
		nodeset, err := ns.NewSharedDBNodeSet(nsInput.Input, keystoneEnv.MustBlockchain())
		if err != nil {
			return errors.Wrap(err, "failed to deploy node set")
		}

		keystoneEnv.NodeInput = append(keystoneEnv.MustNodeInput(), nsInput)
		keystoneEnv.WrappedNodeOutput = append(keystoneEnv.MustWrappedNodeOutput(), &types.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nsInput.Name,
			Capabilities: nsInput.Capabilities,
		})
	}

	return nil
}
