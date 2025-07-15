package registrysyncer

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestCapabilitiesRegistryV1Reader_Creation(t *testing.T) {
	ctx := context.Background()

	// Test that we can create a V1 reader without errors
	reader, err := NewCapabilitiesRegistryV1Reader(
		ctx,
		nil, // We'll pass nil for now to test the creation
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
	)

	require.NoError(t, err)
	require.NotNil(t, reader)
}

func TestCapabilitiesRegistryReaderFactory_Creation(t *testing.T) {
	factory := NewCapabilitiesRegistryReaderFactory()
	require.NotNil(t, factory)
}
