package don

import (
	"fmt"
)

// ChainlinkValuesConfig represents the complete Helm values configuration for Chainlink deployment
type ChainlinkValuesConfig struct {
	BootNodeCount    int                  `yaml:"bootNodeCount"`
	NodeCount        int                  `yaml:"nodeCount"`
	FullnameOverride string               `yaml:"fullnameOverride"`
	NetworkPolicy    NetworkPolicyConfig  `yaml:"networkPolicy"`
	Common           CommonConfig         `yaml:"common"`
	Overrides        []NodeOverrideConfig `yaml:"overrides"`
	Rollout          RolloutConfig        `yaml:"rollout"`
	ServiceAccount   ServiceAccountConfig `yaml:"serviceAccount"`
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
	Enabled bool `yaml:"enabled"`
}

// RolloutConfig represents rollout settings
type RolloutConfig struct {
	Enabled bool `yaml:"enabled"`
}

// ServiceAccountConfig represents service account settings
type ServiceAccountConfig struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
}

// CommonConfig represents shared configuration for all nodes
type CommonConfig struct {
	Affinity           map[string]interface{} `yaml:"affinity"`
	Chainlink          ChainlinkConfig        `yaml:"chainlink"`
	ChainlinkNode      ChainlinkNodeConfig    `yaml:"chainlinkNode"`
	Image              ImageConfig            `yaml:"image"`
	ImagePullSecrets   []ImagePullSecret      `yaml:"imagePullSecrets"`
	Ingress            IngressConfig          `yaml:"ingress"`
	NodeSelector       map[string]interface{} `yaml:"nodeSelector"`
	PodSecurityContext PodSecurityContext     `yaml:"podSecurityContext"`
	RequiredLabels     map[string]string      `yaml:"requiredLabels"`
	Resources          ResourcesConfig        `yaml:"resources"`
	Service            ServiceConfig          `yaml:"service"`
	ServiceMonitor     ServiceMonitorConfig   `yaml:"servicemonitor"`
	Tolerations        []interface{}          `yaml:"tolerations"`
}

// ChainlinkConfig represents Chainlink v2 configuration
type ChainlinkConfig struct {
	V2Config map[string]string `yaml:"v2Config"`
}

// ChainlinkNodeConfig represents Chainlink node specific configuration
type ChainlinkNodeConfig struct {
	Enabled  bool                  `yaml:"enabled"`
	Metadata ChainlinkNodeMetadata `yaml:"metadata"`
	Spec     ChainlinkNodeSpec     `yaml:"spec"`
}

// ChainlinkNodeMetadata represents node metadata
type ChainlinkNodeMetadata struct {
	Annotations map[string]string `yaml:"annotations"`
}

// ChainlinkNodeSpec represents node specification
type ChainlinkNodeSpec struct {
	Credentials CredentialsConfig `yaml:"credentials"`
	Database    DatabaseConfig    `yaml:"database"`
}

// CredentialsConfig represents node credentials
type CredentialsConfig struct {
	Config CredentialDetails `yaml:"config"`
}

// CredentialDetails represents credential details
type CredentialDetails struct {
	API      APICredential      `yaml:"api"`
	Keystore KeystoreCredential `yaml:"keystore"`
	VRF      VRFCredential      `yaml:"vrf"`
}

// APICredential represents API credentials
type APICredential struct {
	Key      string `yaml:"key"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
}

// KeystoreCredential represents keystore credentials
type KeystoreCredential struct {
	Key      string `yaml:"key"`
	Password string `yaml:"password"`
}

// VRFCredential represents VRF credentials
type VRFCredential struct {
	Key      string `yaml:"key"`
	Password string `yaml:"password"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Config DatabaseDetails `yaml:"config"`
}

// DatabaseDetails represents database connection details
type DatabaseDetails struct {
	Database string      `yaml:"database"`
	Host     interface{} `yaml:"host"`
	Password string      `yaml:"password"`
	Port     string      `yaml:"port"`
	User     string      `yaml:"user"`
}

// ImageConfig represents Docker image configuration
type ImageConfig struct {
	PullPolicy string `yaml:"pullPolicy"`
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

// ImagePullSecret represents image pull secret
type ImagePullSecret struct {
	Name string `yaml:"name"`
}

// IngressConfig represents ingress configuration
type IngressConfig struct {
	Annotations      map[string]string `yaml:"annotations"`
	Enabled          bool              `yaml:"enabled"`
	Hosts            []IngressHost     `yaml:"hosts"`
	IngressClassName string            `yaml:"ingressClassName"`
}

// IngressHost represents ingress host configuration
type IngressHost struct {
	Host        string        `yaml:"host"`
	Paths       []IngressPath `yaml:"paths"`
	UseNodeName bool          `yaml:"useNodeName"`
}

// IngressPath represents ingress path configuration
type IngressPath struct {
	Path     string `yaml:"path"`
	PathType string `yaml:"pathType"`
}

// PodSecurityContext represents pod security context
type PodSecurityContext struct {
	RunAsGroup   int  `yaml:"runAsGroup"`
	RunAsNonRoot bool `yaml:"runAsNonRoot"`
	RunAsUser    int  `yaml:"runAsUser"`
}

// ResourcesConfig represents resource constraints
type ResourcesConfig struct {
	Limits   ResourceLimits   `yaml:"limits"`
	Requests ResourceRequests `yaml:"requests"`
}

// ResourceLimits represents resource limits
type ResourceLimits struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// ResourceRequests represents resource requests
type ResourceRequests struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// ServiceConfig represents service configuration
type ServiceConfig struct {
	Private ServicePrivate `yaml:"private"`
}

// ServicePrivate represents private service configuration
type ServicePrivate struct {
	Type string `yaml:"type"`
}

// ServiceMonitorConfig represents service monitor configuration
type ServiceMonitorConfig struct {
	Enabled bool `yaml:"enabled"`
}

// NodeOverrideConfig represents per-node override configuration
type NodeOverrideConfig struct {
	ChainlinkNode *ChainlinkNodeOverride `yaml:"chainlinkNode,omitempty"`
	Chainlink     *ChainlinkOverride     `yaml:"chainlink,omitempty"`
	V2Secret      map[string]string      `yaml:"v2Secret,omitempty"`
}

// ChainlinkNodeOverride represents node-specific overrides
type ChainlinkNodeOverride struct {
	Spec ChainlinkNodeOverrideSpec `yaml:"spec"`
}

// ChainlinkNodeOverrideSpec represents node override spec
type ChainlinkNodeOverrideSpec struct {
	Database DatabaseConfig `yaml:"database"`
}

// ChainlinkOverride represents Chainlink-specific overrides
type ChainlinkOverride struct {
	V2Config map[string]string `yaml:"v2Config"`
}

// NewChainlinkClusterValuesConfig creates a new ChainlinkValuesConfig with default values
func NewChainlinkClusterValuesConfig() *ChainlinkValuesConfig {
	return &ChainlinkValuesConfig{
		BootNodeCount:    0,
		NodeCount:        0,
		FullnameOverride: "base",
		NetworkPolicy: NetworkPolicyConfig{
			Enabled: false,
		},
		Common: CommonConfig{
			Affinity: make(map[string]interface{}),
			Chainlink: ChainlinkConfig{
				V2Config: make(map[string]string),
			},
			ChainlinkNode: ChainlinkNodeConfig{
				Enabled: false,
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
				RunAsNonRoot: true,
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
				Enabled: false,
			},
			Tolerations: []interface{}{},
		},
		Rollout: RolloutConfig{
			Enabled: false,
		},
		ServiceAccount: ServiceAccountConfig{
			Enabled: false,
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
