package contracts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/deployment"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"

	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"

	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"

	keystonenode "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/node"
)

func ConfigureKeystone(t *testing.T, keystoneEnv *types.KeystoneEnvironment) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.DONTopology, "DON topology must be set")
	require.NotEmpty(t, keystoneEnv.ChainSelector, "chain selector must be set")
	require.NotNil(t, keystoneEnv.Environment, "environment must be set")

	donCapabilities := make([]keystone_changeset.DonCapabilities, 0, len(keystoneEnv.DONTopology))

	for _, donTopology := range keystoneEnv.DONTopology {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		// check what capabilities each DON has and register them with Capabilities Registry contract
		if flags.HasFlag(donTopology.Flags, types.CronCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "cron-trigger",
					Version:        "1.0.0",
					CapabilityType: 0, // TRIGGER
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		if flags.HasFlag(donTopology.Flags, types.CustomComputeCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "custom-compute",
					Version:        "1.0.0",
					CapabilityType: 1, // ACTION
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		if flags.HasFlag(donTopology.Flags, types.OCR3Capability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "offchain_reporting",
					Version:        "1.0.0",
					CapabilityType: 2, // CONSENSUS
					ResponseType:   0, // REPORT
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		if flags.HasFlag(donTopology.Flags, types.WriteEVMCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "write_geth-testnet",
					Version:        "1.0.0",
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		// Add support for new capabilities here as needed

		donPeerIDs := make([]string, len(donTopology.DON.Nodes)-1)
		for i, node := range donTopology.DON.Nodes {
			if i == 0 {
				continue
			}

			p2pID, err := keystonenode.ToP2PID(node, keystonenode.NoOpTransformFn)
			require.NoError(t, err, "failed to get p2p id for node %s", node.Name)

			donPeerIDs[i-1] = p2pID
		}

		// we only need to assign P2P IDs to NOPs, since `ConfigureInitialContractsChangeset` method
		// will take care of creating DON to Nodes mapping
		nop := keystone_changeset.NOP{
			Name:  fmt.Sprintf("NOP for %s DON", donTopology.NodeOutput.NodeSetName),
			Nodes: donPeerIDs,
		}

		donName := donTopology.NodeOutput.NodeSetName + "-don"
		donCapabilities = append(donCapabilities, keystone_changeset.DonCapabilities{
			Name:         donName,
			F:            1,
			Nops:         []keystone_changeset.NOP{nop},
			Capabilities: capabilities,
		})
	}

	var transmissionSchedule []int

	for _, donTopology := range keystoneEnv.DONTopology {
		if flags.HasFlag(donTopology.Flags, types.OCR3Capability) {
			// this schedule makes sure that all worker nodes are transmitting OCR3 reports
			transmissionSchedule = []int{len(donTopology.DON.Nodes) - 1}
			break
		}
	}

	require.NotEmpty(t, transmissionSchedule, "transmission schedule must not be empty")

	// values supplied by Alexandr Yepishev as the expected values for OCR3 config
	oracleConfig := keystone_changeset.OracleConfig{
		DeltaProgressMillis:               5000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		TransmissionSchedule:              transmissionSchedule,
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationAcceptMillis:           1000,
		MaxDurationTransmitMillis:         1000,
		MaxFaultyOracles:                  1,
		MaxQueryLengthBytes:               1000000,
		MaxObservationLengthBytes:         1000000,
		MaxReportLengthBytes:              1000000,
		MaxRequestBatchSize:               1000,
		UniqueReports:                     true,
	}

	cfg := keystone_changeset.InitialContractsCfg{
		RegistryChainSel: keystoneEnv.ChainSelector,
		Dons:             donCapabilities,
		OCR3Config:       &oracleConfig,
	}

	_, err := keystone_changeset.ConfigureInitialContractsChangeset(*keystoneEnv.Environment, cfg)
	require.NoError(t, err, "failed to configure initial contracts")
}

func DeployKeystone(t *testing.T, testLogger zerolog.Logger, keystoneEnv *types.KeystoneEnvironment) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.Environment, "environment must be set")
	require.NotEmpty(t, keystoneEnv.ChainSelector, "chain selector must be set")

	keystoneEnv.KeystoneContractAddresses = &types.KeystoneContractAddresses{}

	keystoneEnv.KeystoneContractAddresses.ForwarderAddress = deployKeystoneForwarder(t, testLogger, keystoneEnv.Environment, keystoneEnv.ChainSelector)
	keystoneEnv.KeystoneContractAddresses.OCR3CapabilityAddress = deployOCR3(t, testLogger, keystoneEnv.Environment, keystoneEnv.ChainSelector)
	keystoneEnv.KeystoneContractAddresses.CapabilitiesRegistryAddress = deployCapabilitiesRegistry(t, testLogger, keystoneEnv.Environment, keystoneEnv.ChainSelector)
	keystoneEnv.KeystoneContractAddresses.WorkflowRegistryAddress = deployWorkflowRegistry(t, testLogger, keystoneEnv.Environment, keystoneEnv.ChainSelector)
}

func deployOCR3(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployOCR3(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy OCR3 Capability contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var ocr3capabilityAddr common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "OCR3Capability") {
			ocr3capabilityAddr = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed OCR3Capability contract at %s", ocr3capabilityAddr.Hex())
			break
		}
	}
	require.NotEmpty(t, ocr3capabilityAddr, "failed to find OCR3Capability address in the address book")

	return ocr3capabilityAddr
}

