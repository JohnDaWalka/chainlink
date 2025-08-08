package jobs

import (
	"fmt"

	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

var (
	DefaultAllowedPorts = []int{80, 443}
)

type HandlerType string

const (
	WebAPIHandlerType HandlerType = "web-api-capabilities"
	HTTPHandlerType   HandlerType = "http-capabilities"
)

func BootstrapOCR3(nodeID string, name string, ocr3CapabilityAddress string, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "bootstrap"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	contractID = "%s"
	contractConfigTrackerPollInterval = "1s"
	contractConfigConfirmations = 1
	relay = "evm"
	[relayConfig]
	chainID = %d
	providerType = "ocr3-capability"
`,
			uuid,
			"ocr3-bootstrap-"+name,
			ocr3CapabilityAddress,
			chainID),
	}
}

const (
	EmptyStdCapConfig = "\"\""
)

func WorkerStandardCapability(nodeID, name, command, config, oracleFactoryConfig string) *jobv1.ProposeJobRequest {
	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "standardcapabilities"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	forwardingAllowed = false
	command = "%s"
	config = %s
	%s
`,
			uuid.NewString(),
			name,
			command,
			config,
			oracleFactoryConfig),
	}
}

func WorkerOCR3(nodeID string, ocr3CapabilityAddress, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData cre.OCRPeeringData, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "offchainreporting2"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	contractID = "%s"
	ocrKeyBundleID = "%s"
	p2pv2Bootstrappers = [
		"%s@%s",
	]
	relay = "evm"
	pluginType = "plugin"
	transmitterID = "%s"
	[relayConfig]
	chainID = "%d"
	[pluginConfig]
	command = "/usr/local/bin/chainlink-ocr3-capability"
	ocrVersion = 3
	pluginName = "ocr-capability"
	providerType = "ocr3-capability"
	telemetryType = "plugin"
	[onchainSigningStrategy]
	strategyName = 'multi-chain'
	[onchainSigningStrategy.config]
	evm = "%s"
`,
			uuid,
			cre.OCR3Capability,
			ocr3CapabilityAddress,
			ocr2KeyBundleID,
			ocrPeeringData.OCRBootstraperPeerID,
			fmt.Sprintf("%s:%d", ocrPeeringData.OCRBootstraperHost, ocrPeeringData.Port),
			nodeEthAddress,
			chainID,
			ocr2KeyBundleID,
		),
	}
}

func WorkerVaultOCR3(nodeID string, vaultCapabilityAddress, nodeEthAddress, ocr2KeyBundleID string, ocrPeeringData cre.OCRPeeringData, chainID uint64, masterPublicKey string, encryptedPrivateKeyShare string) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "offchainreporting2"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	contractID = "%s"
	ocrKeyBundleID = "%s"
	p2pv2Bootstrappers = [
		"%s@%s",
	]
	relay = "evm"
	pluginType = "%s"
	transmitterID = "%s"
	[relayConfig]
	chainID = "%d"
	[pluginConfig]
	requestExpiryDuration = "60s"
	[pluginConfig.dkg]
	masterPublicKey = "%s"
	encryptedPrivateKeyShare = "%s"
`,
			uuid,
			"Vault OCR3 Capability",
			vaultCapabilityAddress,
			ocr2KeyBundleID,
			ocrPeeringData.OCRBootstraperPeerID,
			fmt.Sprintf("%s:%d", ocrPeeringData.OCRBootstraperHost, ocrPeeringData.Port),
			types.VaultPlugin,
			nodeEthAddress,
			chainID,
			masterPublicKey,
			encryptedPrivateKeyShare,
		),
	}
}
