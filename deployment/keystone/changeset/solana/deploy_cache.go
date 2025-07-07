package solana

import (
    "fmt"

    "github.com/Masterminds/semver/v3"
    "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
    cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
    "github.com/smartcontractkit/chainlink-deployments-framework/operations"
    "github.com/smartcontractkit/chainlink/deployment/helpers"
    seq "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence"
    "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence/operation"
)

const (
    CacheContract datastore.ContractType = "DataFeedsCache"
    CacheState    datastore.ContractType = "DataFeedsCacheState"
)

type DeployCacheRequest struct {
    ChainSel    uint64
    BuildConfig *helpers.BuildSolanaConfig
    Qualifier   string
    LabelSet    datastore.LabelSet
    Version     string
}

var _ cldf.ChangeSetV2[*DeployCacheRequest] = DeployCache{}

type DeployCache struct{}

func (cs DeployCache) VerifyPreconditions(env cldf.Environment, req *DeployCacheRequest) error {
    if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
        return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
    }
    if _, err := semver.NewVersion(req.Version); err != nil {
        return err
    }
    return nil
}

func (cs DeployCache) Apply(env cldf.Environment, req *DeployCacheRequest) (cldf.ChangesetOutput, error) {
    var out cldf.ChangesetOutput

    if req.BuildConfig != nil {
        // You may want to define a specific build params for the cache if needed
        err := helpers.BuildSolana(env, *req.BuildConfig, keystoneBuildParams)
        if err != nil {
            return out, fmt.Errorf("failed build solana artifacts: %w", err)
        }
    }

    out.DataStore = datastore.NewMemoryDataStore()
    version := semver.MustParse(req.Version)
    ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
    if !ok {
        return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
    }

    deploySeqInput := seq.DeployCacheSeqInput{
        ChainSel:    req.ChainSel,
        ProgramName: "data_feeds_cache", // Update if the program name is different
    }

    deps := operation.Deps{
        Datastore: env.DataStore,
        Env:       env,
        Chain:     ch,
    }

    deploySeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployCacheSeq, deps, deploySeqInput)
    if err != nil {
        return out, err
    }

    // Save programID
    err = out.DataStore.Addresses().Add(
        datastore.AddressRef{
            Address:       deploySeqReport.Output.ProgramID.String(),
            ChainSelector: req.ChainSel,
            Type:          CacheContract,
            Version:       version,
            Qualifier:     req.Qualifier,
            Labels:        req.LabelSet,
        },
    )
    if err != nil {
        return out, err
    }
    // Save StateID
    err = out.DataStore.Addresses().Add(
        datastore.AddressRef{
            Address:       deploySeqReport.Output.State.String(),
            ChainSelector: req.ChainSel,
            Type:          CacheState,
            Version:       version,
            Qualifier:     req.Qualifier,
            Labels:        req.LabelSet,
        },
    )
    if err != nil {
        return out, err
    }

    return out, nil
}