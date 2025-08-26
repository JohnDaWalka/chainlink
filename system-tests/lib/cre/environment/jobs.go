package environment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func StartJD(lggr zerolog.Logger, nixShell *nix.Shell, jdInput jd.Input, infraType libtypes.InfraType) (*jd.Output, error) {
	startTime := time.Now()
	lggr.Info().Msg("Starting Job Distributor")
	lggr.Info().Msgf("StartJD: Infrastructure type: %s", infraType)
	lggr.Info().Msgf("StartJD: JD input image: %s", jdInput.Image)
	lggr.Info().Msgf("StartJD: JD input CSA encryption key length: %d", len(jdInput.CSAEncryptionKey))

	var jdOutput *jd.Output
	if infraType == libtypes.CRIB {
		lggr.Info().Msg("StartJD: Using CRIB infrastructure, deploying JD with devspace")
		deployCribJdInput := &cretypes.DeployCribJdInput{
			JDInput:        &jdInput,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var jdErr error
		jdInput.Out, jdErr = crib.DeployJd(deployCribJdInput)
		if jdErr != nil {
			lggr.Error().Err(jdErr).Msg("StartJD: Failed to deploy JD with devspace")
			return nil, pkgerrors.Wrap(jdErr, "failed to deploy JD with devspace")
		}
		lggr.Info().Msg("StartJD: Successfully deployed JD with devspace")
	}

	lggr.Info().Msg("StartJD: Creating Job Distributor container")
	var jdErr error
	jdOutput, jdErr = CreateJobDistributor(&jdInput)
	if jdErr != nil {
		lggr.Error().Err(jdErr).Msgf("StartJD: Failed to start JD container for image %s", jdInput.Image)
		jdErr = fmt.Errorf("failed to start JD container for image %s: %w", jdInput.Image, jdErr)

		// useful end user messages
		if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
			lggr.Error().Msg("StartJD: Docker pull access denied - ensure you have built the local image or are logged into AWS")
			jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
		}
		return nil, jdErr
	}

	lggr.Info().Msgf("Job Distributor started in %.2f seconds", time.Since(startTime).Seconds())
	lggr.Info().Msg("StartJD: Job Distributor container created successfully")

	return jdOutput, nil
}

func SetupJobs(
	lggr zerolog.Logger,
	jdInput jd.Input,
	nixShell *nix.Shell,
	registryChainBlockchainOutput *blockchain.Output,
	topology *cretypes.Topology,
	infraType libtypes.InfraType,
	capabilitiesAwareNodeSets []*cretypes.CapabilitiesAwareNodeSet,
) (*jd.Output, []*cretypes.WrappedNodeOutput, error) {
	lggr.Info().Msgf("SetupJobs: Starting with %d capabilities aware node sets", len(capabilitiesAwareNodeSets))
	lggr.Info().Msgf("SetupJobs: Infrastructure type: %s", infraType)
	lggr.Info().Msgf("SetupJobs: JD input image: %s", jdInput.Image)

	var jdOutput *jd.Output
	jdAndDonsErrGroup := &errgroup.Group{}

	jdAndDonsErrGroup.Go(func() error {
		lggr.Info().Msg("SetupJobs: Starting Job Distributor in background")
		var startJDErr error
		jdOutput, startJDErr = StartJD(lggr, nixShell, jdInput, infraType)
		if startJDErr != nil {
			lggr.Error().Err(startJDErr).Msg("SetupJobs: Failed to start Job Distributor")
			return pkgerrors.Wrap(startJDErr, "failed to start Job Distributor")
		}
		lggr.Info().Msg("SetupJobs: Job Distributor started successfully")
		return nil
	})

	nodeSetOutput := make([]*cretypes.WrappedNodeOutput, 0, len(capabilitiesAwareNodeSets))

	jdAndDonsErrGroup.Go(func() error {
		lggr.Info().Msg("SetupJobs: Starting DONs in background")
		var startDonsErr error
		nodeSetOutput, startDonsErr = StartDONs(lggr, nixShell, topology, infraType, registryChainBlockchainOutput, capabilitiesAwareNodeSets)
		if startDonsErr != nil {
			lggr.Error().Err(startDonsErr).Msg("SetupJobs: Failed to start DONs")
			return pkgerrors.Wrap(startDonsErr, "failed to start DONs")
		}
		lggr.Info().Msgf("SetupJobs: DONs started successfully, got %d node set outputs", len(nodeSetOutput))
		return nil
	})

	lggr.Info().Msg("SetupJobs: Waiting for Job Distributor and DONs to complete")
	if jdAndDonErr := jdAndDonsErrGroup.Wait(); jdAndDonErr != nil {
		lggr.Error().Err(jdAndDonErr).Msg("SetupJobs: Failed to start Job Distributor or DONs")
		return nil, nil, pkgerrors.Wrap(jdAndDonErr, "failed to start Job Distributor or DONs")
	}

	lggr.Info().Msg("SetupJobs: Job Distributor and DONs started successfully")
	return jdOutput, nodeSetOutput, nil
}
