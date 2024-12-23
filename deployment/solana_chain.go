package deployment

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// SolChain represents a Solana chain.
type SolChain struct {
	// Selectors used as canonical chain identifier.
	Selector uint64
	Client   SolClient
	Confirm  func() (uint64, error)
}

func (c SolChain) String() string {
	chainInfo, err := ChainInfo(c.Selector)
	if err != nil {
		// we should never get here, if the selector is invalid it should not be in the environment
		panic(err)
	}
	return fmt.Sprintf("%s (%d)", chainInfo.ChainName, chainInfo.ChainSelector)
}

func (c SolChain) Name() string {
	chainInfo, err := ChainInfo(c.Selector)
	if err != nil {
		// we should never get here, if the selector is invalid it should not be in the environment
		panic(err)
	}
	if chainInfo.ChainName == "" {
		return fmt.Sprintf("%d", c.Selector)
	}
	return chainInfo.ChainName
}

type SolClient interface {
}

type ContractDeploySolana struct {
	ProgramID *solana.PublicKey // We leave this incase a Go binding doesn't have Address()
	Tv        TypeAndVersion
	Err       error
}

func DeploySolContract(
	lggr logger.Logger,
	chain SolChain,
	addressBook AddressBook,
	deploy func(chain SolChain) ContractDeploySolana,
) (*ContractDeploySolana, error) {
	contractDeploy := deploy(chain)
	if contractDeploy.Err != nil {
		lggr.Errorw("Failed to deploy contract", "chain", chain.String(), "err", contractDeploy.Err)
		return nil, contractDeploy.Err
	}
	_, err := chain.Confirm()
	if err != nil {
		lggr.Errorw("Failed to confirm deployment", "chain", chain.String(), "Contract", contractDeploy.Tv.String(), "err", err)
		return nil, err
	}
	lggr.Infow("Deployed contract", "Contract", contractDeploy.Tv.String(), "addr", contractDeploy.ProgramID, "chain", chain.String())
	err = addressBook.Save(chain.Selector, "fill in address", contractDeploy.Tv)
	if err != nil {
		lggr.Errorw("Failed to save contract address", "Contract", contractDeploy.Tv.String(), "addr", contractDeploy.ProgramID, "chain", chain.String(), "err", err)
		return nil, err
	}
	return &contractDeploy, nil
}
