package changeset

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

type DeployLinkTokenConfig struct {
	// LinkTokenByChain is a mapping from chain selector to
	// the type of link token that should be deployed.
	LinkTokenByChain map[uint64]deployment.ContractType
}

var _ deployment.ChangeSet[DeployLinkTokenConfig] = DeployLinkToken

// DeployLinkToken deploys a link token contract to the chains identified by the config.
func DeployLinkToken(e deployment.Environment, config DeployLinkTokenConfig) (deployment.ChangesetOutput, error) {
	for chain, cType := range config.LinkTokenByChain {
		if _, ok := e.Chains[chain]; !ok {
			if !ok {
				return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
			}
		}
		if cType != types.StaticLinkToken && cType != types.BurnMintLinkToken {
			return deployment.ChangesetOutput{}, fmt.Errorf("invalid link token type")
		}
	}
	newAddresses := deployment.NewMemoryAddressBook()
	for chain, cType := range config.LinkTokenByChain {
		var err error
		switch cType {
		case types.StaticLinkToken:
			_, err = deployment.DeployContract[*link_token.LinkToken](e.Logger, e.Chains[chain], newAddresses,
				func(chain deployment.Chain) deployment.ContractDeploy[*link_token.LinkToken] {
					linkTokenAddr, tx, linkToken, err2 := link_token.DeployLinkToken(
						chain.DeployerKey,
						chain.Client,
					)
					return deployment.ContractDeploy[*link_token.LinkToken]{
						Address:  linkTokenAddr,
						Contract: linkToken,
						Tx:       tx,
						Tv:       deployment.NewTypeAndVersion(types.StaticLinkToken, deployment.Version1_0_0),
						Err:      err2,
					}
				})
		case types.BurnMintLinkToken:
			_, err = deployment.DeployContract(e.Logger, e.Chains[chain], newAddresses,
				func(chain deployment.Chain) deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
					linkTokenAddr, tx2, linkToken, err2 := burn_mint_erc677.DeployBurnMintERC677(
						chain.DeployerKey,
						chain.Client,
						"Link Token",
						"LINK",
						uint8(18),
						big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
					)
					return deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
						linkTokenAddr, linkToken, tx2, deployment.NewTypeAndVersion(types.BurnMintLinkToken, deployment.Version1_0_0), err2,
					}
				})
		default:
			return deployment.ChangesetOutput{}, fmt.Errorf("impossible already validated type")
		}
		if err != nil {
			e.Logger.Errorw("Failed to deploy link token", "err", err)
			return deployment.ChangesetOutput{AddressBook: newAddresses}, err
		}
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}
