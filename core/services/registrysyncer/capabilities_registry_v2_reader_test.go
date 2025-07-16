package registrysyncer

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCapabilitiesRegistryV2Reader_Creation(t *testing.T) {
	ctx := context.Background()

	// Test that we can create a V2 reader without errors
	reader, err := NewCapabilitiesRegistryV2Reader(
		ctx,
		nil, // We'll pass nil for now to test the creation
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
	)

	require.NoError(t, err)
	require.NotNil(t, reader)

	// Verify the reader implements the interface
	assert.Implements(t, (*CapabilitiesRegistryReader)(nil), reader)

	// Verify the address is set correctly
	assert.Equal(t, common.HexToAddress("0x1234567890123456789012345678901234567890"), reader.Address())
}

func TestParseCapabilityID(t *testing.T) {
	testCases := []struct {
		name                 string
		input                string
		expectedID           string
		expectedLabelledName string
		expectedVersion      string
		shouldError          bool
	}{
		{
			name:                 "valid capability ID",
			input:                "data-streams-reports@1.0.0",
			expectedID:           "data-streams-reports@1.0.0",
			expectedLabelledName: "data-streams-reports",
			expectedVersion:      "1.0.0",
			shouldError:          false,
		},
		{
			name:                 "another valid capability ID",
			input:                "workflow-engine@2.1.0",
			expectedID:           "workflow-engine@2.1.0",
			expectedLabelledName: "workflow-engine",
			expectedVersion:      "2.1.0",
			shouldError:          false,
		},
		{
			name:        "invalid format - no version",
			input:       "data-streams-reports",
			shouldError: true,
		},
		{
			name:        "invalid format - multiple @",
			input:       "data-streams-reports@1.0.0@extra",
			shouldError: true,
		},
		{
			name:        "empty input",
			input:       "",
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, labelledName, version, err := parseCapabilityID(tc.input)

			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedID, id)
			assert.Equal(t, tc.expectedLabelledName, labelledName)
			assert.Equal(t, tc.expectedVersion, version)
		})
	}
}

func TestParseCapabilityMetadata(t *testing.T) {
	testCases := []struct {
		name             string
		input            []byte
		expectedCapType  uint8
		expectedRespType uint8
		shouldError      bool
	}{
		{
			name:             "valid metadata",
			input:            []byte(`{"capabilityType": 1, "responseType": 2}`),
			expectedCapType:  1,
			expectedRespType: 2,
			shouldError:      false,
		},
		{
			name:             "empty metadata",
			input:            []byte{},
			expectedCapType:  0,
			expectedRespType: 0,
			shouldError:      false,
		},
		{
			name:             "invalid JSON",
			input:            []byte(`{"capabilityType": 1, "responseType": 2`),
			expectedCapType:  0,
			expectedRespType: 0,
			shouldError:      false, // Should not error, just return defaults
		},
		{
			name:             "missing fields",
			input:            []byte(`{"otherField": "value"}`),
			expectedCapType:  0,
			expectedRespType: 0,
			shouldError:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			capType, respType, err := parseCapabilityMetadata(tc.input)

			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedCapType, capType)
			assert.Equal(t, tc.expectedRespType, respType)
		})
	}
}
