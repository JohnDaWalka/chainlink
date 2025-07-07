package solana

import (
    "testing"
		"github.com/Masterminds/semver/v3"

		chain_selectors "github.com/smartcontractkit/chain-selectors"
    "github.com/smartcontractkit/wsrpc/logger"
    "github.com/stretchr/testify/require"
    "go.uber.org/zap/zapcore"

    cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
    "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
    commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
    "github.com/smartcontractkit/chainlink/deployment/environment/memory"
    "github.com/smartcontractkit/chainlink/deployment/helpers"
)

func TestDeployCache(t *testing.T) {
    skipInCI(t)
    t.Parallel()

    lggr := logger.Test(t)
    cfg := memory.MemoryEnvironmentConfig{
        Nodes:     1,
        SolChains: 1,
    }

    env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
    solSel := env.BlockChains.ListChainSelectors(cldfchain.WithFamily(chain_selectors.FamilySolana))[0]

    chain := env.BlockChains.SolanaChains()[solSel]
    chain.ProgramsPath = getProgramsPath()
    env.BlockChains = cldfchain.NewBlockChains(map[uint64]cldfchain.BlockChain{solSel: chain})

    t.Run("should deploy cache", func(t *testing.T) {
        configuredChangeset := commonchangeset.Configure(DeployCache{},
            &DeployCacheRequest{
                ChainSel:  solSel,
                Qualifier: testQualifier,
                Version:   "1.0.0",
                BuildConfig: &helpers.BuildSolanaConfig{
                    GitCommitSha:   "ba5a33ab378020fac73bda72b6bc2f9ae6bddb83",
                    DestinationDir: getProgramsPath(),
                    LocalBuild:     helpers.LocalBuildConfig{BuildLocally: true, CreateDestinationDir: true},
                },
            },
        )

        var err error
        env, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
        require.NoError(t, err)

        // Check that the cache program and state addresses are present in the datastore
        ds := env.DataStore
        version := "1.0.0"
        cacheKey := datastore.NewAddressRefKey(solSel, CacheContract, mustParseVersion(version), testQualifier)
        cacheStateKey := datastore.NewAddressRefKey(solSel, CacheState, mustParseVersion(version), testQualifier)

        cacheAddr, err := ds.Addresses().Get(cacheKey)
        require.NoError(t, err)
        require.NotEmpty(t, cacheAddr.Address)

        cacheStateAddr, err := ds.Addresses().Get(cacheStateKey)
        require.NoError(t, err)
        require.NotEmpty(t, cacheStateAddr.Address)
    })

		t.Run("should pass upgrade authority", func(t *testing.T) {
        configuredChangeset := commonchangeset.Configure(SetCacheUpgradeAuthority{},
            &SetCacheUpgradeAuthorityRequest{
                ChainSel:            solSel,
                Qualifier:           testQualifier,
                Version:             "1.0.0",
                NewUpgradeAuthority: chain.DeployerKey.PublicKey().String(),
            },
        )

        var err error
        _, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
        require.NoError(t, err)
    })
}

func ParseSemver(v string) *semver.Version {
    ver, err := semver.NewVersion(v)
    if err != nil {
        panic(err)
    }
    return ver
}

func mustParseVersion(v string) *semver.Version {
    return ParseSemver(v)
}