func deployCapabilitiesRegistry(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployCapabilityRegistry(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy Capabilities Registry contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var capabilitiesRegistryAddr common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "CapabilitiesRegistry") {
			capabilitiesRegistryAddr = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed Capabilities Registry contract at %s", capabilitiesRegistryAddr.Hex())
			break
		}
	}
	require.NotEmpty(t, capabilitiesRegistryAddr, "failed to find Capabilities Registry address in the address book")

	return capabilitiesRegistryAddr
}

func deployKeystoneForwarder(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployForwarder(*ctfEnv, keystone_changeset.DeployForwarderRequest{
		ChainSelectors: []uint64{chainSelector},
	})
	require.NoError(t, err, "failed to deploy forwarder contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var forwarderAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "KeystoneForwarder") {
			forwarderAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed KeystoneForwarder contract at %s", forwarderAddress.Hex())
			break
		}
	}
	require.NotEmpty(t, forwarderAddress, "failed to find KeystoneForwarder address in the address book")

	return forwarderAddress
}

func deployWorkflowRegistry(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	require.NotNil(t, ctfEnv, "environment must not be nil")

	output, err := workflow_registry_changeset.Deploy(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy workflow registry contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var workflowRegistryAddr common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "WorkflowRegistry") {
			workflowRegistryAddr = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed WorkflowRegistry contract at %s", workflowRegistryAddr.Hex())
		}
	}
	require.NotEmpty(t, workflowRegistryAddr, "failed to find WorkflowRegistry address in the address book")

	return workflowRegistryAddr
}

func ConfigureWorkflowRegistry(t *testing.T, testLogger zerolog.Logger, keystoneEnv *types.KeystoneEnvironment) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.Environment, "environment must be set")
	require.NotEmpty(t, keystoneEnv.ChainSelector, "chain selector must be set")
	require.NotNil(t, keystoneEnv.SethClient, "seth client must be set")
	require.NotEmpty(t, keystoneEnv.WorkflowDONID, "workflow DON ID must be set")

	_, err := workflow_registry_changeset.UpdateAllowedDons(*keystoneEnv.Environment, &workflow_registry_changeset.UpdateAllowedDonsRequest{
		RegistryChainSel: keystoneEnv.ChainSelector,
		DonIDs:           []uint32{keystoneEnv.WorkflowDONID},
		Allowed:          true,
	})
	require.NoError(t, err, "failed to update allowed Dons")

	_, err = workflow_registry_changeset.UpdateAuthorizedAddresses(*keystoneEnv.Environment, &workflow_registry_changeset.UpdateAuthorizedAddressesRequest{
		RegistryChainSel: keystoneEnv.ChainSelector,
		Addresses:        []string{keystoneEnv.SethClient.MustGetRootKeyAddress().Hex()},
		Allowed:          true,
	})
	require.NoError(t, err, "failed to update authorized addresses")
}

