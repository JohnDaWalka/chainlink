package environment

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	libdevenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type SetupJobsSeqInput struct {
	Topology                  *cretypes.Topology
	CapabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet
}

type SetupJobsSeqOutput struct {
	CapabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet
}

type SetupJobsSeqDeps struct {
	JdInput    jd.Input
	NixShell   *nix.Shell
	InfraInput libtypes.InfraInput
}

var SetupJobsSeq = operations.NewSequence[SetupJobsSeqInput, SetupJobsSeqOutput, SetupJobsSeqDeps](
	"setup-jobs-seq",
	semver.MustParse("1.0.0"),
	"Setup Jobs Sequence",
	func(b operations.Bundle, deps SetupJobsSeqDeps, input SetupJobsSeqInput) (SetupJobsSeqOutput, error) {
		// TODO: divide segments into operations

		// Create JD
		if deps.InfraInput.InfraType == libtypes.CRIB {
			deployCribJdInput := &cretypes.DeployCribJdInput{
				JDInput:        &deps.JdInput,
				NixShell:       deps.NixShell,
				CribConfigsDir: cribConfigsDir,
			}

			var jdErr error
			deps.JdInput.Out, jdErr = crib.DeployJd(deployCribJdInput)
			if jdErr != nil {
				return SetupJobsSeqOutput{}, pkgerrors.Wrap(jdErr, "failed to deploy JD with devspace")
			}
		}

		jdOutput, jdErr := CreateJobDistributor(&deps.JdInput)
		if jdErr != nil {
			jdErr = fmt.Errorf("failed to start JD container for image %s: %w", deps.JdInput.Image, jdErr)

			// useful end user messages
			if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
				jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
			}
			return SetupJobsSeqOutput{}, jdErr
		}

		// starting DONs
		if deps.InfraInput.InfraType == libtypes.CRIB {
			deployCribDonsInput := &cretypes.DeployCribDonsInput{
				Topology:       input.Topology,
				NodeSetInputs:  input.CapabilitiesAwareNodeSets,
				NixShell:       deps.NixShell,
				CribConfigsDir: cribConfigsDir,
			}

			var devspaceErr error
			input.CapabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
			if devspaceErr != nil {
				return SetupJobsSeqOutput{}, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with devspace")
			}
		}

		nodeSetOutput := make([]*cretypes.WrappedNodeOutput, 0, len(input.CapabilitiesAwareNodeSets))
		for _, nodeSetInput := range input.CapabilitiesAwareNodeSets {
			nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, homeChainOutput.BlockchainOutput)
			if nodesetErr != nil {
				return SetupJobsSeqOutput{}, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
			}

			nodeSetOutput = append(nodeSetOutput, &cretypes.WrappedNodeOutput{
				Output:       nodeset,
				NodeSetName:  nodeSetInput.Name,
				Capabilities: nodeSetInput.Capabilities,
			})
		}

		// Prepare the CLD environment that's required by the keystone changeset
		// Ugly glue hack ¯\_(ツ)_/¯
		fullCldInput := &cretypes.FullCLDEnvironmentInput{
			JdOutput:          jdOutput,
			BlockchainOutputs: bcOuts,
			SethClients:       sethClients,
			NodeSetOutput:     nodeSetOutput,
			ExistingAddresses: allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			Topology:          topology,
		}

		// We need to use TLS for CRIB, because it exposes HTTPS endpoints
		var creds credentials.TransportCredentials
		if deps.InfraInput.InfraType == libtypes.CRIB {
			creds = credentials.NewTLS(&tls.Config{
				MinVersion: tls.VersionTLS12,
			})
		} else {
			creds = insecure.NewCredentials()
		}

		fullCldOutput, cldErr := libdevenv.BuildFullCLDEnvironment(ctx, singeFileLogger, fullCldInput, creds)
		if cldErr != nil {
			return SetupJobsSeqOutput{}, pkgerrors.Wrap(cldErr, "failed to build full CLD environment")
		}

		// Funding Chainlink nodes
		// Fund the nodes
		concurrentNonceMap, concurrentNonceMapErr := NewConcurrentNonceMap(ctx, blockchainsOutput)
		if concurrentNonceMapErr != nil {
			return nil, pkgerrors.Wrap(concurrentNonceMapErr, "failed to create concurrent nonce map")
		}

		// Decrement the nonce for each chain, because we will increment it in the next loop
		for _, bcOut := range blockchainsOutput {
			concurrentNonceMap.Decrement(bcOut.ChainID)
		}

		errGroup := &errgroup.Group{}
		for _, metaDon := range fullCldOutput.DonTopology.DonsWithMetadata {
			for _, bcOut := range blockchainsOutput {
				for _, node := range metaDon.DON.Nodes {
					errGroup.Go(func() error {
						nodeAddress := node.AccountAddr[strconv.FormatUint(bcOut.ChainID, 10)]
						if nodeAddress == "" {
							return nil
						}

						nonce := concurrentNonceMap.Increment(bcOut.ChainID)

						_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, bcOut.SethClient, libtypes.FundsToSend{
							ToAddress:  common.HexToAddress(nodeAddress),
							Amount:     big.NewInt(5000000000000000000),
							PrivateKey: bcOut.SethClient.MustGetRootPrivateKey(),
							Nonce:      ptr.Ptr(nonce),
						})
						if fundingErr != nil {
							return pkgerrors.Wrapf(fundingErr, "failed to fund node %s", nodeAddress)
						}
						return nil
					})
				}
			}
		}

		if err := errGroup.Wait(); err != nil {
			return SetupJobsSeqOutput{}, pkgerrors.Wrap(err, "failed to fund nodes")
		}

		// Creating jobs with Job Distributor
		donToJobSpecs := make(cretypes.DonsToJobSpecs)

		for _, jobSpecGeneratingFn := range input.JobSpecFactoryFunctions {
			singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cretypes.JobSpecFactoryInput{
				CldEnvironment:   fullCldOutput.Environment,
				BlockchainOutput: homeChainOutput.BlockchainOutput,
				DonTopology:      fullCldOutput.DonTopology,
				AddressBook:      allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			})
			if jobSpecsErr != nil {
				return SetupJobsSeqOutput{}, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
			}
			mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
		}

		createJobsInput := cretypes.CreateJobsInput{
			CldEnv:        fullCldOutput.Environment,
			DonTopology:   fullCldOutput.DonTopology,
			DonToJobSpecs: donToJobSpecs,
		}

		jobsErr := libdon.CreateJobs(ctx, testLogger, createJobsInput)
		if jobsErr != nil {
			return SetupJobsSeqOutput{}, pkgerrors.Wrap(jobsErr, "failed to create jobs")
		}

		// TODO: return new `CapabilitiesAwareNodeSets`

		return SetupJobsSeqOutput{}, nil
	},
)
