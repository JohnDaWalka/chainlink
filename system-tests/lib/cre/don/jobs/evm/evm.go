package evm

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/pkg/errors"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
)

const flag = cre.EVMCapability
const evmConfigTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}","logTriggerPollInterval":{{.LogTriggerPollInterval}}, "creForwarderAddress":"{{.CreForwarderAddress}}","receiverGasMinimum":{{.ReceiverGasMinimum}},"nodeAddress":"{{.NodeAddress}}"}'`

// buildRuntimeValues creates runtime-generated  values for any keys not specified in TOML
func buildRuntimeValues(chainID uint64, networkFamily, creForwarderAddress, nodeAddress string) map[string]any {
	return map[string]any{
		"ChainID":             chainID,
		"NetworkFamily":       networkFamily,
		"CreForwarderAddress": creForwarderAddress,
		"NodeAddress":         nodeAddress,
	}
}

var JobSpecFn = func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
	configGen := func(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error) {
		cs, ok := chainsel.EvmChainIdToChainSelector()[chainID]
		if !ok {
			return "", fmt.Errorf("chain selector not found for chainID: %d", chainID)
		}

		creForwarderKey := datastore.NewAddressRefKey(
			cs,
			datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
			semver.MustParse("1.0.0"),
			"",
		)
		creForwarderAddress, err := input.CldEnvironment.DataStore.Addresses().Get(creForwarderKey)
		if err != nil {
			return "", errors.Wrap(err, "failed to get CRE Forwarder address")
		}

		logger.Debug().Msgf("Found CRE Forwarder contract on chain %d at %s", chainID, creForwarderAddress.Address)

		// Build runtime fallbacks for any missing values
		runtimeFallbacks := buildRuntimeValues(chainID, "evm", creForwarderAddress.Address, nodeAddress)

		// Apply runtime fallbacks only for keys not specified by user
		templateData := don.ApplyRuntimeValues(mergedConfig, runtimeFallbacks)

		// Parse and execute template
		tmpl, err := template.New("evmConfig").Parse(evmConfigTemplate)
		if err != nil {
			return "", errors.Wrapf(err, "failed to parse %s config template", flag)
		}

		var configBuffer bytes.Buffer
		if err := tmpl.Execute(&configBuffer, templateData); err != nil {
			return "", errors.Wrapf(err, "failed to execute %s config template", flag)
		}

		configStr := configBuffer.String()

		if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
			return "", errors.Wrapf(err, "%s template validation failed", flag)
		}

		return configStr, nil
	}

	return ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
		input.DonTopology,
		input.CldEnvironment.DataStore,
		input.CapabilitiesAwareNodeSets,
		input.InfraInput,
		"capability_evm",
		flag,
		ocr.CapabilityAppliesPerChainsFn,
		ocr.EnabledPerChainFn,
		configGen,
		ocr.ConfigMergePerChainFn,
		input.CapabilityConfigs,
	)
}
