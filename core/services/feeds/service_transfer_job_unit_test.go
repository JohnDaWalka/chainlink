package feeds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Service_TransferJob_AllowedPartners_Unit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		allowedPartners []string
		targetName      string
		expectAllowed   bool
	}{
		{
			name:            "Target in allowed list",
			allowedPartners: []string{"Target-Manager", "Other-Manager"},
			targetName:      "Target-Manager",
			expectAllowed:   true,
		},
		{
			name:            "Target not in allowed list",
			allowedPartners: []string{"Other-Manager", "Another-Manager"},
			targetName:      "Target-Manager",
			expectAllowed:   false,
		},
		{
			name:            "Empty allowed list (blocks all transfers)",
			allowedPartners: []string{},
			targetName:      "Target-Manager",
			expectAllowed:   false,
		},
		{
			name:            "Nil allowed list (blocks all transfers)",
			allowedPartners: nil,
			targetName:      "Target-Manager",
			expectAllowed:   false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			targetAllowed := false

			for _, partner := range tc.allowedPartners {
				if partner == tc.targetName {
					targetAllowed = true
					break
				}
			}

			assert.Equal(t, tc.expectAllowed, targetAllowed,
				"Expected allowed=%v for target=%s with partners=%v",
				tc.expectAllowed, tc.targetName, tc.allowedPartners)
		})
	}
}
