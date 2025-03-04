package changeset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
)

var _ deployment.ViewState = ViewKeystone

func ViewKeystone(e deployment.Environment) (json.Marshaler, error) {
	lggr := e.Logger
	state, err := GetContractSets(e.Logger, &GetContractSetsRequest{
		Chains:      e.Chains,
		AddressBook: e.ExistingAddresses,
	})
	// this error is unrecoverable
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	var (
		viewErrs error
		mu       sync.Mutex
		wg       sync.WaitGroup
		finished int
		started  int
	)
	chainViews := make(map[string]KeystoneChainView)
	appendError := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			viewErrs = errors.Join(viewErrs, err)
		}
	}
	want := len(state.ContractSets)
	statMap := make(map[string]bool)
	// TODO, set a reasonable timeout in the caller env in CLD
	ctx, cancel := context.WithTimeout(e.GetContext(), 3*time.Minute)
	defer cancel()
	for chainSel, contracts := range state.ContractSets {
		wg.Add(1)
		go func(chainSel uint64, contracts ContractSet) {
			defer wg.Done()

			chainid, err := chainsel.ChainIdFromSelector(chainSel)
			if err != nil {
				err2 := fmt.Errorf("failed to resolve chain id for selector %d: %w", chainSel, err)
				lggr.Error(err2)
				appendError(err2)
				return
			}
			chainName, err := chainsel.NameFromChainId(chainid)
			if err != nil {
				err2 := fmt.Errorf("failed to resolve chain name for chain id %d: %w", chainid, err)
				lggr.Error(err2)
				appendError(err2)
				return
			}

			mu.Lock()
			e.Logger.Debugf("contract view start for %s, %d/%d", chainName, started, want)
			started++
			statMap[chainName] = false
			mu.Unlock()
			defer func() {
				mu.Lock()
				finished++
				statMap[chainName] = true
				e.Logger.Debugf("contract view done for %s done/want %d/%d, stat %v\n", chainName, finished, want, statMap)
				mu.Unlock()
			}()
			v, err := contracts.View(ctx, e.Logger)
			if err != nil {
				err2 := fmt.Errorf("failed to view chain %s: %w", chainName, err)
				lggr.Error(err2)
				appendError(err2)
				// don't return here, we want to view all chains
			}
			mu.Lock()
			chainViews[chainName] = v
			mu.Unlock()
		}(chainSel, contracts)
	}
	// wait for all chains to be viewed
	e.Logger.Debugf("waiting for all chains to be viewed")
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()
	select {
	case <-ctx.Done():
		e.Logger.Debugf("timed out waiting for all chains to be viewed")
		break
	case <-doneCh:
		e.Logger.Debugf("all chains viewed")
		break
	}
	nopsView, err := commonview.GenerateNopsView(e.NodeIDs, e.Offchain)
	if err != nil {
		err2 := fmt.Errorf("failed to view nops: %w", err)
		lggr.Error(err2)
		viewErrs = errors.Join(viewErrs, err2)
	}
	return &KeystoneView{
		Chains: chainViews,
		Nops:   nopsView,
	}, viewErrs
}
