package medianpoc

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type frozenTimeClock struct{}

func (frozenTimeClock) Now() time.Time {
	return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
}

func Test_PendleDeviationFunc(t *testing.T) {
	defaultMultiplier, ok := new(big.Int).SetString("10000000000000000000", 10) // 10e18
	require.True(t, ok)

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
			name:             "test 0 (one block after) - DID UPDATE",
			expiresInSeconds: 13857541.0,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.187152977881070687 * 10e18),
			newVal:           big.NewInt(0.164498448931278907 * 10e18),
			expected:         true,
		},
		{
			name:             "test 1 (same block) - DID UPDATE",
			expiresInSeconds: 1.3857552999999971712 * 10e7,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.187152977881070687 * 10e18),
			newVal:           big.NewInt(0.164513876918347290 * 10e18),
			expected:         true,
		},
		{
			name:             "test 5 (previous block) - DID NOT UPDATE",
			expiresInSeconds: 13857565.0,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.187152977881070687 * 10e18),
			newVal:           big.NewInt(0.164529304905415591 * 10e18),
			expected:         false,
		},
		{
			name:             "test 6 (previous block) - DID NOT UPDATE",
			expiresInSeconds: 13825932.999999998137354851,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.164498448931278907 * 10e18),
			newVal:           big.NewInt(0.141815210766795902 * 10e18),
			expected:         false,
		},
		{
			name:             "test 7 (previous block) - DID NOT UPDATE",
			expiresInSeconds: 11564869.0,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.141802025539163407 * 10e18),
			newVal:           big.NewInt(0.114669254016829994 * 10e18),
			expected:         false,
		},
		{
			name:             "test 8 (previous block) - DID NOT UPDATE",
			expiresInSeconds: 3574128.999999998603016138,
			thresholdPPB:     1e7,
			oldVal:           big.NewInt(0.114668136842518698 * 10e18),
			newVal:           big.NewInt(0.202460645024667513 * 10e18),
			expected:         false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			clock := frozenTimeClock{}
			expiresAt := float64(clock.Now().Unix()) + tc.expiresInSeconds
			actual, err := makePendleDeviationFunc(expiresAt, clock, defaultMultiplier)(nil, tc.thresholdPPB, tc.oldVal, tc.newVal)
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
