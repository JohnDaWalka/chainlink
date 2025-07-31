package jobs

// BootstrapperCfg is the configuration for a bootstrapper node.
type BootstrapperCfg struct {
	Name       string `json:"name" toml:"name"`
	NOP        string `json:"nop" toml:"nop"`
	CSAKey     string `json:"csa_key" toml:"csa_key"`
	P2PID      string `json:"p2p_id" toml:"p2p_id"`
	OCRUrl     string `json:"ocr_url" toml:"ocr_url"`
	DON2DONUrl string `json:"don2don_url" toml:"don2don_url"`
}
