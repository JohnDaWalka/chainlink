package changeset

import (
	"math/big"
	"time"

	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
)

// MCMSConfig defines timelock duration.
type MCMSConfig struct {
	MinDelay time.Duration
}

// MakeTxOptsAndHandlerForContract creates transaction opts and a handler that either submits and confirms a transaction or outputs it as an MCMS operation,
// depending on whether or not an MCMS config is supplied.
func MakeTxOptsAndHandlerForContract(
	contractAddress common.Address,
	chain deployment.Chain,
	mcmsConfig *MCMSConfig,
) (*bind.TransactOpts, func(tx *types.Transaction, err error) (*mcms.Operation, error), error) {
	// Set the transaction opts based whether or not the MCMS config exists
	var opts *bind.TransactOpts
	if mcmsConfig != nil {
		opts = deployment.SimTransactOpts()
	} else {
		opts = chain.DeployerKey
	}

	// Create handler
	handler := func(tx *types.Transaction, err error) (*mcms.Operation, error) {
		if opts.From == chain.DeployerKey.From {
			if _, err = deployment.ConfirmIfNoError(chain, tx, err); err != nil {
				return nil, fmt.Errorf("failed to confirm transaction with hash %s on %s: %w", tx.Hash(), chain.String(), err)
			}
			return nil, nil
		}
		return &mcms.Operation{
			To:    contractAddress,
			Data:  tx.Data(),
			Value: big.NewInt(0),
		}, nil
	}

	return opts, handler, nil
}
