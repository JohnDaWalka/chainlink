package changeset

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs/offchain"
	jobs2 "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/jobs"
)

const (
	defaultOCRJobSpecsTimeout = 120 * time.Second
)

type CsDistributeOCRJobSpecsConfig struct {
	ChainSelectorEVM   uint64
	ChainSelectorAptos uint64
	DONFilter          *offchain.DONFilter
	// DON2ContractID is a map of DON Name to OCR contract ID.  Multiple OCR contracts can be deployed
	// for a given chain, this map allows specifying the contract ID for each DON.
	// The map must be supplied if there are multiple OCR contracts for a chain.
	DON2ContractID   map[string]string
	BootstrapperCfgs []jobs.BootstrapperCfg
	DomainKey        string
	EnvLabel         string
}

var CsDistributeOCRJobSpecs cldf.ChangeSetV2[CsDistributeOCRJobSpecsConfig] = CsDistributeOCRJobSpecsImpl{}

type CsDistributeOCRJobSpecsImpl struct{}

func (c CsDistributeOCRJobSpecsImpl) Apply(e cldf.Environment, cfg CsDistributeOCRJobSpecsConfig) (cldf.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(e.GetContext(), defaultOCRJobSpecsTimeout)
	defer cancel()

	nodes, err := offchain.FetchNodesFromJD(ctx, e.Offchain, cfg.DONFilter)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow don nodes: %w", err)
	}
	if len(nodes) != cfg.DONFilter.Size {
		return cldf.ChangesetOutput{}, fmt.Errorf("expected %d nodes, got %d", cfg.DONFilter.Size, len(nodes))
	}
	nodesByID := make(map[string]*nodev1.Node)
	nodeIDs := make([]string, 0, len(nodes))
	for _, node := range nodes {
		nodesByID[node.Id] = node
		nodeIDs = append(nodeIDs, node.Id)
	}

	contractID, err := getOCRContractID(cfg.DONFilter.DONName, cfg.DON2ContractID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get OCR contract ID: %w", err)
	}

	addresses := e.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(cfg.ChainSelectorEVM),
		datastore.AddressRefByAddress(contractID),
	)
	if len(addresses) == 0 {
		return cldf.ChangesetOutput{}, fmt.Errorf("no addresses found for chain selector %d and contract ID %s", cfg.ChainSelectorEVM, contractID)
	}
	addr := addresses[0]
	if cldf.ContractType(addr.Type) != OCR3Capability {
		return cldf.ChangesetOutput{}, fmt.Errorf("address %s for chain selector %d is not of type OCR3Capability", addr.Address, cfg.ChainSelectorEVM)
	}
	e.Logger.Debugw("found OCR contract ID", "contractID", contractID)

	btURLs := make([]string, 0, len(cfg.BootstrapperCfgs))
	for _, bootCfg := range cfg.BootstrapperCfgs {
		btURLs = append(btURLs, bootCfg.OCRUrl)
	}

	seqReport, errs := operations.ExecuteSequence(
		e.OperationsBundle,
		jobs2.DistributeOCRJobSpecSeq,
		jobs2.DistributeOCRJobSpecSeqDeps{
			NodeIDs:  nodeIDs,
			Offchain: e.Offchain,
		},
		jobs2.DistributeOCRJobSpecSeqInput{
			ContractID:           contractID,
			EnvironmentLabel:     cfg.EnvLabel,
			DomainKey:            cfg.DomainKey,
			DONName:              cfg.DONFilter.DONName,
			ChainSelectorEVM:     cfg.ChainSelectorEVM,
			ChainSelectorAptos:   cfg.ChainSelectorAptos,
			BootstrapperOCR3Urls: btURLs,
		},
	)
	if errs != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute distribute OCR job specs sequence: %w", errs)
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{seqReport.ToGenericReport()},
	}, errs
}

func (c CsDistributeOCRJobSpecsImpl) VerifyPreconditions(_ cldf.Environment, cfg CsDistributeOCRJobSpecsConfig) error {
	if cfg.DONFilter == nil {
		// Cannot get DON Name from the filter
		return errors.New("DON filter is nil")
	}
	return nil
}

func getOCRContractID(name string, don2ContractID map[string]string) (string, error) {
	if don2ContractID == nil {
		return "", errors.New("no map of DON name to contract ID provided")
	}

	contractID, ok := don2ContractID[name]
	if !ok {
		return "", fmt.Errorf("OCR contract ID not found for DON %s", name)
	}

	return contractID, nil
}
