package vault

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/tdh2/go/tdh2/lib/group/nist"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault/sanmarinodkg/dummydkg"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault/sanmarinodkg/tdh2shim"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

const flag = cre.VaultCapability

func New(chainID uint64) (*capabilities.Capability, error) {
	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(jobSpec(chainID)),
		capabilities.WithGatewayJobHandlerConfigFn(handlerConfig),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
		capabilities.WithValidateFn(func(c *capabilities.Capability) error {
			if chainID == 0 {
				return fmt.Errorf("chainID is required, got %d", chainID)
			}
			return nil
		}),
	)
}

func EncryptSecret(secret, masterPublicKeyStr string) (string, error) {
	masterPublicKey := tdh2easy.PublicKey{}
	masterPublicKeyBytes, err := hex.DecodeString(masterPublicKeyStr)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode master public key")
	}
	err = masterPublicKey.Unmarshal(masterPublicKeyBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal master public key")
	}
	cipher, err := tdh2easy.Encrypt(&masterPublicKey, []byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt secret")
	}
	cipherBytes, err := cipher.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal encrypted secrets to bytes")
	}
	return hex.EncodeToString(cipherBytes), nil
}

func jobSpec(chainID uint64) cre.JobSpecFn {
	return func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
		if input.DonTopology == nil {
			return nil, errors.New("topology is nil")
		}
		donToJobSpecs := make(cre.DonsToJobSpecs)

		// return early if no DON has the vault capability
		if !don.AnyDonHasCapability(input.DonTopology.Dons.DonMetadata, flag) {
			return donToJobSpecs, nil
		}

		vaultOCR3Key := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"capability_vault_plugin",
		)
		vaultCapabilityAddress, err := input.CldEnvironment.DataStore.Addresses().Get(vaultOCR3Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get Vault capability address")
		}

		dkgKey := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"capability_vault_dkg",
		)
		dkgAddress, err := input.CldEnvironment.DataStore.Addresses().Get(dkgKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get DKG address")
		}

		for _, donWithMetadata := range input.DonTopology.Dons.DonMetadata {
			if !flags.HasFlag(donWithMetadata.Flags, flag) {
				continue
			}

			// create job specs for the worker nodes
			workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: cre.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
			if err != nil {
				// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
				return nil, errors.Wrap(err, "failed to find worker nodes")
			}

			bootstrapNode, bootErr := input.DonTopology.BootstrapNode()
			if bootErr != nil {
				return nil, errors.Wrap(bootErr, "failed to find bootstrap node")
			}

			bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
			}

			_, ocrPeeringData, peeringErr := cre.PeeringCfgs(bootstrapNode)
			if peeringErr != nil {
				return nil, errors.Wrap(peeringErr, "failed to find peering data")
			}

			// create job specs for the bootstrap node
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "vault-capability", vaultCapabilityAddress.Address, chainID))

			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				ethKey, ok := workerNode.Keys.EVM[chainID]
				if !ok {
					return nil, fmt.Errorf("node %s does not have EVM key for chainID %d", nodeID, chainID)
				}

				ocr2KeyBundlesPerFamily, ocr2kbErr := node.ExtractBundleKeysPerFamily(workerNode)
				if ocr2kbErr != nil {
					return nil, errors.Wrap(ocr2kbErr, "failed to get ocr2 key bundle id from labels")
				}

				evmOCR2KeyBundle, ok := ocr2KeyBundlesPerFamily["evm"]
				if !ok {
					return nil, errors.New("key bundle ID for evm family is not found")
				}

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerVaultOCR3(nodeID, vaultCapabilityAddress.Address, dkgAddress.Address, ethKey.PublicAddress.Hex(), evmOCR2KeyBundle, ocrPeeringData, chainID))
			}
		}

		return donToJobSpecs, nil
	}
}

func handlerConfig(donMetadata *cre.DonMetadata) (cre.HandlerTypeToConfig, error) {
	if !flags.HasFlag(donMetadata.Flags, flag) {
		return nil, nil
	}

	return map[string]string{coregateway.VaultHandlerType: `
ServiceName = "vault"
[gatewayConfig.Dons.Handlers.Config]
requestTimeoutSec = 70
[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
`}, nil
}

func encryptPrivateShare(offchain cldf_offchain.Client, nodeID string, sk *tdh2easy.PrivateShare) (string, error) {
	nodeResp, err := offchain.GetNode(context.Background(), &nodev1.GetNodeRequest{
		Id: nodeID,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to get node from jd")
	}
	wk := nodeResp.GetNode().GetWorkflowKey()
	if wk == "" {
		return "", errors.New("node must contain a workflow key")
	}

	wkb, err := hex.DecodeString(wk)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode workflow key from hex")
	}

	skb, err := sk.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal private share")
	}

	wkbSized := [32]byte(wkb)
	sealed, err := box.SealAnonymous(nil, skb, &wkbSized, cryptorand.Reader)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt private share")
	}

	return hex.EncodeToString(sealed), nil
}

// satisfy unused lint
var _ = encryptPrivateShare

func dkgKeys(n, t int) (string, []*tdh2easy.PrivateShare, error) {
	instanceID, recipCfg, recipSecKeys, err := dummydkg.NewDKGSetup(n, t, "REPLACE_ME_WITH_RANDOM_SEED")
	if err != nil {
		return "", nil, err
	}

	group := nist.NewP256()
	result, err := dummydkg.NewDKGResult(instanceID, recipCfg, group)
	if err != nil {
		return "", nil, err
	}

	shares := []*tdh2easy.PrivateShare{}
	for _, share := range recipSecKeys {
		s, ierr := tdh2shim.TDH2PrivateShareFromDKGResult(result, share)
		if ierr != nil {
			return "", nil, errors.Wrap(ierr, "failed to convert DKG share to TDH2 share")
		}

		shares = append(shares, s)
	}

	pk, err := tdh2shim.TDH2PublicKeyFromDKGResult(result)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to convert DKG result to TDH2 public key")
	}

	pkb, err := pk.Marshal()
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to marshal TDH2 public key")
	}

	pks := hex.EncodeToString(pkb)
	framework.L.Info().Msg("Generated MasterPublicKey for n=" + strconv.Itoa(n) + ", t=" + strconv.Itoa(t) + ". Key = " + pks)
	return pks, shares, nil
}

// satisfy unused lint
var _ = dkgKeys

func registerWithV1(donFlags []string, _ *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if flags.HasFlag(donFlags, flag) {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "vault",
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}
