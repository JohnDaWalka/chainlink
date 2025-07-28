package modsec

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
)

func validSpec() *job.ModsecSpec {
	return &job.ModsecSpec{
		SourceChainID:           "1338",
		SourceChainFamily:       string(chaintype.EVM),
		DestChainID:             "1337",
		DestChainFamily:         string(chaintype.EVM),
		OnRampAddress:           "0x1234567890123456789012345678901234567890",
		OffRampAddress:          "0x0987654321098765432109876543210987654321",
		CCIPMessageSentEventSig: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		modifier    func(*job.ModsecSpec)
		expectedErr string
	}{
		{
			name:        "valid spec",
			modifier:    func(spec *job.ModsecSpec) {},
			expectedErr: "",
		},
		{
			name: "empty source chain id",
			modifier: func(spec *job.ModsecSpec) {
				spec.SourceChainID = ""
			},
			expectedErr: "source chain id () or family (evm) is empty",
		},
		{
			name: "empty source chain family",
			modifier: func(spec *job.ModsecSpec) {
				spec.SourceChainFamily = ""
			},
			expectedErr: "source chain id (1338) or family () is empty",
		},
		{
			name: "empty dest chain id",
			modifier: func(spec *job.ModsecSpec) {
				spec.DestChainID = ""
			},
			expectedErr: "dest chain id () or family (evm) is empty",
		},
		{
			name: "empty dest chain family",
			modifier: func(spec *job.ModsecSpec) {
				spec.DestChainFamily = ""
			},
			expectedErr: "dest chain id (1337) or family () is empty",
		},
		{
			name: "source and dest chain same",
			modifier: func(spec *job.ModsecSpec) {
				spec.DestChainID = "1338"
			},
			expectedErr: "source chain id (1338) and dest chain id (1338) are the same",
		},
		{
			name: "empty onramp address",
			modifier: func(spec *job.ModsecSpec) {
				spec.OnRampAddress = ""
			},
			expectedErr: "on ramp address is empty",
		},
		{
			name: "invalid onramp address",
			modifier: func(spec *job.ModsecSpec) {
				spec.OnRampAddress = "not-an-address"
			},
			expectedErr: "on ramp address (not-an-address) is not a valid address",
		},
		{
			name: "empty offramp address",
			modifier: func(spec *job.ModsecSpec) {
				spec.OffRampAddress = ""
			},
			expectedErr: "off ramp address is empty",
		},
		{
			name: "invalid offramp address",
			modifier: func(spec *job.ModsecSpec) {
				spec.OffRampAddress = "not-an-address"
			},
			expectedErr: "off ramp address (not-an-address) is not a valid address",
		},
		{
			name: "empty ccip message sent event sig",
			modifier: func(spec *job.ModsecSpec) {
				spec.CCIPMessageSentEventSig = ""
			},
			expectedErr: "ccip message sent event sig is empty",
		},
		{
			name: "invalid ccip message sent event sig",
			modifier: func(spec *job.ModsecSpec) {
				spec.CCIPMessageSentEventSig = "0x1234"
			},
			expectedErr: "ccip message sent event sig is not 32 bytes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spec := validSpec()
			tc.modifier(spec)
			err := validate(spec)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			}
		})
	}
}

func TestValidatedModsecSpec(t *testing.T) {
	t.Parallel()

	validToml := `
		type = "modsec"
		schemaVersion = 1
		sourceChainID = "1338"
		sourceChainFamily = "evm"
		destChainID = "1337"
		destChainFamily = "evm"
		onRampAddress = "0x1234567890123456789012345678901234567890"
		offRampAddress = "0x0987654321098765432109876543210987654321"
		ccipMessageSentEventSig = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	`

	testCases := []struct {
		name        string
		toml        string
		expectedErr string
	}{
		{
			name:        "valid spec",
			toml:        validToml,
			expectedErr: "",
		},
		{
			name:        "invalid toml",
			toml:        "invalid-toml",
			expectedErr: "toml error on load",
		},
		{
			name: "unmarshal error on job",
			toml: `
				type = "modsec"
				schemaVersion = "not-a-number"
			`,
			expectedErr: "toml unmarshal error on job",
		},
		{
			name: "wrong job type",
			toml: `
				type = "not-modsec"
				schemaVersion = 1
			`,
			expectedErr: "the only supported type is currently 'modsec', got not-modsec",
		},
		{
			name: "validation error",
			toml: `
				type = "modsec"
				schemaVersion = 1
				sourceChainID = ""
			`,
			expectedErr: "source chain id () or family () is empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidatedModsecSpec(tc.toml)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			}
		})
	}
}
