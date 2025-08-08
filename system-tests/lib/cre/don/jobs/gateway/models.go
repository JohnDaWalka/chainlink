package gateway

type GatewayJobSpec struct {
	Type              string        `toml:"type"`
	SchemaVersion     int           `toml:"schemaVersion"`
	ExternalJobID     string        `toml:"externalJobID"`
	Name              string        `toml:"name"`
	ForwardingAllowed bool          `toml:"forwardingAllowed"`
	GatewayConfig     GatewayConfig `toml:"gatewayConfig"`
}

type GatewayConfig struct {
	ConnectionManagerConfig ConnectionManagerConfig `toml:"ConnectionManagerConfig"`
	DONs                    []DON                   `toml:"Dons"`
	NodeServerConfig        NodeServerConfig        `toml:"NodeServerConfig"`
	UserServerConfig        UserServerConfig        `toml:"UserServerConfig"`
	HTTPClientConfig        HTTPClientConfig        `toml:"HTTPClientConfig"`
}

type NodeServerConfig struct {
	HandshakeTimeoutMillis int    `toml:"HandshakeTimeoutMillis"`
	MaxRequestBytes        int    `toml:"MaxRequestBytes"`
	Path                   string `toml:"Path"`
	Port                   int    `toml:"Port"`
	ReadTimeoutMillis      int    `toml:"ReadTimeoutMillis"`
	RequestTimeoutMillis   int    `toml:"RequestTimeoutMillis"`
	WriteTimeoutMillis     int    `toml:"WriteTimeoutMillis"`
}

type UserServerConfig struct {
	ContentTypeHeader    string   `toml:"ContentTypeHeader"`
	MaxRequestBytes      int      `toml:"MaxRequestBytes"`
	Path                 string   `toml:"Path"`
	Port                 int      `toml:"Port"`
	ReadTimeoutMillis    int      `toml:"ReadTimeoutMillis"`
	RequestTimeoutMillis int      `toml:"RequestTimeoutMillis"`
	WriteTimeoutMillis   int      `toml:"WriteTimeoutMillis"`
	CORSEnabled          bool     `toml:"CORSEnabled"`
	CORSAllowedOrigins   []string `toml:"CORSAllowedOrigins"`
}

type HTTPClientConfig struct {
	MaxResponseBytes int      `toml:"MaxResponseBytes"`
	AllowedPorts     []int    `toml:"AllowedPorts,omitempty"`
	AllowedIps       []string `toml:"AllowedIps,omitempty"`
	AllowedIPsCIDR   []string `toml:"AllowedIPsCIDR,omitempty"`
}

type ConnectionManagerConfig struct {
	AuthChallengeLen          int    `toml:"AuthChallengeLen"`
	AuthGatewayId             string `toml:"AuthGatewayId"`
	AuthTimestampToleranceSec int    `toml:"AuthTimestampToleranceSec"`
	HeartbeatIntervalSec      int    `toml:"HeartbeatIntervalSec"`
}

type DON struct {
	DonId         string
	F             int
	HandlerName   string
	Members       []DONMember
	HandlerConfig HandlerConfig
	Handlers      []Handler
}

type Handler struct {
	Name        string
	ServiceName string
	Config      any
}

type HTTPCapabilitiesConfig struct {
	MaxAllowedMessageAgeSec int               `toml:"MaxAllowedMessageAgeSec"`
	AuthPullIntervalSec     int               `toml:"authPullIntervalSec"`
	NodeRateLimiter         RateLimiterConfig `toml:"NodeRateLimiter"`
	UserRateLimiter         RateLimiterConfig `toml:"UserRateLimiter"`
}

type HandlerConfig struct {
	MaxAllowedMessageAgeSec int
	NodeRateLimiter         RateLimiterConfig
}

type RateLimiterConfig struct {
	GlobalBurst    int
	GlobalRPS      int
	PerSenderBurst int
	PerSenderRPS   int
}

type DONMember struct {
	Address string
	Name    string
}
