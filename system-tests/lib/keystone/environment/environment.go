package environment

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

func BuildTopologyAndCLDEnvironment(lgr logger.Logger, keystoneEnv *types.KeystoneEnvironment) error {
	err := buildChainlinkDeploymentEnv(lgr, keystoneEnv)
	if err != nil {
		return errors.Wrap(err, "failed to build chainlink deployment environment")
	}
	err = don.BuildDONTopology(keystoneEnv)
	if err != nil {
		return errors.Wrap(err, "failed to build DON topology")
	}

	return nil
}

func buildChainlinkDeploymentEnv(lgr logger.Logger, keystoneEnv *types.KeystoneEnvironment) error {
	if keystoneEnv == nil {
		return errors.New("keystone environment must be set")
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
	if len(keystoneEnv.Blockchain.Nodes) < 1 {
		return errors.New("expected at least one node in the blockchain output")
	}
	if len(keystoneEnv.WrappedNodeOutput) < 1 {
		return errors.New("expected at least one node in the wrapped node output")
	}

	envs := make([]*deployment.Environment, len(keystoneEnv.WrappedNodeOutput))
	keystoneEnv.Dons = make([]*devenv.DON, len(keystoneEnv.WrappedNodeOutput))

	for i, nodeOutput := range keystoneEnv.WrappedNodeOutput {
		// assume that each nodeset has only one bootstrap node
		nodeInfo, err := getNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, 1)
		if err != nil {
			return errors.Wrap(err, "failed to get node info")
		}

		jdConfig := devenv.JDConfig{
			GRPC:     keystoneEnv.JD.HostGRPCUrl,
			WSRPC:    keystoneEnv.JD.DockerWSRPCUrl,
			Creds:    insecure.NewCredentials(),
			NodeInfo: nodeInfo,
		}

		devenvConfig := devenv.EnvironmentConfig{
			JDConfig: jdConfig,
			Chains: []devenv.ChainConfig{
				{
					ChainID:   keystoneEnv.SethClient.Cfg.Network.ChainID,
					ChainName: keystoneEnv.SethClient.Cfg.Network.Name,
					ChainType: strings.ToUpper(keystoneEnv.Blockchain.Family),
					WSRPCs: []devenv.CribRPCs{{
						External: keystoneEnv.Blockchain.Nodes[0].HostWSUrl,
						Internal: keystoneEnv.Blockchain.Nodes[0].DockerInternalWSUrl,
					}},
					HTTPRPCs: []devenv.CribRPCs{{
						External: keystoneEnv.Blockchain.Nodes[0].HostHTTPUrl,
						Internal: keystoneEnv.Blockchain.Nodes[0].DockerInternalHTTPUrl,
					}},
					DeployerKey: keystoneEnv.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
				},
			},
		}

		env, don, err := devenv.NewEnvironment(context.Background, lgr, devenvConfig)
		if err != nil {
			return errors.Wrap(err, "failed to create environment")
		}

		envs[i] = env
		keystoneEnv.Dons[i] = don
	}

	var nodeIDs []string
	for _, env := range envs {
		nodeIDs = append(nodeIDs, env.NodeIDs...)
	}

	// we assume that all DONs run on the same chain and that there's only one chain
	// also, we don't care which instance of offchain client we use, because we have
	// only one instance of offchain client and we have just configured it to work
	// with nodes from all DONs
	keystoneEnv.Environment = &deployment.Environment{
		Name:              envs[0].Name,
		Logger:            envs[0].Logger,
		ExistingAddresses: envs[0].ExistingAddresses,
		Chains:            envs[0].Chains,
		Offchain:          envs[0].Offchain,
		OCRSecrets:        envs[0].OCRSecrets,
		GetContext:        envs[0].GetContext,
		NodeIDs:           nodeIDs,
	}

	return nil
}

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
// TODO: remove this once the PR is merged and import his function
// IMPORTANT ADDITION:  prefix to differentiate between the different DONs
func getNodeInfo(nodeOut *ns.Output, prefix string, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.DockerP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}
		if i <= bootstrapNodeCount {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("%s_bootstrap-%d", prefix, i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		} else {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("%s_node-%d", prefix, i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		}
	}
	return nodeInfo, nil
}
