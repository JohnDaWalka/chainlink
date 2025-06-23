package environment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

type SetupJobsSeqDeps struct {
	Logger                    zerolog.Logger
	JdInput                   jd.Input
	NixShell                  *nix.Shell
	HomeChainBlockchainOutput *blockchain.Output
	Topology                  *cretypes.Topology
}

type SetupJobsSeqInput struct {
	InfraType                 libtypes.InfraType
	CapabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet
}

type SetupJobsSeqOutput struct {
	JdOutput      *jd.Output
	NodeSetOutput []*cretypes.WrappedNodeOutput
}

var SetupJobsSeq = operations.NewSequence[SetupJobsSeqInput, SetupJobsSeqOutput, SetupJobsSeqDeps](
	"setup-jobs-seq",
	semver.MustParse("1.0.0"),
	"Setup Jobs Sequentially",
	func(b operations.Bundle, deps SetupJobsSeqDeps, input SetupJobsSeqInput) (SetupJobsSeqOutput, error) {
		var jdOutput *jd.Output
		jdAndDonsErrGroup := &errgroup.Group{}

		jdAndDonsErrGroup.Go(func() error {
			startJDInput := StartJDOpInput{
				InfraType: input.InfraType,
			}
			startJDDeps := StartJDOpDeps{
				Logger:   deps.Logger,
				JdInput:  deps.JdInput,
				NixShell: deps.NixShell,
			}
			startJDOutput, startJDErr := operations.ExecuteOperation(b, StartJDOp, startJDDeps, startJDInput)
			if startJDErr != nil {
				return pkgerrors.Wrap(startJDErr, "failed to start Job Distributor")
			}

			jdOutput = startJDOutput.Output.JdOutput

			return nil
		})

		nodeSetOutput := make([]*cretypes.WrappedNodeOutput, 0, len(input.CapabilitiesAwareNodeSets))

		jdAndDonsErrGroup.Go(func() error {
			startDONsInput := StartDONsOpInput(input)
			startDONsDeps := StartDONsOpDeps{
				Logger:                    deps.Logger,
				Topology:                  deps.Topology,
				NixShell:                  deps.NixShell,
				HomeChainBlockchainOutput: deps.HomeChainBlockchainOutput,
			}
			startDONsOutput, startDonsErr := operations.ExecuteOperation(b, StartDONsOp, startDONsDeps, startDONsInput)
			if startDonsErr != nil {
				return pkgerrors.Wrap(startDonsErr, "failed to start DONs")
			}

			nodeSetOutput = startDONsOutput.Output.NodeSetOutput

			return nil
		})

		if jdAndDonErr := jdAndDonsErrGroup.Wait(); jdAndDonErr != nil {
			return SetupJobsSeqOutput{}, pkgerrors.Wrap(jdAndDonErr, "failed to start Job Distributor or DONs")
		}

		return SetupJobsSeqOutput{
			JdOutput:      jdOutput,
			NodeSetOutput: nodeSetOutput,
		}, nil
	},
)

type StartJDOpDeps struct {
	Logger   zerolog.Logger
	JdInput  jd.Input
	NixShell *nix.Shell
}

type StartJDOpInput struct {
	InfraType libtypes.InfraType
}

type StartJDOpOutput struct {
	JdOutput *jd.Output
}

var StartJDOp = operations.NewOperation[StartJDOpInput, StartJDOpOutput, StartJDOpDeps](
	"start-jd-op",
	semver.MustParse("1.0.0"),
	"Start Job Distributor",
	func(b operations.Bundle, deps StartJDOpDeps, input StartJDOpInput) (StartJDOpOutput, error) {
		jdStartTime := time.Now()
		deps.Logger.Info().Msg("Starting Job Distributor")

		var jdOutput *jd.Output
		if input.InfraType == libtypes.CRIB {
			deployCribJdInput := &cretypes.DeployCribJdInput{
				JDInput:        &deps.JdInput,
				NixShell:       deps.NixShell,
				CribConfigsDir: cribConfigsDir,
			}

			var jdErr error
			deps.JdInput.Out, jdErr = crib.DeployJd(deployCribJdInput)
			if jdErr != nil {
				return StartJDOpOutput{}, pkgerrors.Wrap(jdErr, "failed to deploy JD with devspace")
			}
		}

		var jdErr error
		jdOutput, jdErr = CreateJobDistributor(&deps.JdInput)
		if jdErr != nil {
			jdErr = fmt.Errorf("failed to start JD container for image %s: %w", deps.JdInput.Image, jdErr)

			// useful end user messages
			if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
				jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
			}
			return StartJDOpOutput{}, jdErr
		}

		deps.Logger.Info().Msgf("Job Distributor started in %.2f seconds", time.Since(jdStartTime).Seconds())

		return StartJDOpOutput{JdOutput: jdOutput}, nil
	},
)