func DeployFeedsConsumer(t *testing.T, testLogger zerolog.Logger, keystoneEnv *types.KeystoneEnvironment) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.Environment, "environment must be set")
	require.NotEmpty(t, keystoneEnv.ChainSelector, "chain selector must be set")
	require.NotNil(t, keystoneEnv.SethClient, "seth client must be set")
	require.NotNil(t, keystoneEnv.KeystoneContractAddresses, "keystone contract addresses must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.ForwarderAddress, "forwarder address must be set")

	output, err := keystone_changeset.DeployFeedsConsumer(*keystoneEnv.Environment, &keystone_changeset.DeployFeedsConsumerRequest{
		ChainSelector: keystoneEnv.ChainSelector,
	})
	require.NoError(t, err, "failed to deploy feeds_consumer contract")

	err = keystoneEnv.Environment.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := keystoneEnv.Environment.ExistingAddresses.AddressesForChain(keystoneEnv.ChainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", keystoneEnv.ChainSelector)

	var feedsConsumerAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "FeedConsumer") {
			feedsConsumerAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed FeedConsumer contract at %s", feedsConsumerAddress.Hex())
			break
		}
	}

	require.NotEmpty(t, feedsConsumerAddress, "failed to find FeedConsumer address in the address book")
	keystoneEnv.KeystoneContractAddresses.FeedsConsumerAddress = feedsConsumerAddress
}

func ConfigureFeedsConsumer(t *testing.T, testLogger zerolog.Logger, workflowName string, keystoneEnv *types.KeystoneEnvironment) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.Environment, "environment must be set")
	require.NotEmpty(t, keystoneEnv.ChainSelector, "chain selector must be set")
	require.NotNil(t, keystoneEnv.SethClient, "seth client must be set")
	require.NotNil(t, keystoneEnv.KeystoneContractAddresses, "keystone contract addresses must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.ForwarderAddress, "forwarder address must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.FeedsConsumerAddress, "feeds consumer address must be set")

	// configure Keystone Feeds Consumer contract, so it can accept reports from the forwarder contract,
	// that come from our workflow that is owned by the root private key
	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(keystoneEnv.KeystoneContractAddresses.FeedsConsumerAddress, keystoneEnv.SethClient.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	// Prepare hex-encoded and truncated workflow name
	var workflowNameBytes [10]byte
	var HashTruncateName = func(name string) string {
		// Compute SHA-256 hash of the input string
		hash := sha256.Sum256([]byte(name))

		// Encode as hex to ensure UTF8
		var hashBytes []byte = hash[:]
		resultHex := hex.EncodeToString(hashBytes)

		// Truncate to 10 bytes
		truncated := []byte(resultHex)[:10]
		return string(truncated)
	}

	truncated := HashTruncateName(workflowName)
	copy(workflowNameBytes[:], []byte(truncated))

	_, decodeErr := keystoneEnv.SethClient.Decode(feedsConsumerInstance.SetConfig(
		keystoneEnv.SethClient.NewTXOpts(),
		[]common.Address{keystoneEnv.KeystoneContractAddresses.ForwarderAddress}, // allowed senders
		[]common.Address{keystoneEnv.SethClient.MustGetRootKeyAddress()},         // allowed workflow owners
		// here we need to use hex-encoded workflow name converted to []byte
		[][10]byte{workflowNameBytes}, // allowed workflow names
	))
	require.NoError(t, decodeErr, "failed to set config for feeds consumer")
}
