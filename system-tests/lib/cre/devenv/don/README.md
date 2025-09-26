# Chainlink DON Values Generator

A Go library for generating Helm values configurations for Chainlink Distributed Oracle Network (DON) deployments. Provides a type-safe, fluent API for creating complex Kubernetes configurations.

> **Note**: All keys, passwords, and other secret-like strings in this package are generated test data and are not real secrets.

## Architecture

### Core Components

**ChainlinkValuesConfig** - Main configuration representing a complete Helm values file
```go
type ChainlinkValuesConfig struct {
    BootNodeCount    int
    NodeCount        int
    FullnameOverride string
    NetworkPolicy    NetworkPolicyConfig
    Common           CommonConfig
    Overrides        []NodeOverrideConfig
    Rollout          RolloutConfig
    ServiceAccount   ServiceAccountConfig
}
```

**NodeValuesConfig** - Individual node configuration (bootstrap or regular)
```go
type NodeValuesConfig struct {
    IsBootstrap     bool
    Image           string
    AppInstanceName string
    Config          string          // TOML configuration
    SecretsOverride string          // TOML secrets
    DatabaseHost    string
}
```

**Builder API**
- `NewChainlinkClusterValuesConfig()` - Creates cluster configuration
- `NewNodeValuesConfig()` - Creates regular node configuration  
- `NewBootNodeValuesConfig()` - Creates bootstrap node configuration

## Usage

```go
// Create nodes
bootNode := NewBootNodeValuesConfig().
    SetImage("chainlink:latest").
    SetAppInstanceName("crib-bt-0").
    SetConfig(bootstrapConfigTOML).
    SetSecretsOverride(bootstrapSecretsTOML).
    Build()

regularNode := NewNodeValuesConfig().
    SetImage("chainlink:latest").
    SetAppInstanceName("crib-node-0").
    SetConfig(nodeConfigTOML).
    SetSecretsOverride(nodeSecretsTOML).
    Build()

// Create cluster
cluster := NewChainlinkClusterValuesConfig().
    SetNodes([]NodeValuesConfig{bootNode, regularNode}).
    SetFullnameOverride("my-chainlink-cluster").
    Build()

// Generate YAML
yamlData, err := yaml.Marshal(cluster)
```

## Generated Structure

```yaml
bootNodeCount: 1              # Auto-calculated from node types
nodeCount: 2                  # Auto-calculated from total nodes
fullnameOverride: "my-cluster"
networkPolicy:
  enabled: false              # Explicit boolean values preserved
common:
  chainlinkNode:
    enabled: false
    metadata:
      annotations:
        chainlinknode.k8s.chain.link/disable-tls: "false"
    spec:
      # ... shared configuration
  image:
    repository: "chainlink:latest"
    pullPolicy: "Always"
    tag: "develop"
  # ... other shared settings
overrides:
  - chainlinkNode:             # Database override for node 0
      spec:
        database:
          config:
            host: base-db-bt-0
  - chainlink:                 # Configuration override for node 0
      v2Config:
        99-config-override.toml: |
          # Node-specific TOML configuration
      v2Secret:
        99-secrets-override.toml: |
          # Node-specific secrets
  # ... additional node overrides
```

## Key Features

- **Automatic Node Counting**: Calculates `bootNodeCount` and `nodeCount` from provided configurations
- **Database Host Generation**: Auto-generates hosts (`base-db-bt-{index}` for bootstrap, `base-db-{index}` for regular)
- **Configuration Inheritance**: Common settings shared, node-specific overrides applied via `overrides` array
- **Clean YAML Output**: Only meaningful values included using `omitempty` with boolean pointer pattern
- **TOML Integration**: Embed TOML configuration and secrets as strings

## Implementation Details

### Boolean Pointer Pattern
Uses pointers to distinguish between explicit `false` and unset values:

```go
type NetworkPolicyConfig struct {
    Enabled *bool `yaml:"enabled,omitempty"`
}

func BoolPtr(b bool) *bool { return &b }

// Usage
config.NetworkPolicy.Enabled = BoolPtr(false)  // Explicit false → included
config.NetworkPolicy.Enabled = nil             // Unset → omitted
```

### Override System
Alternating database and configuration overrides for each node:
1. Database host override (per node)
2. Configuration + secrets override (per node)

This creates a predictable Helm pattern for node-specific customizations.

## Testing

Snapshot testing with `go-snaps` ensures output consistency:

```go
func TestValuesGenerating(t *testing.T) {
    // ... create configuration
    yamlData, err := yaml.Marshal(config)
    snaps.MatchSnapshot(t, string(yamlData))
}
```

Run tests: `go test -v`
Update snapshots: `go test -v --update-snapshots`
