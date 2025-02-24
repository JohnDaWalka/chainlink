package types

type NodeEthKeySelector struct {
	ChainSelector uint64 `toml:"ChainSelector"`
}
type NodeEthKey struct {
	JSON     string             `toml:"JSON"`
	Password string             `toml:"Password"`
	Selector NodeEthKeySelector `toml:"Selector"`
}

type NodeP2PKey struct {
	JSON     string `toml:"JSON"`
	Password string `toml:"Password"`
}

type NodeSecret struct {
	EthKey NodeEthKey `toml:"EthKey"`
	P2PKey NodeP2PKey `toml:"P2PKey"`
	// Add more fields as needed to reflect 'Secrets' struct from /core/config/toml/types.go
	// We can't use the original struct, because it's using custom types that serlialize secrets to 'xxxxx'
}
