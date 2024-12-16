package smoke

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	gethwrappers "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/dual-transmission"
)

type Cfg struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	Contracts   *actions.Input    `toml:"contracts" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
}

var bootstrapJobSpec = `
type = "bootstrap"
schemaVersion = 1
name = "smoke OEV bootstrap"
externalJobID = "%s"
contractID = "%s"
relay = "%s"

[relayConfig]
chainID = %s
`
var oevJobSpec = `
type = "offchainreporting2"
schemaVersion = 1
name = "OEV job"
externalJobID = "%s"
forwardingAllowed = true
maxTaskDuration = "0s"
contractID = "%s"
relay = "%s"
pluginType = "median"
transmitterID = "%s"
observationSource = """
 //randomness
    val1 [type="memo" value="10"]
    val2 [type="memo" value="20"]
    val3 [type="memo" value="30"]
    val4 [type="memo" value="40"]
    val5 [type="memo" value="50"]
    val6 [type="memo" value="60"]
    val7 [type="memo" value="70"]
    val8 [type="memo" value="80"]
    val9 [type="memo" value="90"]

    random1 [type="any"]
    random2 [type="any"]
    random3 [type="any"]

    val1 -> random1
    val2 -> random2
    val3 -> random3
    val4 -> random1
    val5 -> random2
    val6 -> random3
    val7 -> random1
    val8 -> random2
    val9 -> random3


    // data source 1
    ds1_multiply [type="multiply" times=100]

     // data source 2
    ds2_multiply [type="multiply" times=100]


    // data source 3
    ds3_multiply [type="multiply" times=100]


    random1 -> ds1_multiply -> answer
    random2 -> ds2_multiply -> answer
    random3 -> ds3_multiply -> answer

    answer [type=median]
"""

[relayConfig]
chainID = %s
enableDualTransmission = true

[relayConfig.dualTransmission]
contractAddress = "%s"
transmitterAddress = "%s"

[relayConfig.dualTransmission.meta]
hint = [ "calldata" ]
refund = [ "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90" ]

[pluginConfig]
juelsPerFeeCoinSource = """
juels_per_fee_coin [type="sum" values=<[0]>];
"""
`

func TestSmoke(t *testing.T) {
	lggr := logging.GetTestLogger(t)
	pkey := os.Getenv("PRIVATE_KEY")
	require.NotEmpty(t, pkey, "private key is empty")
	in, err := framework.Load[Cfg](t)
	require.NoError(t, err)

	//1. Set up nodes & docker environment (5 nodes)
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	//Replace customURL
	in.NodeSet.NodeSpecs[0].Node.UserConfigOverrides = fmt.Sprintf(in.NodeSet.NodeSpecs[0].Node.UserConfigOverrides, bc.Nodes[0].DockerInternalHTTPUrl)

	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err)

	// connecting clients
	sethClient, err := seth.NewClientBuilder().
		WithGasPriceEstimations(true, 0, seth.Priority_Fast).
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		Build()
	require.NoError(t, err)
	nodeClients, err := clclient.New(out.CLNodes)
	require.NoError(t, err)
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	//2. Create secondary addresses
	primaryAddresses := make([]common.Address, 0)
	secondaryAddresses := make([]common.Address, 0)
	//Create secondary ETH key
	for i := range workerNodes {
		primary, err := workerNodes[i].PrimaryEthAddress()
		require.NoError(t, err, "Cannot get primary key")
		primaryAddresses = append(primaryAddresses, common.HexToAddress(primary))

		key, _, err := workerNodes[i].CreateTxKey("evm", in.BlockchainA.ChainID) //TODO: @george-dorin Remove hardcoded evm
		require.NoError(t, err, "Cannot create secondary key")
		secondaryAddresses = append(secondaryAddresses, common.HexToAddress(key.Data.Attributes.Address))
	}

	//3. Restart
	out, err = ns.UpgradeNodeSet(in.NodeSet, bc, time.Second*10)
	require.NoError(t, err, "Cannot restart nodes")
	nodeClients, err = clclient.New(out.CLNodes)
	require.NoError(t, err, "Cannot create clients")
	bootstrapNode, workerNodes = nodeClients[0], nodeClients[1:]

	//4. Fund addresses
	err = ns.FundNodes(sethClient.Client, nodeClients, pkey, 0.2)
	require.NoError(t, err)
	for i := range secondaryAddresses {
		err = ns.SendETH(sethClient.Client, pkey, secondaryAddresses[i].String(), big.NewFloat(0.2))
		require.NoError(t, err, "Cannot fund secondary address")
	}

	//5. Deploy link contract
	linkContract, err := contracts.DeployLinkTokenContract(lggr, sethClient)
	require.NoError(t, err, "Error loading/deploying link token contract")

	//6. Deploy forwarders
	var operators []common.Address
	operators, forwarders, _ := actions.DeployForwarderContracts(
		t, sethClient, common.HexToAddress(linkContract.Address()), len(workerNodes),
	)
	require.Equal(t, len(workerNodes), len(operators), "Number of operators should match number of nodes")
	require.Equal(t, len(workerNodes), len(forwarders), "Number of authorized forwarders should match number of nodes")

	//7. Configure forwarders
	require.NoError(t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	for i := range workerNodes {
		actions.AcceptAuthorizedReceiversOperator(
			t, lggr, sethClient, operators[i], forwarders[i], []common.Address{primaryAddresses[i], secondaryAddresses[i]})
		require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")

		chainIDBigInt := new(big.Int)
		chainIDBigInt, ok := chainIDBigInt.SetString(in.BlockchainA.ChainID, 10)
		require.True(t, ok, "ChainID cannot be converted to big.Int")
		_, _, err = workerNodes[i].TrackForwarder(chainIDBigInt, forwarders[i])
		require.NoError(t, err, "Cannot track forwarders")
	}

	//8. Deploy dual agg
	//in.Contracts.URL = bc.Nodes[0].HostWSUrl
	ocrOffchainOptions := contracts.DefaultOffChainAggregatorOptions()
	oevContract, err := actions.NewDualAggregatorDeployment(sethClient, in.Contracts, linkContract.Address(), ocrOffchainOptions)
	require.NoError(t, err)
	_, err = gethwrappers.NewDualAggregator(oevContract.Addresses[0], sethClient.Client)
	require.NoError(t, err)

	//9. Create jobs
	//Bootstrap
	response, _, err2 := bootstrapNode.CreateJobRaw(fmt.Sprintf(bootstrapJobSpec, uuid.New().String(), oevContract.Addresses[0].String(), "evm", in.BlockchainA.ChainID))
	require.NoError(t, err2)
	require.Empty(t, response.Errors)

	//Feed job
	for i := range workerNodes {
		response, _, err2 = bootstrapNode.CreateJobRaw(fmt.Sprintf(oevJobSpec, uuid.New().String(), oevContract.Addresses[0].String(), "evm", primaryAddresses[i], in.BlockchainA.ChainID, oevContract.Addresses[0].String(), secondaryAddresses[i]))
		require.NoError(t, err2)
		require.Empty(t, response.Errors)
	}

	//10. Configure dual agg
	//ocrOffchainOptions := contracts2.DefaultOffChainAggregatorOptions()
	//actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	//dualAggContract.SetConfig()

	t.Run("test OEV", func(t *testing.T) {
		//dualAggContract.DualAggregatorFilterer.FilterSecondaryRoundIdUpdated()

		//Check round ID and check if we have transmitSecondary event
	})
}
