package ops

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OpEVMDeployLinkToken(t *testing.T) {
	t.Parallel()

	var (
		chainID       uint64 = 11155111
		chainSelector uint64 = 16015286601757825753
	)

	tests := []struct {
		name string
		give OpEVMDeployLinkTokenInput
		want OpEvmDeployLinkTokenOutput
	}{
		{
			name: "test",
			give: OpEVMDeployLinkTokenInput{
				ChainSelector: chainSelector,
				ChainName:     "test-chain",
			},
			want: OpEvmDeployLinkTokenOutput{
				Type:    LinkTokenTypeAndVersion1.Type.String(),
				Version: LinkTokenTypeAndVersion1.Version.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				chains, bindings = memory.NewMemoryChainsWithChainIDs(t, []uint64{chainID}, 1)
				chain            = chains[chainSelector]
				auth             = bindings[chainSelector][0]
				deps             = OpEVMDeployLinkTokenDeps{
					Auth:        auth,
					Backend:     chain.Client,
					ConfirmFunc: chain.Confirm,
				}
			)

			got, err := operations.ExecuteOperation(
				optest.NewBundle(t), OpEVMDeployLinkToken, deps, tt.give,
			)
			require.NoError(t, err)

			assert.NotEmpty(t, got.Output.Address.String())
			assert.Equal(t, tt.want.Type, got.Output.Type)
			assert.Equal(t, tt.want.Version, got.Output.Version)
		})
	}
}
