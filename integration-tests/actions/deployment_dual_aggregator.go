package actions

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	gethwrappers "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/dual-transmission"
)

type Input struct {
	URL string  `toml:"url"`
	Out *Output `toml:"out"`
}

type Output struct {
	UseCache  bool             `toml:"use_cache"`
	Addresses []common.Address `toml:"addresses"`
}

func NewDualAggregatorDeployment(c *seth.Client, in *Input, linkContractAddress string, offchainOptions contracts.OffchainOptions) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	abi, err := gethwrappers.DualAggregatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	dd, err := c.DeployContract(c.NewTXOpts(),
		"DualAggregator",
		*abi,
		common.FromHex(gethwrappers.DualAggregatorMetaData.Bin),
		common.HexToAddress(linkContractAddress),
		offchainOptions.MinimumAnswer,
		offchainOptions.MaximumAnswer,
		offchainOptions.BillingAccessController,   // billingAccessController
		offchainOptions.RequesterAccessController, // requesterAccessController
		offchainOptions.Decimals,
		offchainOptions.Description,
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // secondary proxy
		uint32(30), // cutOffTime
		uint32(20), // maxSyncIterations
	)
	if err != nil {
		return nil, err
	}
	out := &Output{
		UseCache: true,
		// save all the addresses to output, so it can be cached
		Addresses: []common.Address{dd.Address},
	}
	in.Out = out
	return out, nil
}
