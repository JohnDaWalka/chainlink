package consensus

import (
	"bytes"
	"html/template"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const flag = cre.ConsensusCapabilityV2
const configTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}","nodeAddress":"{{.NodeAddress}}"}'`

func New() (*capabilities.Capability, error) {
	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(jobSpec),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
	)
}

func registerWithV1(donFlags []string, _ *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, flag) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "consensus",
				Version:        "1.0.0",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}

func buildRuntimeValues(chainID uint64, networkFamily, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":       chainID,
		"NetworkFamily": networkFamily,
		"NodeAddress":   nodeAddress,
	}
}

type jobConfigGenerator struct {
	input *cre.JobSpecInput
}

func (c *jobConfigGenerator) Generate(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error) {
	// Build runtime fallbacks for any missing values
	runtimeFallbacks := buildRuntimeValues(chainID, "evm", nodeAddress)

	// Apply runtime fallbacks only for keys not specified by user
	templateData, aErr := don.ApplyRuntimeValues(mergedConfig, runtimeFallbacks)
	if aErr != nil {
		return "", errors.Wrap(aErr, "failed to apply runtime values")
	}

	// Parse and execute template
	tmpl, err := template.New("consensusConfig").Parse(configTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse consensus config template")
	}

	var configBuffer bytes.Buffer
	if err := tmpl.Execute(&configBuffer, templateData); err != nil {
		return "", errors.Wrap(err, "failed to execute consensus config template")
	}

	return configBuffer.String(), nil
}

func jobSpec(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
	return ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
		input.DonTopology,
		input.CldEnvironment.DataStore,
		input.CapabilitiesAwareNodeSets,
		input.InfraInput,
		"capability_consensus",
		flag,
		&ocr.CapabilityEnablerPerDon{},
		&ocr.RegistryChainOnlyProvider{},
		&jobConfigGenerator{input: input},
		&ocr.ConfigMergerPerDon{},
		input.CapabilityConfigs,
	)
}
