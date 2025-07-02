package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobDistributorConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	p := cfg.JobDistributor()
	assert.Equal(t, "test-node", p.DisplayName())

	transfers := p.AllowedJobTransfers()
	require.Len(t, transfers, 3)

	assert.Equal(t, "54227538d9352e0a24550a80ab6a7af6e4f1ffbb8a604e913cbb81c484a7f97d", transfers[0].From)
	assert.Equal(t, "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", transfers[0].To)

	assert.Equal(t, "37346b7ea98af21e1309847e00f772826ac3689fe990b1920d01efc58ad2f250", transfers[1].From)
	assert.Equal(t, "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", transfers[1].To)

	assert.Equal(t, "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", transfers[2].From)
	assert.Equal(t, "54227538d9352e0a24550a80ab6a7af6e4f1ffbb8a604e913cbb81c484a7f97d", transfers[2].To)
}

func TestJobDistributorConfigEmptyAllowedTransfers(t *testing.T) {
	tomlString := `
[JobDistributor]
DisplayName = 'test-node'
`
	opts := GeneralConfigOpts{
		ConfigStrings: []string{tomlString},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	p := cfg.JobDistributor()
	assert.Equal(t, "test-node", p.DisplayName())
	assert.Empty(t, p.AllowedJobTransfers())
}
