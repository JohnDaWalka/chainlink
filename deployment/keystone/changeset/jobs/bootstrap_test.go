package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveBootstrapJob(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg BootstrapCfg
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				cfg: BootstrapCfg{
					JobName:       "OCR3 MultiChain Capability Bootstrap (for DON workflow_1)",
					ExternalJobID: "f1ac5211-ab79-4c31-ba1c-0997b72db466",
					ContractID:    "0x123",
					ChainID:       "11155111",
				},
			},
			want: `type = "bootstrap"
schemaVersion = 1
name = "OCR3 MultiChain Capability Bootstrap (for DON workflow_1)"
externalJobID = "f1ac5211-ab79-4c31-ba1c-0997b72db466"
contractID = "0x123"
relay = "evm"

[relayConfig]
chainID = 11155111
providerType = "ocr3-capability"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ResolveBootstrapJob(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveBootstrapJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
