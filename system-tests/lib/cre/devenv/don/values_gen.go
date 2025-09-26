package don

import (
	"fmt"
)

// ChainlinkValuesConfig represents the complete Helm values configuration for Chainlink deployment
type ChainlinkValuesConfig struct {
	BootNodeCount    int                  `yaml:"bootNodeCount,omitempty"`
	NodeCount        int                  `yaml:"nodeCount,omitempty"`
	FullnameOverride string               `yaml:"fullnameOverride,omitempty"`
	NetworkPolicy    NetworkPolicyConfig  `yaml:"networkPolicy,omitempty"`
	Common           CommonConfig         `yaml:"common,omitempty"`
	Overrides        []NodeOverrideConfig `yaml:"overrides,omitempty"`
	Rollout          RolloutConfig        `yaml:"rollout,omitempty"`
	ServiceAccount   ServiceAccountConfig `yaml:"serviceAccount,omitempty"`
}

// NodeValuesConfig represents configuration for individual nodes
type NodeValuesConfig struct {
	IsBootstrap     bool
	Image           string
	AppInstanceName string
	Config          string
	SecretsOverride string
	DatabaseHost    string
}

// NetworkPolicyConfig represents network policy settings
type NetworkPolicyConfig struct {
	Enabled *bool `yaml:"enabled,omitempty"`
}

// RolloutConfig represents rollout settings
type RolloutConfig struct {
	Enabled *bool `yaml:"enabled,omitempty"`
}

// ServiceAccountConfig represents service account settings
type ServiceAccountConfig struct {
	Enabled *bool  `yaml:"enabled,omitempty"`
	Name    string `yaml:"name,omitempty"`
}

// CommonConfig represents shared configuration for all nodes
type CommonConfig struct {
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	Chainlink          ChainlinkConfig        `yaml:"chainlink,omitempty"`
	ChainlinkNode      ChainlinkNodeConfig    `yaml:"chainlinkNode,omitempty"`
	Image              ImageConfig            `yaml:"image,omitempty"`
	ImagePullSecrets   []ImagePullSecret      `yaml:"imagePullSecrets,omitempty"`
	Ingress            IngressConfig          `yaml:"ingress,omitempty"`
	NodeSelector       map[string]interface{} `yaml:"nodeSelector,omitempty"`
	PodSecurityContext PodSecurityContext     `yaml:"podSecurityContext,omitempty"`
	RequiredLabels     map[string]string      `yaml:"requiredLabels,omitempty"`
	Resources          ResourcesConfig        `yaml:"resources,omitempty"`
	Service            ServiceConfig          `yaml:"service,omitempty"`
	ServiceMonitor     ServiceMonitorConfig   `yaml:"servicemonitor,omitempty"`
	Tolerations        []interface{}          `yaml:"tolerations,omitempty"`
}

// ChainlinkConfig represents Chainlink v2 configuration
type ChainlinkConfig struct {
	V2Config map[string]string `yaml:"v2Config,omitempty"`
}

// ChainlinkNodeConfig represents Chainlink node specific configuration
type ChainlinkNodeConfig struct {
	Enabled  *bool                 `yaml:"enabled,omitempty"`
	Metadata ChainlinkNodeMetadata `yaml:"metadata,omitempty"`
	Spec     ChainlinkNodeSpec     `yaml:"spec,omitempty"`
}

