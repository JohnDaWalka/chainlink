package v2

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/consensus"
)

const flag = cre.ConsensusCapabilityV2

type Consensus struct{}

func (c *Consensus) Flag() cre.CapabilityFlag {
	return flag
}

func (c *Consensus) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "consensus",
			Version:        "1.0.0-alpha",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const ContractQualifier = "capability_consensus"

func (c *Consensus) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	_, ocr3ContractAddr, ocrErr := contracts.DeployOCR3Contract(testLogger, ContractQualifier, creEnv.RegistryChainSelector, creEnv.CldfEnvironment, creEnv.ContractVersions)
	if ocrErr != nil {
		return fmt.Errorf("failed to deploy OCR3 (consensus v2) contract %w", ocrErr)
	}

	jobsErr := createJobs(
		ctx,
		don,
		dons,
		*ocr3ContractAddr,
		creEnv,
	)
	if jobsErr != nil {
		return fmt.Errorf("failed to create OCR3 jobs: %w", jobsErr)
	}

	// wait for LP to be started (otherwise it won't pick up contract's configuration events)
	if err := consensus.WaitForLogPollerToBeHealthy(don); err != nil {
		return errors.Wrap(err, "failed while waiting for Log Poller to become healthy")
	}

	ocr3Config, ocr3confErr := contracts.DefaultOCR3Config()
	if ocr3confErr != nil {
		return fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
	}

	_, ocr3Err := operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env: creEnv.CldfEnvironment,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress: ocr3ContractAddr,
			ChainSelector:   creEnv.RegistryChainSelector,
			DON:             don.KeystoneDONConfig(),
			Config:          don.ResolveORC3Config(ocr3Config),
			DryRun:          false,
		},
	)

	if ocr3Err != nil {
		return errors.Wrap(ocr3Err, "failed to configure OCR3 contract")
	}

	return nil
}

const configTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}","nodeAddress":"{{.NodeAddress}}"}'`

func createJobs(
	ctx context.Context,
	don *cre.Don,
	dons *cre.Dons,
	contractAddress common.Address,
	creEnv *cre.Environment,
) error {
	jobSpecs := []*jobv1.ProposeJobRequest{}
	capabilityConfig, ok := creEnv.CapabilityConfigs[flag]
	if !ok {
		return fmt.Errorf("%s config not found in capabilities config: %v", flag, creEnv.CapabilityConfigs)
	}

	bootstrapNode, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	var nodeSet cre.NodeSetWithCapabilityConfigs
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == don.Name {
			nodeSet = ns
			break
		}
	}
	if nodeSet == nil {
		return fmt.Errorf("could not find node set for Don named '%s'", don.Name)
	}

	chainID, cErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if cErr != nil {
		return fmt.Errorf("failed to get chain ID from selector %d: %w", creEnv.RegistryChainSelector, cErr)
	}

	jobSpecs = append(jobSpecs, ocr.BootstrapJobSpec(bootstrapNode.JobDistributorDetails.NodeID, flag, contractAddress.Hex(), chainID))
	chainIDStr := strconv.FormatUint(chainID, 10)
	templateData := envconfig.ResolveCapabilityConfigForDON(flag, capabilityConfig.Config, nodeSet.GetCapabilityConfigOverrides())

	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryPath, creEnv.Provider)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get command for cron capability")
	}

	for _, workerNode := range workerNodes {
		evmKey, ok := workerNode.Keys.EVM[chainID]
		if !ok {
			return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
		}
		nodeAddress := evmKey.PublicAddress.Hex()

		evmKeyBundle, ok := workerNode.Keys.OCR2BundleIDs[chainselectors.FamilyEVM] // we can always expect evm bundle key id present since evm is the registry chain
		if !ok {
			return errors.New("failed to get key bundle id for evm family")
		}

		strategyName := "single-chain"
		if len(workerNode.Keys.OCR2BundleIDs) > 1 {
			strategyName = "multi-chain"
		}

		oracleFactoryConfigInstance := job.OracleFactoryConfig{
			Enabled:            true,
			ChainID:            chainIDStr,
			BootstrapPeers:     []string{fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrapNode.Keys.PeerID(), "p2p_"), bootstrapNode.Host, cre.OCRPeeringPort)},
			OCRContractAddress: contractAddress.Hex(),
			OCRKeyBundleID:     evmKeyBundle,
			TransmitterID:      nodeAddress,
			OnchainSigning: job.OnchainSigningStrategy{
				StrategyName: strategyName,
				Config:       workerNode.Keys.OCR2BundleIDs,
			},
		}

		// TODO: merge with jobConfig?
		type OracleFactoryConfigWrapper struct {
			OracleFactory job.OracleFactoryConfig `toml:"oracle_factory"`
		}
		wrapper := OracleFactoryConfigWrapper{OracleFactory: oracleFactoryConfigInstance}

		var oracleBuffer bytes.Buffer
		if errEncoder := toml.NewEncoder(&oracleBuffer).Encode(wrapper); errEncoder != nil {
			return errors.Wrap(errEncoder, "failed to encode oracle factory config to TOML")
		}
		oracleStr := strings.ReplaceAll(oracleBuffer.String(), "\n", "\n\t")

		runtimeFallbacks := buildRuntimeValues(chainID, "evm", nodeAddress)
		var aErr error
		templateData, aErr = credon.ApplyRuntimeValues(templateData, runtimeFallbacks)
		if aErr != nil {
			return errors.Wrap(aErr, "failed to apply runtime values")
		}

		tmpl, err := template.New("consensusConfig").Parse(configTemplate)
		if err != nil {
			return errors.Wrap(err, "failed to parse consensus config template")
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return errors.Wrap(err, "failed to execute consensus config template")
		}

		configStr := configBuffer.String()

		if err := credon.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return errors.Wrapf(err, "%s template validation failed", flag)
		}

		jobSpec := standardcapability.WorkerJobSpec(workerNode.JobDistributorDetails.NodeID, flag, command, configStr, oracleStr)
		jobSpec.Labels = []*ptypes.Label{{Key: cre.CapabilityLabelKey, Value: ptr.Ptr(flag)}}
		jobSpecs = append(jobSpecs, jobSpec)
	}

	jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, jobSpecs)
	if jobErr != nil {
		return fmt.Errorf("failed to create EVM OCR3 jobs for don %s: %w", don.Name, jobErr)
	}

	return nil
}

func buildRuntimeValues(chainID uint64, networkFamily, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":       chainID,
		"NetworkFamily": networkFamily,
		"NodeAddress":   nodeAddress,
	}
}
