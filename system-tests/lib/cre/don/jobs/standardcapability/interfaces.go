package factory

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry"
)

// Type aliases for cleaner function signatures
type RuntimeValuesExtractor func(chainID uint64, nodeMetadata *cre.NodeMetadata) map[string]any
type RuntimeValueBuilder func(runtimeValues, mergedConfig map[string]any) map[string]any
type CommandBuilder func(input *cre.JobSpecInput, capabilityConfig cre.CapabilityConfig) (string, error)

// NoOpExtractor is a no-operation runtime values extractor for DON-level capabilities
// that don't need runtime values extraction from node metadata
var NoOpExtractor RuntimeValuesExtractor = func(chainID uint64, nodeMetadata *cre.NodeMetadata) map[string]any {
	return map[string]any{} // Return empty map - DON-level capabilities typically don't need runtime values
}

// BinaryPathBuilder constructs the container path for capability binaries by combining
// the default container directory with the base name of the capability's binary path
var BinaryPathBuilder CommandBuilder = func(input *cre.JobSpecInput, capabilityConfig cre.CapabilityConfig) (string, error) {
	containerPath, pathErr := creregistry.DefaultContainerDirectory(input.InfraInput.Type)
	if pathErr != nil {
		return "", errors.Wrapf(pathErr, "failed to get default container directory for infra type %s", input.InfraInput.Type)
	}

	return filepath.Join(containerPath, filepath.Base(capabilityConfig.BinaryPath)), nil
}

// JobSpecFactory is the base interface for all job specification factories
type JobSpecFactory interface {
	GenerateJobSpecs(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error)
}

// DonLevelJobSpecFactory generates job specs for capabilities that operate at the DON level
// without chain-specific configuration (e.g., cron, mock, custom-compute, web-api-*)
type DonLevelJobSpecFactory interface {
	JobSpecFactory
	CapabilityFlag() cre.CapabilityFlag
	ConfigTemplate() string
	RuntimeValuesExtractor() RuntimeValuesExtractor
	CommandBuilder() CommandBuilder
}

// ChainSpecificJobSpecFactory generates job specs for capabilities that require
// per-chain configuration and deployment (e.g., read-contract, log-event-trigger)
type ChainSpecificJobSpecFactory interface {
	JobSpecFactory
	CapabilityFlag() cre.CapabilityFlag
	ConfigTemplate() string
	RuntimeValuesExtractor() RuntimeValuesExtractor
	CommandBuilder() CommandBuilder
}
