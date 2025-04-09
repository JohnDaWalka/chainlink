package devenv

import (
	"context"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/types"
)

type CCIPEnvironmentBuilder struct {
	jdOutput          *jd.Output
	blockchainOutputs types.ChainIDToBlockchainOutputs
	// todo: replace this with an array of transactOpts
	sethClients       []*seth.Client
	nodeSetOutput     *simple_node_set.Output
	existingAddresses deployment.AddressBook
	credentials       credentials.TransportCredentials
	logger            logger.Logger
	errs              []string
}

func NewCCIPEnvironmentBuilder(lgr logger.Logger) *CCIPEnvironmentBuilder {
	b := &CCIPEnvironmentBuilder{
		logger: lgr,
	}

	if lgr == nil {
		b.errs = append(b.errs, "logger not set")
	}
	return b
}

func (b *CCIPEnvironmentBuilder) WithJobDistributor(jdOutput *jd.Output, jdTransportCredentials credentials.TransportCredentials) *CCIPEnvironmentBuilder {
	if jdTransportCredentials == nil {
		b.errs = append(b.errs, "jd credentials not set")
	}
	if jdOutput == nil {
		b.errs = append(b.errs, "jd output not set")
		return b
	}
	if jdOutput.ExternalGRPCUrl == "" {
		b.errs = append(b.errs, "external gRPC url not set")
	}
	if jdOutput.InternalWSRPCUrl == "" {
		b.errs = append(b.errs, "internal wsRPC url not set")
	}

	b.jdOutput = jdOutput
	b.credentials = jdTransportCredentials
	return b
}

func (b *CCIPEnvironmentBuilder) WithBlockchains(blockchainOutputs types.ChainIDToBlockchainOutputs) *CCIPEnvironmentBuilder {
	if len(blockchainOutputs) == 0 {
		b.errs = append(b.errs, "blockchain outputs not set")
	}
	b.blockchainOutputs = blockchainOutputs
	return b
}

func (b *CCIPEnvironmentBuilder) WithSethClients(sethClients []*seth.Client) *CCIPEnvironmentBuilder {
	if len(sethClients) == 0 {
		b.errs = append(b.errs, "seth clients not set")
	}
	b.sethClients = sethClients
	return b
}

func (b *CCIPEnvironmentBuilder) WithNodeSet(nodeSetOutput *simple_node_set.Output) *CCIPEnvironmentBuilder {
	if nodeSetOutput == nil {
		b.errs = append(b.errs, "node set output not set")
	}
	b.nodeSetOutput = nodeSetOutput
	return b
}

func (b *CCIPEnvironmentBuilder) WithExistingAddresses(existingAddresses deployment.AddressBook) *CCIPEnvironmentBuilder {
	b.existingAddresses = existingAddresses
	return b
}

func (b *CCIPEnvironmentBuilder) Build() (*deployment.Environment, *DON, error) {
	if len(b.errs) > 0 {
		return nil, nil, errors.New("validation errors: " + strings.Join(b.errs, ", "))
	}
	if b.blockchainOutputs == nil {
		return nil, nil, errors.New("blockchain outputs not set")
	}
	if b.nodeSetOutput == nil {
		return nil, nil, errors.New("nodeSetOutput not set")
	}
	if b.jdOutput == nil {
		return nil, nil, errors.New("jd output not set")
	}

	chains := chainsFromBlockchainOutputs(b.sethClients, b.blockchainOutputs)

	// In CCIP we assume that there is one bootstrap node
	// prefix is hardcoded to simplify the setup
	allNodesInfo, err := GetNodeInfo(b.nodeSetOutput, "ccip", 1)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get node info")
	}

	jdConfig := JDConfig{
		GRPC:     b.jdOutput.ExternalGRPCUrl,
		WSRPC:    b.jdOutput.InternalWSRPCUrl,
		Creds:    b.credentials,
		NodeInfo: allNodesInfo,
	}

	devenvConfig := EnvironmentConfig{
		JDConfig: jdConfig,
		Chains:   chains,
	}

	b.logger.Infow("creating CLD environment")
	env, don, err := NewEnvironment(context.Background, b.logger, devenvConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create a CLD environment")
	}

	env.ExistingAddresses = b.existingAddresses

	return env, don, nil
}
