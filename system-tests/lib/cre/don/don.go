package don

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func CreateJobs(testLogger zerolog.Logger, input cretypes.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, don := range input.DonTopology.Dons {
		if jobSpecs, ok := input.DonToJobSpecs[don.ID]; ok {
			createErr := jobs.Create(input.CldEnv.Offchain, don.DON, don.Flags, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", don.ID)
		}
	}

	return nil
}

func BuildTopology(nodeSetInput []*cretypes.CapabilitiesAwareNodeSet) (*cretypes.Topology, error) {
	donsWithMetadata := make([]*cretypes.DonMetadata, len(nodeSetInput))

	// one DON to do everything
	if len(nodeSetInput) == 1 {
		flags, err := flags.NodeSetFlags(nodeSetInput[0])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert string flags to bitmap for nodeset %s", nodeSetInput[0].Name)
		}

		donsWithMetadata[0] = &cretypes.DonMetadata{
			ID:            1,
			Flags:         flags,
			NodesMetadata: make([]*cretypes.NodeMetadata, len(nodeSetInput[0].NodeSpecs)),
			Name:          nodeSetInput[0].Name,
		}
	} else {
		for i := range nodeSetInput {
			flags, err := flags.NodeSetFlags(nodeSetInput[i])
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert string flags to bitmap for nodeset %s", nodeSetInput[i].Name)
			}

			donsWithMetadata[i] = &cretypes.DonMetadata{
				ID:            libc.MustSafeUint32(i + 1),
				Flags:         flags,
				NodesMetadata: make([]*cretypes.NodeMetadata, len(nodeSetInput[i].NodeSpecs)),
				Name:          nodeSetInput[0].Name,
			}
		}
	}

	for i, donMetadata := range donsWithMetadata {
		for j := range donMetadata.NodesMetadata {
			nodeWithLabels := cretypes.NodeMetadata{}
			nodeType := devenv.NodeLabelValuePlugin
			if nodeSetInput[i].BootstrapNodeIndex != -1 && j == nodeSetInput[i].BootstrapNodeIndex {
				nodeType = devenv.NodeLabelValueBootstrap
			}
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &ptypes.Label{
				Key:   devenv.NodeLabelKeyType,
				Value: ptr.Ptr(nodeType),
			})
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &ptypes.Label{
				Key:   libnode.IndexKey,
				Value: ptr.Ptr(fmt.Sprint(j)),
			})
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &ptypes.Label{
				Key: libnode.HostLabelKey,
				// TODO this will only work with Docker, for CRIB we need a different approach
				Value: ptr.Ptr(fmt.Sprintf("%s-node%d", donMetadata.Name, j)),
			})

			donsWithMetadata[i].NodesMetadata[j] = &nodeWithLabels
		}
	}

	maybeID, err := flags.OneDonMetadataWithFlag(donsWithMetadata, cretypes.WorkflowDON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DON ID")
	}

	return &cretypes.Topology{
		Metadata:      donsWithMetadata,
		WorkflowDONID: maybeID.ID,
	}, nil
}

func GenereteKeys(input *cretypes.GenerateKeysInput) (*cretypes.Topology, *cretypes.GenerateKeysOutput, error) {
	if input == nil {
		return nil, nil, errors.New("input is nil")
	}

	if err := input.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "input validation failed")
	}

	output := &cretypes.GenerateKeysOutput{
		EVMKeys: make(cretypes.DonsToEVMKeys),
		P2PKeys: make(cretypes.DonsToP2PKeys),
	}

	for _, donMetadata := range input.Topology.Metadata {
		if input.GenerateP2PKeys {
			p2pKeys, err := crypto.GenerateP2PKeys(input.Password, len(donMetadata.NodesMetadata))
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to generate P2P keys")
			}
			output.P2PKeys[donMetadata.ID] = p2pKeys

			for idx, node := range donMetadata.NodesMetadata {
				node.Labels = append(node.Labels, &ptypes.Label{
					Key:   devenv.NodeLabelP2PIDType,
					Value: ptr.Ptr(p2pKeys.PeerIDs[idx]),
				})
			}
		}

		if input.EMVKeysToGenerate != nil {
			evmKeys, err := crypto.GenerateEVMKeys(input.Password, len(donMetadata.NodesMetadata))
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to generate EVM keys")
			}
			for _, chain := range input.EMVKeysToGenerate {
				evmKeys.Chains = append(evmKeys.Chains, types.EVMKeysToChains{
					ChainSelector: chain.ChainSelector,
					ChainName:     chain.ChainName,
				})
			}

			output.EVMKeys[donMetadata.ID] = evmKeys

			for idx, node := range donMetadata.NodesMetadata {
				node.Labels = append(node.Labels, &ptypes.Label{
					Key:   libnode.EthAddressKey,
					Value: ptr.Ptr(evmKeys.PublicAddresses[idx].Hex()),
				})
			}

		}
	}

	return input.Topology, output, nil
}

// In order to whitelist host IP in the gateway, we need to resolve the host.docker.internal to the host IP,
// and since CL image doesn't have dig or nslookup, we need to use curl.
func ResolveHostDockerInternaIP(testLogger zerolog.Logger, containerName string) (string, error) {
	if isCurlInstalled(containerName) {
		return resolveDockerHostWithCurl(containerName)
	} else if isNsLookupInstalled(containerName) {
		return resolveDockerHostWithNsLookup(containerName)
	}

	return "", errors.New("neither curl nor nslookup is installed")
}

func isNsLookupInstalled(containerName string) bool {
	cmd := []string{"which", "nslookup"}
	output, err := framework.ExecContainer(containerName, cmd)

	if err != nil || output == "" {
		return false
	}

	return true
}

func resolveDockerHostWithNsLookup(containerName string) (string, error) {
	cmd := []string{"nslookup", "host.docker.internal"}
	output, err := framework.ExecContainer(containerName, cmd)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`host.docker.internal(\n|\r)Address:\s+([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", errors.New("failed to extract IP address from curl output")
	}

	return matches[2], nil
}

func isCurlInstalled(containerName string) bool {
	cmd := []string{"which", "curl"}
	output, err := framework.ExecContainer(containerName, cmd)

	if err != nil || output == "" {
		return false
	}

	return true
}

func resolveDockerHostWithCurl(containerName string) (string, error) {
	cmd := []string{"curl", "-v", "http://host.docker.internal"}
	output, err := framework.ExecContainer(containerName, cmd)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`.*Trying ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+).*`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", errors.New("failed to extract IP address from curl output")
	}

	return matches[1], nil
}
