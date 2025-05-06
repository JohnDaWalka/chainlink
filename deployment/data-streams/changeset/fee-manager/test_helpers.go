package fee_manager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"
	"github.com/stretchr/testify/require"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

type DataStreamsTestEnvSetupOutput struct {
	Env               deployment.Environment
	LinkTokenAddress  common.Address
	FeeManagerAddress common.Address
}

type DataStreamsTestEnvOptions struct {
	DeployFeeManager bool
	DeployLinkToken  bool
	DeployMCMS       bool
}

func NewDefaultOptions() DataStreamsTestEnvOptions {
	return DataStreamsTestEnvOptions{
		DeployLinkToken:  true,
		DeployMCMS:       true,
		DeployFeeManager: true,
	}
}

func DeployTestEnvironment(t *testing.T, opts DataStreamsTestEnvOptions) (DataStreamsTestEnvSetupOutput, error) {
	t.Helper()

	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      opts.DeployMCMS,
		ShouldDeployLinkToken: opts.DeployLinkToken,
	})
	e := testEnv.Environment

	feeManagerAddress := common.HexToAddress("0x044304C47eD3B1C1357569960A537056AFE8c815")

	if opts.DeployFeeManager {
		// FM checks LinkToken is ERC20 - but accepts any address for NativeTokenAddress
		cc := DeployFeeManagerConfig{
			ChainsToDeploy: map[uint64]DeployFeeManager{testutil.TestChain.Selector: {
				LinkTokenAddress:     testEnv.LinkTokenState.LinkToken.Address(),
				NativeTokenAddress:   common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
				VerifierProxyAddress: common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e"),
				RewardManagerAddress: common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2"),
			}},
		}

		env, err := commonchangesets.Apply(t, e, nil,
			commonchangesets.Configure(
				DeployFeeManagerChangeset,
				cc,
			),
		)
		require.NoError(t, err)

		record, err := env.DataStore.Addresses().Get(ds.NewAddressRefKey(testutil.TestChain.Selector, ds.ContractType(types.FeeManager), &deployment.Version0_5_0, ""))
		require.NoError(t, err)
		feeManagerAddress = common.HexToAddress(record.Address)
		e = env
	}

	return DataStreamsTestEnvSetupOutput{
		Env:               e,
		LinkTokenAddress:  testEnv.LinkTokenState.LinkToken.Address(),
		FeeManagerAddress: feeManagerAddress,
	}, nil
}

func GetState(
	t *testing.T,
	inDs ds.MutableDataStore[ds.DefaultMetadata, ds.DefaultMetadata],
	chainSelector uint64,
	contractAddress common.Address,
) *metadata.GenericContractMetadata[v0_5.FeeManagerView] {
	envDatastore, err := ds.FromDefault[metadata.SerializedContractMetadata, ds.DefaultMetadata](inDs.Seal())
	require.NoError(t, err)

	cm, err := envDatastore.ContractMetadata().Get(
		ds.NewContractMetadataKey(chainSelector, contractAddress.String()),
	)
	require.NoError(t, err, "Failed to get contract metadata")

	contractMetadata, err := metadata.DeserializeMetadata[v0_5.FeeManagerView](cm.Metadata)
	require.NoError(t, err, "Failed to convert contract metadata")
	require.NotNil(t, contractMetadata, "Failed to get contract metadata")
	return contractMetadata
}
