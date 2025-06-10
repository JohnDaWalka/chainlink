package sui

import (
	"encoding/hex"
	"fmt"
	"strings"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	offrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_offramp"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[v1_6.SetOCR3OffRampConfig] = SetOCR3Offramp{}

type SetOCR3Offramp struct{}

// Apply implements deployment.ChangeSetV2.
func (s SetOCR3Offramp) Apply(e cldf.Environment, config v1_6.SetOCR3OffRampConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Sui onchain state: %w", err)
	}

	ab := cldf.NewMemoryAddressBook()

	for _, remoteSelector := range config.RemoteChainSels {
		suiChains := e.BlockChains.SuiChains()
		suiChain := suiChains[remoteSelector]
		suiSigner := rel.NewPrivateKeySigner(suiChain.DeployerKey)

		deps := SuiDeps{
			AB: ab,
			SuiChain: sui_ops.OpTxDeps{
				Client: *suiChain.Client,
				Signer: suiSigner,
				GetTxOpts: func() bind.TxOpts {
					b := uint64(300_000_000)
					return bind.TxOpts{
						GasBudget: &b,
					}
				},
			},
			CCIPOnChainState: state,
		}

		// 4 random signer addresses
		signerAddresses := []string{
			"0x3f6d6a9e3f7707485bf51c02a6bc6cb6e17dffe7f3e160b3c5520d55d1de8398",
			"0xdbc00fee80d2f1d061ea2c2227247bda49def7fa6c389d313b59bb59ffa73de1",
			"0x776b1e5d91dfecff55c18e949463395a53913867b26d471eea0b7fc382b431cf",
			"0x96c55059838728052c0d2466f5ad702bb6cbb6f01a0dd98939fe7604ecfee92c",
		}

		signerAddrBytes := make([][]byte, 0, len(signerAddresses))

		for _, addr := range signerAddresses {
			addrHex := strings.TrimPrefix(addr, "0x")

			addrBytes, err := hex.DecodeString(addrHex)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to decode address %s: %w", addr, err)
			}

			signerAddrBytes = append(signerAddrBytes, addrBytes)
		}

		// TODO: THIS INPUT NEEDS TO BE ACC
		setOCR3ConfigInput := offrampops.SetOCR3ConfigInput{
			OffRampPackageId: state.SuiChains[remoteSelector].OffRampAddress.String(),
			OffRampStateId:   state.SuiChains[remoteSelector].OffRampStateObjectId.String(),
			OwnerCapObjectId: state.SuiChains[remoteSelector].OffRampOwnerCapId.String(),
			// Sample config digest
			ConfigDigest: []byte{
				0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
				0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
				0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
				0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
			},
			OCRPluginType:                  byte(0),
			BigF:                           byte(1),
			IsSignatureVerificationEnabled: true,
			Signers:                        signerAddrBytes,
			Transmitters:                   signerAddresses,
		}

		_, err := operations.ExecuteOperation(e.OperationsBundle, offrampops.SetOCR3ConfigOp, deps.SuiChain, setOCR3ConfigInput)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (s SetOCR3Offramp) VerifyPreconditions(e cldf.Environment, config v1_6.SetOCR3OffRampConfig) error {
	return nil
}
