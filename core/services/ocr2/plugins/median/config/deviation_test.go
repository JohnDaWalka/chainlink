package config

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/assert"
)

type frozenTimeClock struct{}

func (frozenTimeClock) Now() time.Time {
	return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
}

func valFromString(t *testing.T, v string) *big.Int {
	d, err := decimal.NewFromString(v)
	require.NoError(t, err)
	d = d.Mul(decimal.New(1, 18))
	return d.BigInt()
}

func Test_DeviationFunctionDefinition(t *testing.T) {
	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Run("invalid string", func(t *testing.T) {
			var d DeviationFunctionDefinition
			err := d.UnmarshalJSON([]byte(`"invalid"`))
			require.EqualError(t, err, "json: cannot unmarshal string into Go value of type map[string]interface {}")
		})
		t.Run("pendle - valid", func(t *testing.T) {
			expiresAt := float64(13857541.0) + float64(time.Now().Unix())
			var d DeviationFunctionDefinition
			err := d.UnmarshalJSON([]byte(fmt.Sprintf(`{"type": "pendle", "expiresAt": %f}`, expiresAt)))
			require.NoError(t, err)
			f := d.Func()
			require.NotNil(t, f)
			// Test the actual deviation function behavior
			assert.True(t, f(nil, 1e7, big.NewInt(0.187152977881070687*1e18), big.NewInt(0.160000000000000000*1e18)))
			assert.False(t, f(nil, 1e7, big.NewInt(0.187152977881070687*1e18), big.NewInt(0.177777777777777777*1e18)))
		})
		t.Run("pendle - with multiplier", func(t *testing.T) {
			var d DeviationFunctionDefinition
			err := d.UnmarshalJSON([]byte(fmt.Sprintf(`{"type": "pendle", "expiresAt": 123456.890, "multiplier": "1000"}`, expiresAt)))
			require.NoError(t, err)
		})
	})
}

func Test_PendleDeviationFunc(t *testing.T) {
	tcs := []struct {
		name string

		expiresInSeconds float64
		multiplierVal    *big.Int
		thresholdPPB     uint64
		oldVal           *big.Int
		newVal           *big.Int

		err      string
		expected bool
	}{
		{
			name:   "nil oldVal errors",
			oldVal: nil,
			newVal: big.NewInt(2),
			err:    "oldVal and newVal must be non-nil",
		},
		{
			name:   "nil newVal errors",
			oldVal: big.NewInt(1),
			newVal: nil,
			err:    "oldVal and newVal must be non-nil",
		},
		{
			name:             "test 0 (one block after) - SHOULD UPDATE",
			expiresInSeconds: 13857541.0,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.187152977881070687 * 1e18),
			newVal:           big.NewInt(0.164498448931278907 * 1e18),
			expected:         true,
		},
		{
			name:             "test 1 (same block) - SHOULD UPDATE",
			expiresInSeconds: 1.385755300000000056744 * 1e7,
			thresholdPPB:     1e7,
			oldVal:           valFromString(t, "0.187152977881070686771991518071"),
			newVal:           valFromString(t, "0.164498448931278906659514404964"),
			expected:         true,
		},
		{
			name:             "test 2 (same block) - SHOULD UPDATE",
			expiresInSeconds: 13825908.999999998137354850769,
			thresholdPPB:     1e7,
			oldVal:           valFromString(t, "0.164498448931278906659514404964"),
			newVal:           valFromString(t, "0.141802025539163406575582371261"),
			expected:         true,
		},
		{
			name:             "test 3 (same block) - SHOULD UPDATE",
			expiresInSeconds: 11564844.999999998137354850769,
			thresholdPPB:     1e7,
			oldVal:           valFromString(t, "0.141802025539163406575582371261"),
			newVal:           valFromString(t, "0.114668136842518697537940397524"),
			expected:         true,
		},
		{
			name:             "test 4 edge case - SHOULD UPDATE",
			expiresInSeconds: 11564856.999999998137354850769,
			thresholdPPB:     1e7,
			oldVal:           valFromString(t, "0.141802025539163406575582371261"),
			newVal:           valFromString(t, "0.114668695429674297181499298404"),
			expected:         false,
		},
		{
			name:             "test 5 (previous block) - SHOULD NOT UPDATE",
			expiresInSeconds: 13857565.0,
			thresholdPPB:     1e7,
			oldVal:           valFromString(t, "0.187152977881070687"),
			newVal:           valFromString(t, "0.164529304905415591"),
			expected:         false,
		},
		{
			name:             "test 6 (previous block) - SHOULD NOT UPDATE",
			expiresInSeconds: 13825932.999999998137354851,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.164498448931278907 * 1e18),
			newVal:           big.NewInt(0.141815210766795902 * 1e18),
			expected:         false,
		},
		{
			name:             "test 7 (previous block) - SHOULD NOT UPDATE",
			expiresInSeconds: 11564869.0,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.141802025539163407 * 1e18),
			newVal:           big.NewInt(0.114669254016829994 * 1e18),
			expected:         false,
		},
		{
			name:             "test 8 (previous block) - SHOULD NOT UPDATE",
			expiresInSeconds: 3574128.999999998603016138,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.114668136842518698 * 1e18),
			newVal:           big.NewInt(0.202460645024667513 * 1e18),
			expected:         false,
		},
		{
			name:             "test 9 EDGE CASE - SHOULD NOT UPDATE",
			expiresInSeconds: 3574116.99999999860301613807678,
			thresholdPPB:     1e7,
			oldVal:           valFromString(t, "0.114668136842518697537940397524"),
			newVal:           valFromString(t, "0.202462661212593902915202193071"),
			expected:         false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			clock := frozenTimeClock{}
			expiresAt := float64(clock.Now().Unix()) + tc.expiresInSeconds
			actual, err := makePendleDeviationFunc(expiresAt, clock, DefaultMultiplier)(nil, tc.thresholdPPB, tc.oldVal, tc.newVal)
			if tc.err != "" {
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
				if actual != tc.expected {
					t.Fatalf("expected %v, got %v", tc.expected, actual)
				}
			}
		})
	}
}
