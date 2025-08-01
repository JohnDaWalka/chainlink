package jobs

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

const bootstrapPth = "bootstrap.tmpl"

// BootstrapperCfg is the configuration for a bootstrapper node.
type BootstrapperCfg struct {
	Name       string `json:"name" toml:"name"`
	NOP        string `json:"nop" toml:"nop"`
	CSAKey     string `json:"csa_key" toml:"csa_key"`
	P2PID      string `json:"p2p_id" toml:"p2p_id"`
	OCRUrl     string `json:"ocr_url" toml:"ocr_url"`
	DON2DONUrl string `json:"don2don_url" toml:"don2don_url"`
}

type BootstrapCfg struct {
	JobName       string
	ExternalJobID string
	ContractID    string // contract ID of the ocr3 contract
	ChainID       string
}

func (c BootstrapCfg) Validate() error {
	if c.ContractID == "" {
		return errors.New("ContractID is empty")
	}

	return nil
}

func ResolveBootstrapJob(cfg BootstrapCfg) (string, error) {
	t, err := template.New("s").ParseFS(tmplFS, bootstrapPth)
	if err != nil {
		return "", fmt.Errorf("failed to parse ocr3_spec.tmpl: %w", err)
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, bootstrapPth, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return b.String(), nil
}

func BootstrapExternalJobID(donName string, evmChainSel uint64) (string, error) {
	return ExternalJobID(donName+"-bootstrap", evmChainSel)
}
