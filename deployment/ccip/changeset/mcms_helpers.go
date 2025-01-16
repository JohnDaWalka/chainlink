package changeset

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/chainlink/deployment"
)

type MCMSConfig struct {
	MinDelay time.Duration
}

// makeTxOptsAndHandlerForContract creates transaction opts and a handler that either confirms a transaction or adds it to an MCMS proposal,
// all depending on the current owner of the target contract.
func makeTxOptsAndHandlerForContract(
	contractAddress common.Address,
	chain deployment.Chain,
	deployerKey *bind.TransactOpts,
	mcmsConfig *MCMSConfig,
) (*bind.TransactOpts, func(tx *types.Transaction) (*mcms.Operation, error), error) {
	// Set the transaction opts based whether or not the MCMS config was applied
	var opts *bind.TransactOpts
	if mcmsConfig != nil {
		opts = deployment.SimTransactOpts()
	} else {
		opts = deployerKey
	}

	// Create handler
	handler := func(tx *types.Transaction) (*mcms.Operation, error) {
		if opts.From == deployerKey.From {
			_, err := chain.Confirm(tx)
			if err != nil {
				return nil, fmt.Errorf("failed to confirm transaction with hash %s on %s: %w", tx.Hash(), chain.Name(), err)
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
