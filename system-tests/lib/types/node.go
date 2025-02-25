package types

type NodeEthKeySelector struct {
	ChainSelector uint64 `toml:"ChainSelector"`
	ChainName     string `toml:"ChainName"`
}
type NodeEthKey struct {
	JSON     string             `toml:"JSON"`
	Password string             `toml:"Password"`
	Selector NodeEthKeySelector `toml:"ChainDetails"`
}

type NodeP2PKey struct {
	JSON     string `toml:"JSON"`
	Password string `toml:"Password"`
}

type NodeEthKeyWrapper struct {
	EthKeys []NodeEthKey `toml:"EthKeys"`
}

type NodeSecret struct {
	EthKeys NodeEthKeyWrapper `toml:"EthKeys"`
	P2PKey  NodeP2PKey        `toml:"P2PKey"`
	// Add more fields as needed to reflect 'Secrets' struct from /core/config/toml/types.go
	// We can't use the original struct, because it's using custom types that serlialize secrets to 'xxxxx'
}
