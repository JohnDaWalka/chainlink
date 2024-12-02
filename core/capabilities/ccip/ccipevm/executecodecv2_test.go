package ccipevm

import (
	"testing"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteCodec(t *testing.T) {
	d := testSetup(t)
	input := randomExecuteReport(t, d)
	c, err := codec.NewCodec(execCodecConfig)
	require.NoError(t, err)

	result, err := c.Encode(testutils.Context(t), input, "ExecPluginReport")
	require.NoError(t, err)
	require.NotNil(t, result)
	output := cciptypes.ExecutePluginReport{}
	err = c.Decode(testutils.Context(t), result, &output, "ExecPluginReport")
	require.NoError(t, err)
	require.Equal(t, input, output)
}

func TestExecutePluginCodecV2(t *testing.T) {
	d := testSetup(t)

	testCases := []struct {
		name   string
		report func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport
		expErr bool
	}{
		{
			name:   "base report",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport { return report },
			expErr: false,
		},
		{
			name: "reports have empty msgs",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].Messages = []cciptypes.Message{}
				report.ChainReports[4].Messages = []cciptypes.Message{}
				return report
			},
			expErr: false,
		},
		{
			name: "reports have empty offchain token data",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].OffchainTokenData = [][][]byte{}
				report.ChainReports[4].OffchainTokenData[1] = [][]byte{}
				return report
			},
			expErr: false,
		},
	}

	ctx := testutils.Context(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cd := NewExecutePluginCodecV2()
			report := tc.report(randomExecuteReport(t, d))
			bytes, err := cd.Encode(ctx, report)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			testSetup(t)

			// ignore msg hash in comparison
			for i := range report.ChainReports {
				for j := range report.ChainReports[i].Messages {
					report.ChainReports[i].Messages[j].Header.MsgHash = cciptypes.Bytes32{}
					report.ChainReports[i].Messages[j].Header.OnRamp = cciptypes.UnknownAddress{}
					report.ChainReports[i].Messages[j].FeeToken = cciptypes.UnknownAddress{}
					report.ChainReports[i].Messages[j].ExtraArgs = cciptypes.Bytes{}
					report.ChainReports[i].Messages[j].FeeTokenAmount = cciptypes.BigInt{}
				}
			}

			// decode using the codec
			codecDecoded, err := cd.Decode(ctx, bytes)
			require.NoError(t, err)
			assert.Equal(t, report, codecDecoded)
		})
	}
}
