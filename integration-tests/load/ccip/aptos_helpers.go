package ccip

import (
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func fundAdditionalAptosKeys(
	t *testing.T,
	signer aptos.TransactionSigner,
	e cldf.Environment,
	destChains []uint64,
	fundingAmount uint64,
) (map[uint64][]aptos.Account, error) {
	funded := make(map[uint64][]aptos.Account, len(destChains))

	for _, chain := range e.BlockChains.AptosChains() {
		numAccounts := len(destChains)
		funded[chain.ChainSelector()] = make([]aptos.Account, 0, numAccounts)

		for range numAccounts {
			account, err := aptos.NewEd25519Account()

			memory.FundAptosAccount(t, signer, account.AccountAddress(), fundingAmount, chain.Client)
			if err != nil {
				return map[uint64][]aptos.Account{}, err
			}

			funded[chain.ChainSelector()] = append(funded[chain.ChainSelector()], *account)
		}
	}
	return funded, nil
}