// ChainlinkNodeMetadata represents node metadata
type ChainlinkNodeMetadata struct {
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// ChainlinkNodeSpec represents node specification
type ChainlinkNodeSpec struct {
	Credentials CredentialsConfig `yaml:"credentials,omitempty"`
	Database    DatabaseConfig    `yaml:"database,omitempty"`
}

// CredentialsConfig represents node credentials
type CredentialsConfig struct {
	Config CredentialDetails `yaml:"config,omitempty"`
}

// CredentialDetails represents credential details
type CredentialDetails struct {
	API      APICredential      `yaml:"api,omitempty"`
	Keystore KeystoreCredential `yaml:"keystore,omitempty"`
	VRF      VRFCredential      `yaml:"vrf,omitempty"`
}

// APICredential represents API credentials
type APICredential struct {
	Key      string `yaml:"key,omitempty"`
	Password string `yaml:"password,omitempty"`
	User     string `yaml:"user,omitempty"`
}

// KeystoreCredential represents keystore credentials
type KeystoreCredential struct {
	Key      string `yaml:"key,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// VRFCredential represents VRF credentials
type VRFCredential struct {
	Key      string `yaml:"key,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Config DatabaseDetails `yaml:"config,omitempty"`
}

// DatabaseDetails represents database connection details
type DatabaseDetails struct {
	Database string      `yaml:"database,omitempty"`
	Host     interface{} `yaml:"host,omitempty"`
	Password string      `yaml:"password,omitempty"`
	Port     string      `yaml:"port,omitempty"`
	User     string      `yaml:"user,omitempty"`
}

// ImageConfig represents Docker image configuration
type ImageConfig struct {
	PullPolicy string `yaml:"pullPolicy,omitempty"`
	Repository string `yaml:"repository,omitempty"`
	Tag        string `yaml:"tag,omitempty"`
}

// ImagePullSecret represents image pull secret
type ImagePullSecret struct {
	Name string `yaml:"name,omitempty"`
}

// IngressConfig represents ingress configuration
type IngressConfig struct {
	Annotations      map[string]string `yaml:"annotations,omitempty"`
	Enabled          *bool             `yaml:"enabled,omitempty"`
	Hosts            []IngressHost     `yaml:"hosts,omitempty"`
	IngressClassName string            `yaml:"ingressClassName,omitempty"`
}

// IngressHost represents ingress host configuration
type IngressHost struct {
	Host        string        `yaml:"host,omitempty"`
	Paths       []IngressPath `yaml:"paths,omitempty"`
	UseNodeName *bool         `yaml:"useNodeName,omitempty"`
}

// IngressPath represents ingress path configuration
type IngressPath struct {
	Path     string `yaml:"path,omitempty"`
	PathType string `yaml:"pathType,omitempty"`
}

// PodSecurityContext represents pod security context
type PodSecurityContext struct {
	RunAsGroup   int   `yaml:"runAsGroup,omitempty"`
	RunAsNonRoot *bool `yaml:"runAsNonRoot,omitempty"`
	RunAsUser    int   `yaml:"runAsUser,omitempty"`
}

// ResourcesConfig represents resource constraints
type ResourcesConfig struct {
	Limits   ResourceLimits   `yaml:"limits,omitempty"`
	Requests ResourceRequests `yaml:"requests,omitempty"`
}

// ResourceLimits represents resource limits
type ResourceLimits struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// ResourceRequests represents resource requests
type ResourceRequests struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// ServiceConfig represents service configuration
type ServiceConfig struct {
	Private ServicePrivate `yaml:"private,omitempty"`
}

// ServicePrivate represents private service configuration
type ServicePrivate struct {
	Type string `yaml:"type,omitempty"`
}

// ServiceMonitorConfig represents service monitor configuration
type ServiceMonitorConfig struct {
	Enabled *bool `yaml:"enabled,omitempty"`
}

// NodeOverrideConfig represents per-node override configuration
type NodeOverrideConfig struct {
	ChainlinkNode *ChainlinkNodeOverride `yaml:"chainlinkNode,omitempty"`
	Chainlink     *ChainlinkOverride     `yaml:"chainlink,omitempty"`
	V2Secret      map[string]string      `yaml:"v2Secret,omitempty"`
}

// ChainlinkNodeOverride represents node-specific overrides
type ChainlinkNodeOverride struct {
	Spec ChainlinkNodeOverrideSpec `yaml:"spec,omitempty"`
}

// ChainlinkNodeOverrideSpec represents node override spec
type ChainlinkNodeOverrideSpec struct {
	Database DatabaseConfig `yaml:"database,omitempty"`
}

// ChainlinkOverride represents Chainlink-specific overrides
type ChainlinkOverride struct {
	V2Config map[string]string `yaml:"v2Config,omitempty"`
}

// BoolPtr is a helper function to get a pointer to a boolean value
func BoolPtr(b bool) *bool {
	return &b
}

// NewChainlinkClusterValuesConfig creates a new ChainlinkValuesConfig with default values
func NewChainlinkClusterValuesConfig() *ChainlinkValuesConfig {
	return &ChainlinkValuesConfig{
		BootNodeCount:    0,
		NodeCount:        0,
		FullnameOverride: "base",
		NetworkPolicy: NetworkPolicyConfig{
			Enabled: BoolPtr(false),
		},
		Common: CommonConfig{
			Affinity: make(map[string]interface{}),
			Chainlink: ChainlinkConfig{
				V2Config: make(map[string]string),
			},
			ChainlinkNode: ChainlinkNodeConfig{
				Enabled: BoolPtr(false),
				Metadata: ChainlinkNodeMetadata{
					Annotations: map[string]string{
						"chainlinknode.k8s.chain.link/disable-tls": "false",
					},
				},
				Spec: ChainlinkNodeSpec{
					Credentials: CredentialsConfig{
						Config: CredentialDetails{
							API: APICredential{
								Key:      ".api",
								Password: "hWDmgcub2gUhyrG6cxriqt7T",
								User:     "admin@chain.link",
							},
							Keystore: KeystoreCredential{
								Key:      ".keystore",
								Password: "cdz7KhvF4ATje2TjrwGMJh2Q",
							},
							VRF: VRFCredential{
								Key:      ".vrf",
								Password: "cdz7KhvF4ATje2TjrwGMJh2Q",
							},
						},
					},
					Database: DatabaseConfig{
						Config: DatabaseDetails{
							Database: "chainlink",
							Host:     nil,
							Password: "JGVgp7M2Emcg7Av8KKVUgMZb",
							Port:     "5432",
							User:     "chainlink",
						},
					},
				},
			},
			Image: ImageConfig{
				PullPolicy: "Always",
				Tag:        "develop",
			},
			NodeSelector: make(map[string]interface{}),
			PodSecurityContext: PodSecurityContext{
				RunAsGroup:   14933,
				RunAsNonRoot: BoolPtr(true),
				RunAsUser:    14933,
			},
			RequiredLabels: map[string]string{
				"app.chain.link/blockchain":   "multichain",
				"app.chain.link/network":      "multichain",
				"app.chain.link/network-type": "testnet",
				"app.chain.link/product":      "base",
				"app.chain.link/team":         "infra",
				"app.kubernetes.io/component": "chainlink",
			},
			Resources: ResourcesConfig{
				Limits: ResourceLimits{
					CPU:    "1",
					Memory: "1Gi",
				},
				Requests: ResourceRequests{
					CPU:    "0.1",
					Memory: "256Mi",
				},
			},
			Service: ServiceConfig{
				Private: ServicePrivate{
					Type: "ClusterIP",
				},
			},
			ServiceMonitor: ServiceMonitorConfig{
				Enabled: BoolPtr(false),
			},
			Tolerations: []interface{}{},
		},
		Rollout: RolloutConfig{
			Enabled: BoolPtr(false),
		},
		ServiceAccount: ServiceAccountConfig{
			Enabled: BoolPtr(false),
			Name:    "default",
		},
	}
}

// NewNodeValuesConfig creates a new regular node configuration
func NewNodeValuesConfig() *NodeValuesConfig {
	return &NodeValuesConfig{
		IsBootstrap: false,
	}
}

// NewBootNodeValuesConfig creates a new bootstrap node configuration
func NewBootNodeValuesConfig() *NodeValuesConfig {
	return &NodeValuesConfig{
		IsBootstrap: true,
	}
}

// SetImage sets the Docker image for the node
func (n *NodeValuesConfig) SetImage(image string) *NodeValuesConfig {
	n.Image = image
	return n
}

// SetAppInstanceName sets the app instance name for the node
func (n *NodeValuesConfig) SetAppInstanceName(name string) *NodeValuesConfig {
	n.AppInstanceName = name
	return n
}

// SetConfig sets the configuration TOML for the node
func (n *NodeValuesConfig) SetConfig(config string) *NodeValuesConfig {
	n.Config = config
	return n
}

// SetSecretsOverride sets the secrets override TOML for the node
func (n *NodeValuesConfig) SetSecretsOverride(secrets string) *NodeValuesConfig {
	n.SecretsOverride = secrets
	return n
}

// SetDatabaseHost sets the database host for the node
func (n *NodeValuesConfig) SetDatabaseHost(host string) *NodeValuesConfig {
	n.DatabaseHost = host
	return n
}

// Build finalizes the node configuration
func (n *NodeValuesConfig) Build() NodeValuesConfig {
	if n.DatabaseHost == "" {
		if n.IsBootstrap {
			n.DatabaseHost = fmt.Sprintf("base-db-bt-0")
		} else {
			n.DatabaseHost = fmt.Sprintf("base-db-0")
		}
	}
	return *n
}

// SetNodes sets the nodes for the cluster configuration
func (c *ChainlinkValuesConfig) SetNodes(nodes []NodeValuesConfig) *ChainlinkValuesConfig {
	c.NodeCount = len(nodes)
	c.BootNodeCount = 0

	// Count bootstrap nodes
	for _, node := range nodes {
		if node.IsBootstrap {
			c.BootNodeCount++
		}
	}

	// Generate overrides for each node
	c.Overrides = []NodeOverrideConfig{}

	for i, node := range nodes {
		// Database host override
		dbOverride := NodeOverrideConfig{
			ChainlinkNode: &ChainlinkNodeOverride{
				Spec: ChainlinkNodeOverrideSpec{
					Database: DatabaseConfig{
						Config: DatabaseDetails{
							Host: node.DatabaseHost,
						},
					},
				},
			},
		}
		c.Overrides = append(c.Overrides, dbOverride)

		// Chainlink configuration override
		chainlinkOverride := NodeOverrideConfig{
			Chainlink: &ChainlinkOverride{
				V2Config: map[string]string{
					"99-config-override.toml": node.Config,
				},
			},
			V2Secret: map[string]string{
				"99-secrets-override.toml": node.SecretsOverride,
			},
		}
		c.Overrides = append(c.Overrides, chainlinkOverride)

		// Update image in common config from first node
		if i == 0 {
			c.Common.Image.Repository = node.Image
		}
	}

	return c
}

// SetFullnameOverride sets the fullname override
func (c *ChainlinkValuesConfig) SetFullnameOverride(name string) *ChainlinkValuesConfig {
	c.FullnameOverride = name
	return c
}

// Build generates the complete configuration
func (c *ChainlinkValuesConfig) Build() *ChainlinkValuesConfig {
	// todo add validation and custom logice here
	return c
}
