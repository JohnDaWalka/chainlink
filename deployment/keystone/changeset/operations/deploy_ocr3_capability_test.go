package operations_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	keystoneops "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/jobs"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func TestDeployOCR3CapabilitySeq(t *testing.T) {
	t.Parallel()

	t.Run("success - deploy and configure OCR3 capability", func(t *testing.T) {
		// Setup test environment with contract deployment
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		// Get the capabilities registry address
		capRegistryAddr := te.CapabilityRegistryAddressRef()
		require.NotNil(t, capRegistryAddr)

		// Create test DON capabilities
		donCapabilities := []internal.DonCapabilities{
			{
				Name: "wfDon",
				Capabilities: []internal.DONCapabilityWithConfig{
					{
						Capability: kcr.CapabilitiesRegistryCapability{
							LabelledName:   "offchain-consensus",
							Version:        "1.0.2",
							CapabilityType: 2, // CONSENSUS
							ResponseType:   0,
						},
					},
				},
			},
		}

		// Get nodes from the test environment
		wfNodeIDsToP2PIDs := te.GetJDNodeIDsToP2PIDs("wfDon")
		require.NotEmpty(t, wfNodeIDsToP2PIDs)
		wfNodeNamesToP2PIDs := te.GetJDNodeNamesToP2PIDs("wfDon")
		require.NotEmpty(t, wfNodeNamesToP2PIDs)

		wfNodeIDs := make([]string, 0, len(wfNodeIDsToP2PIDs))
		seqNodes := make([]jobs.DistributeOCRJobSpecSeqNode, 0, len(wfNodeIDsToP2PIDs))
		for nodeID, p2pID := range wfNodeIDsToP2PIDs {
			seqNodes = append(seqNodes, jobs.DistributeOCRJobSpecSeqNode{
				ID:       nodeID,
				P2PLabel: p2pID,
			})
			wfNodeIDs = append(wfNodeIDs, nodeID)
		}
		bootstrapNodes := make([]jobs.DistributeBootstrapJobSpecsSeqBootCfg, 0, len(wfNodeIDsToP2PIDs))
		for nodeName, p2pID := range wfNodeNamesToP2PIDs {
			bootstrapNodes = append(bootstrapNodes, jobs.DistributeBootstrapJobSpecsSeqBootCfg{
				NodeName: nodeName,
				P2PID:    p2pID,
			})
		}

		// Create oracle config
		oracleConfig := internal.OracleConfig{
			MaxFaultyOracles:     1,
			DeltaProgressMillis:  5000,
			TransmissionSchedule: []int{4}, // 4 nodes in the DON
		}

		// Create DON configuration for OCR3
		ocr3DON := contracts.ConfigureKeystoneDON{
			Name:    "wfDon",
			NodeIDs: wfNodeIDs,
		}

		donInfos, err := internal.DonInfos(donCapabilities, te.Env.Offchain)
		require.NoError(t, err)

		chain, ok := te.Env.BlockChains.EVMChains()[te.RegistrySelector]
		require.True(t, ok)

		capabilitiesRegistry, err := changeset.LoadCapabilityRegistry(chain, te.Env, capRegistryAddr)
		require.NoError(t, err)

		donToCapabilities, err := internal.MapDonsToCaps(capabilitiesRegistry.Contract, donInfos)
		require.NoError(t, err)

		var capabilities []kcr.CapabilitiesRegistryCapability
		for _, don := range donToCapabilities {
			for _, donCap := range don {
				capabilities = append(capabilities, donCap.CapabilitiesRegistryCapability)
			}
		}

		// Set up dependencies
		deps := keystoneops.DeployOCR3Capability{
			Env: &te.Env,
		}

		// Set up input
		input := keystoneops.DeployOCR3CapabilityInput{
			Nodes:                seqNodes,
			RegistryChainSel:     te.RegistrySelector,
			RegistryRef:          capRegistryAddr,
			OracleConfig:         oracleConfig,
			DONs:                 []contracts.ConfigureKeystoneDON{ocr3DON},
			DomainKey:            "keystone",
			EnvironmentLabel:     "test",
			DONName:              "wfDon",
			ChainSelectorEVM:     te.RegistrySelector,
			ChainSelectorAptos:   0, // Not using Aptos in this test
			BootstrapperOCR3Urls: []string{"12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"},
			BootstrapCfgs:        bootstrapNodes,
		}

		// Execute the sequence
		f := func() context.Context {
			return t.Context()
		}
		bundle := operations.NewBundle(f, te.Env.Logger, operations.NewMemoryReporter())
		report, err := operations.ExecuteSequence(bundle, keystoneops.DeployOCR3CapabilitySeq, deps, input)

		// Verify results
		require.NoError(t, err)
		require.NotNil(t, report)

		output := report.Output

		// Verify OCR3 contract was deployed
		require.NotNil(t, output.Addresses)
		addresses, err := output.Addresses.Fetch()
		require.NoError(t, err)
		require.NotEmpty(t, addresses, "OCR3 contract should be deployed")

		// Verify job specs were created and distributed
		require.NotEmpty(t, output.JobSpecs, "OCR3 job specs should be created")
		assert.Len(t, output.JobSpecs, len(wfNodeIDs), "Should have one job spec per node")

		// Verify bootstrap job spec was created
		require.NotEmpty(t, output.BootstrapSpec, "Bootstrap job spec should be created")

		// Verify no MCMS proposals were created (since not using MCMS)
		assert.Empty(t, output.MCMSTimelockProposals)
		assert.Nil(t, output.BatchOperation)

		// Verify the deployed OCR3 contract exists in the environment
		ocr3Address := common.HexToAddress(addresses[0].Address)
		assert.NotEqual(t, common.Address{}, ocr3Address, "OCR3 contract address should be valid")

		te.Env.Logger.Infow("DeployOCR3CapabilitySeq test completed successfully",
			"ocr3Address", ocr3Address.Hex(),
			"jobSpecsCount", len(output.JobSpecs),
		)
	})
	/*
		t.Run("success - with MCMS", func(t *testing.T) {
			// Setup test environment with MCMS enabled
			te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
				WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
				NumChains:       1,
				UseMCMS:         true,
			})

			// Get the capabilities registry address
			capRegistryAddr := te.CapabilityRegistryAddressRef()
			require.NotNil(t, capRegistryAddr)

			// Create minimal test setup
			donCapabilities := []internal.DonCapabilities{
				{
					Name: "wfDon",
					Capabilities: []internal.RegisteredCapability{
						{
							CapabilitiesRegistryCapability: te.CapabilitiesRegistry().CapabilitiesRegistryCapability{
								LabelledName:   "offchain-consensus",
								Version:        "1.0.0",
								CapabilityType: 2,
								ResponseType:   0,
							},
						},
					},
				},
			}

			wfNodes := te.GetNodes("wfDon")
			oracleConfig := internal.OracleConfig{
				MaxFaultyOracles:     1,
				DeltaProgressMillis:  5000,
				TransmissionSchedule: []int{4},
			}

			ocr3DON := contracts.ConfigureKeystoneDON{
				Name:    "wfDon",
				NodeIDs: te.GetP2PIDs("wfDon").Strings(),
			}

			deps := keystoneops.DeployOCR3Capability{
				Env:             &te.Env,
				Nodes:           wfNodes,
				DonCapabilities: donCapabilities,
			}

			input := keystoneops.DeployOCR3CapabilityInput{
				RegistryChainSel:        te.RegistrySelector,
				RegistryContractAddress: &capRegistryAddr.Address,
				OracleConfig:            oracleConfig,
				DONs:                    []contracts.ConfigureKeystoneDON{ocr3DON},
				MCMSConfig:              &changeset.MCMSConfig{MinDuration: 0},
				DomainKey:               "keystone",
				EnvironmentLabel:        "test",
				DONName:                 "wfDon",
				ChainSelectorEVM:        te.RegistrySelector,
				ChainSelectorAptos:      0,
				BootstrapperOCR3Urls:    []string{},
			}

			// Execute the sequence
			report, err := operations.ExecuteSequence(
				operations.Bundle{Logger: te.Env.Logger},
				keystoneops.DeployOCR3CapabilitySeq,
				deps,
				input,
			)

			// Verify results with MCMS
			require.NoError(t, err)
			require.NotNil(t, report)

			output := report.Output

			// Verify MCMS proposals were created
			assert.NotEmpty(t, output.MCMSTimelockProposals, "MCMS proposals should be created when using MCMS")

			// Verify other outputs
			require.NotNil(t, output.Addresses)
			addresses, err := output.Addresses.Fetch()
			require.NoError(t, err)
			require.NotEmpty(t, addresses)

			require.NotEmpty(t, output.JobSpecs)
			assert.Len(t, output.JobSpecs, len(wfNodes))
		})

		t.Run("failure - invalid registry chain selector", func(t *testing.T) {
			te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
				WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
				NumChains:       1,
			})

			capRegistryAddr := te.CapabilityRegistryAddressRef()
			deps := keystoneops.DeployOCR3Capability{
				Env:             &te.Env,
				Nodes:           te.GetNodes("wfDon"),
				DonCapabilities: []internal.DonCapabilities{},
			}

			input := keystoneops.DeployOCR3CapabilityInput{
				RegistryChainSel:        99999999, // Invalid chain selector
				RegistryContractAddress: &capRegistryAddr.Address,
				OracleConfig:            internal.OracleConfig{},
				DONs:                    []contracts.ConfigureKeystoneDON{},
				DomainKey:               "keystone",
				EnvironmentLabel:        "test",
				DONName:                 "wfDon",
				ChainSelectorEVM:        te.RegistrySelector,
				ChainSelectorAptos:      0,
				BootstrapperOCR3Urls:    []string{},
			}

			// Execute the sequence and expect failure
			_, err := operations.ExecuteSequence(
				operations.Bundle{Logger: te.Env.Logger},
				keystoneops.DeployOCR3CapabilitySeq,
				deps,
				input,
			)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "does not exist in environment")
		})

		t.Run("failure - invalid registry contract address", func(t *testing.T) {
			te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
				WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
				NumChains:       1,
			})

			invalidAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			deps := keystoneops.DeployOCR3Capability{
				Env:             &te.Env,
				Nodes:           te.GetNodes("wfDon"),
				DonCapabilities: []internal.DonCapabilities{},
			}

			input := keystoneops.DeployOCR3CapabilityInput{
				RegistryChainSel:        te.RegistrySelector,
				RegistryContractAddress: &invalidAddr,
				OracleConfig:            internal.OracleConfig{},
				DONs:                    []contracts.ConfigureKeystoneDON{},
				DomainKey:               "keystone",
				EnvironmentLabel:        "test",
				DONName:                 "wfDon",
				ChainSelectorEVM:        te.RegistrySelector,
				ChainSelectorAptos:      0,
				BootstrapperOCR3Urls:    []string{},
			}

			// Execute the sequence and expect failure
			_, err := operations.ExecuteSequence(
				operations.Bundle{Logger: te.Env.Logger},
				keystoneops.DeployOCR3CapabilitySeq,
				deps,
				input,
			)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed to get capabilities registry contract")
		})
	*/
}
