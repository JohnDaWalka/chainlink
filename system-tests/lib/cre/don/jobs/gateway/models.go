package gateway

type GatewayJobSpec struct {
	Type                    string                  `json:"type"`
	SchemaVersion           int                     `json:"schemaVersion"`
	ExternalJobID           string                  `json:"externalJobID"`
	Name                    string                  `json:"name"`
	ForwardingAllowed       bool                    `json:"forwardingAllowed"`
	ConnectionManagerConfig ConnectionManagerConfig `json:"gatewayConfig.ConnectionManagerConfig"`
	DONs                    []DON                   `json:"gatewayConfig.DONs"`
	NodeServerConfig        NodeServerConfig        `json:"gatewayConfig.NodeServerConfig"`
	UserServerConfig        UserServerConfig        `json:"gatewayConfig.UserServerConfig"`
	HTTPClientConfig        HTTPClientConfig        `json:"gatewayConfig.HTTPClientConfig"`
}

type NodeServerConfig struct {
	HandshakeTimeoutMillis int    `json:"HandshakeTimeoutMillis"`
	MaxRequestBytes        int    `json:"MaxRequestBytes"`
	Path                   string `json:"Path"`
	Port                   int    `json:"Port"`
	ReadTimeoutMillis      int    `json:"ReadTimeoutMillis"`
	RequestTimeoutMillis   int    `json:"RequestTimeoutMillis"`
	WriteTimeoutMillis     int    `json:"WriteTimeoutMillis"`
}

type UserServerConfig struct {
	ContentTypeHeader    string   `json:"ContentTypeHeader"`
	MaxRequestBytes      int      `json:"MaxRequestBytes"`
	Path                 string   `json:"Path"`
	Port                 int      `json:"Port"`
	ReadTimeoutMillis    int      `json:"ReadTimeoutMillis"`
	RequestTimeoutMillis int      `json:"RequestTimeoutMillis"`
	WriteTimeoutMillis   int      `json:"WriteTimeoutMillis"`
	CORSEnabled          bool     `json:"CORSEnabled"`
	CORSAllowedOrigins   []string `json:"CORSAllowedOrigins"`
}

type HTTPClientConfig struct {
	MaxResponseBytes int      `json:"MaxResponseBytes"`
	AllowedPorts     []int    `json:"AllowedPorts,omitempty"`
	AllowedIps       []string `json:"AllowedIps,omitempty"`
	AllowedIPsCIDR   []string `json:"AllowedIPsCIDR,omitempty"`
}

type ConnectionManagerConfig struct {
	AuthChallengeLen          int    `json:"AuthChallengeLen"`
	AuthGatewayId             string `json:"AuthGatewayId"`
	AuthTimestampToleranceSec int    `json:"AuthTimestampToleranceSec"`
	HeartbeatIntervalSec      int    `json:"HeartbeatIntervalSec"`
}

type GatewayConfig struct {
	Dons []DON
}

type DON struct {
	DonId         string
	F             int
	HandlerName   string
	Members       []DONMember
	HandlerConfig HandlerConfig
}

type HandlerConfig struct {
	MaxAllowedMessageAgeSec int
	NodeRateLimiter         NodeRateLimiterConfig
}

type NodeRateLimiterConfig struct {
	GlobalBurst    int
	GlobalRPS      int
	PerSenderBurst int
	PerSenderRPS   int
}

type DONMember struct {
	Address string
	Name    string
}
