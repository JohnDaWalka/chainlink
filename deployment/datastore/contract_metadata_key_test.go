package datastore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContractMetadataKey_Equals(t *testing.T) {
	tests := []struct {
		name     string
		key1     ContractMetadataKey
		key2     ContractMetadataKey
		expected bool
	}{
		{
			name:     "Equal keys",
			key1:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			key2:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			expected: true,
		},
		{
			name:     "Different chain",
			key1:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			key2:     NewContractMetadataKey(2, "0x1234567890abcdef"),
			expected: false,
		},
		{
			name:     "Different address",
			key1:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			key2:     NewContractMetadataKey(1, "0xabcdef1234567890"),
			expected: false,
		},
		{
			name:     "Completely different keys",
			key1:     NewContractMetadataKey(1, "0x1234567890abcdef"),
			key2:     NewContractMetadataKey(2, "0xabcdef1234567890"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.key1.Equals(tt.key2))
		})
	}
}

func TestContractMetadataKey(t *testing.T) {
	chain := uint64(1)
	address := "0x1234567890abcdef"

	key := NewContractMetadataKey(chain, address)

	require.Equal(t, chain, key.ChainSelector(), "ChainSelector should return the correct chain ID")
	require.Equal(t, address, key.Address(), "Address should return the correct address")
}