type StartDONsOpDeps struct {
	Logger                    zerolog.Logger
	Topology                  *cretypes.Topology
	NixShell                  *nix.Shell
	HomeChainBlockchainOutput *blockchain.Output
}

type StartDONsOpInput struct {
	InfraType                 libtypes.InfraType
	CapabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet
}

type StartDONsOpOutput struct {
	NodeSetOutput []*cretypes.WrappedNodeOutput
}

var StartDONsOp = operations.NewOperation[StartDONsOpInput, StartDONsOpOutput, StartDONsOpDeps](
	"start-dons-op",
	semver.MustParse("1.0.0"),
	"Start DONs",
	func(b operations.Bundle, deps StartDONsOpDeps, input StartDONsOpInput) (StartDONsOpOutput, error) {
		startTimeDONs := time.Now()
		deps.Logger.Info().Msgf("Starting %d DONs", len(input.CapabilitiesAwareNodeSets))

		if input.InfraType == libtypes.CRIB {
			deps.Logger.Info().Msg("Saving node configs and secret overrides")
			deployCribDonsInput := &cretypes.DeployCribDonsInput{
				Topology:       deps.Topology,
				NodeSetInputs:  input.CapabilitiesAwareNodeSets,
				NixShell:       deps.NixShell,
				CribConfigsDir: cribConfigsDir,
			}

			var devspaceErr error
			input.CapabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
			if devspaceErr != nil {
				return StartDONsOpOutput{}, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with devspace")
			}
		}

		nodeSetOutput := make([]*cretypes.WrappedNodeOutput, 0, len(input.CapabilitiesAwareNodeSets))

		// TODO we could parallelize this as well in the future, but for single DON env this doesn't matter
		for _, nodeSetInput := range input.CapabilitiesAwareNodeSets {
			nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, deps.HomeChainBlockchainOutput)
			if nodesetErr != nil {
				return StartDONsOpOutput{}, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
			}

			nodeSetOutput = append(nodeSetOutput, &cretypes.WrappedNodeOutput{
				Output:       nodeset,
				NodeSetName:  nodeSetInput.Name,
				Capabilities: nodeSetInput.Capabilities,
			})
		}

		deps.Logger.Info().Msgf("DONs started in %.2f seconds", time.Since(startTimeDONs).Seconds())

		return StartDONsOpOutput{NodeSetOutput: nodeSetOutput}, nil
	},
)

type CreateJobsWithJdOpDeps struct {
	Logger                    zerolog.Logger
	SingleFileLogger          common.Logger
	HomeChainBlockchainOutput *blockchain.Output
	AddressBook               deployment.AddressBook
	JobSpecFactoryFunctions   []cretypes.JobSpecFactoryFn
	FullCLDEnvOutput          *cretypes.FullCLDEnvironmentOutput
}

type CreateJobsWithJdOpInput struct {
}

type CreateJobsWithJdOpOutput struct {
}

var CreateJobsWithJdOp = operations.NewOperation[CreateJobsWithJdOpInput, CreateJobsWithJdOpOutput, CreateJobsWithJdOpDeps](
	"create-jobs-op",
	semver.MustParse("1.0.0"),
	"Create Jobs",
	func(b operations.Bundle, deps CreateJobsWithJdOpDeps, input CreateJobsWithJdOpInput) (CreateJobsWithJdOpOutput, error) {
		createJobsStartTime := time.Now()
		deps.Logger.Info().Msg("Creating jobs with Job Distributor")

		donToJobSpecs := make(cretypes.DonsToJobSpecs)

		for _, jobSpecGeneratingFn := range deps.JobSpecFactoryFunctions {
			singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cretypes.JobSpecFactoryInput{
				CldEnvironment:   deps.FullCLDEnvOutput.Environment,
				BlockchainOutput: deps.HomeChainBlockchainOutput,
				DonTopology:      deps.FullCLDEnvOutput.DonTopology,
				AddressBook:      deps.AddressBook,
			})
			if jobSpecsErr != nil {
				return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
			}
			mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
		}

		createJobsInput := cretypes.CreateJobsInput{
			CldEnv:        deps.FullCLDEnvOutput.Environment,
			DonTopology:   deps.FullCLDEnvOutput.DonTopology,
			DonToJobSpecs: donToJobSpecs,
		}

		jobsErr := libdon.CreateJobs(b.GetContext(), deps.Logger, createJobsInput)
		if jobsErr != nil {
			return CreateJobsWithJdOpOutput{}, pkgerrors.Wrap(jobsErr, "failed to create jobs")
		}

		deps.Logger.Info().Msgf("Jobs created in %.2f seconds", time.Since(createJobsStartTime).Seconds())

		return CreateJobsWithJdOpOutput{}, nil
	},
)
