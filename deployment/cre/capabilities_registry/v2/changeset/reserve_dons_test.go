package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestReserveDons(t *testing.T) {
	tenv := test.SetupEnvV2(t, false)
	chainSelector := tenv.RegistrySelector
	env := *tenv.Env

	tests := []struct {
		name        string
		input       changeset.ReserveDonsInput
		expectError bool
		errorMsg    string
		checkOutput func(t *testing.T, output contracts.RegisterDonsOutput)
	}{
		{
			name: "valid input with 3 DONs",
			input: changeset.ReserveDonsInput{
				Address:       tenv.RegistryAddress.String(),
				ChainSelector: chainSelector,
				N:             3,
			},
			expectError: false,
			checkOutput: func(t *testing.T, output contracts.RegisterDonsOutput) {
				require.NotNil(t, output, "output should not be nil")
				//require.NotNil(t, output.DataStore, "datastore should not be nil")
				//addresses := output.DataStore.Addresses().Filter(changeset.AddressRefByQualifierPrefix("reserved-"))
				require.Len(t, output.DONs, 3, "should have 3 reserved DON addresses")

			},
		},
		/*
			{
				name: "invalid input with N=0",
				input: changeset.ReserveDonsInput{
					Address:       tenv.RegistryAddress.String(),
					ChainSelector: chainSelector,
					N:             0,
				},
				expectError: true,
				errorMsg:    "N must be greater than 0",
			},
			{
				name: "invalid input with negative N",
				input: changeset.ReserveDonsInput{
					Address:       tenv.RegistryAddress.String(),
					ChainSelector: chainSelector,
					N:             -1,
				},
				expectError: true,
				errorMsg:    "N must be greater than 0",
			},
			{
				name: "empty address",
				input: changeset.ReserveDonsInput{
					Address:       "",
					ChainSelector: chainSelector,
					N:             2,
				},
				expectError: true,
				errorMsg:    "address is not set",
			},
			{
				name: "invalid chain selector",
				input: changeset.ReserveDonsInput{
					Address:       tenv.RegistryAddress.String(),
					ChainSelector: 999999, // Invalid chain
					N:             2,
				},
				expectError: true,
				errorMsg:    "chain 999999 not found in environment",
			},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := changeset.ReserveDons{}

			// Test VerifyPreconditions first
			err := cs.VerifyPreconditions(env, tt.input)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				return // Don't proceed to Apply if preconditions fail
			}
			require.NoError(t, err, "VerifyPreconditions should pass for valid input")

			// Test Apply for valid inputs only
			if !tt.expectError {
				// Note: This will likely fail without proper contract deployment
				// In a real test, you'd need to deploy a CapabilitiesRegistry contract first
				out, err := cs.Apply(env, tt.input)
				require.NoError(t, err, "Apply should succeed for valid input")
				require.NotNil(t, out, "output should not be nil")
				require.Len(t, out.Reports, 1, "should have one report, corresponding to RegisterDons")
				x := (out.Reports[0].Output).(contracts.RegisterDonsOutput)
				tt.checkOutput(t, x)

				// For now, we just verify the method can be called
				// In a complete test, you'd verify the output structure
				t.Logf("Apply result: err=%v", err)
			}
		})
	}
}

/*
func TestReserveDonsUniqueNames(t *testing.T) {
	tenv := test.SetupEnvV2(t, false)
	chainSelector := tenv.RegistrySelector
	env := *tenv.Env

	changeset := changeset.ReserveDons{}

	input1 := changeset.ReserveDonsInput{
		Address:       "0x1234567890123456789012345678901234567890",
		ChainSelector: chainSelector,
		N:             3,
	}

	input2 := changeset.ReserveDonsInput{
		Address:       "0x1234567890123456789012345678901234567890",
		ChainSelector: chainSelector,
		N:             2,
	}

	t.Run("unique names across calls", func(t *testing.T) {
		// Verify preconditions first
		err := changeset.VerifyPreconditions(env, input1)
		require.NoError(t, err)

		err = changeset.VerifyPreconditions(env, input2)
		require.NoError(t, err)

		// Note: For a complete test, you would:
		// 1. Deploy a CapabilitiesRegistry contract first
		// 2. Apply both changesets
		// 3. Verify the reserved DON names are unique and follow the "reserved-" pattern

		t.Log("Preconditions verified for both inputs")
		t.Log("Complete test would require contract deployment and verification of unique reserved names")
	})
}

func TestReserveDonsVerifyPreconditions(t *testing.T) {
	tenv := test.SetupEnvV2(t, false)
	chainSelector := tenv.RegistrySelector
	env := *tenv.Env
	changeset := changeset.ReserveDons{}

	t.Run("all preconditions pass", func(t *testing.T) {
		input := changeset.ReserveDonsInput{
			Address:       "0x1234567890123456789012345678901234567890",
			ChainSelector: chainSelector,
			N:             5,
		}

		err := changeset.VerifyPreconditions(env, input)
		require.NoError(t, err)
	})

	t.Run("missing address fails", func(t *testing.T) {
		input := changeset.ReserveDonsInput{
			Address:       "",
			ChainSelector: chainSelector,
			N:             5,
		}

		err := changeset.VerifyPreconditions(env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "address is not set")
	})

	t.Run("zero N fails", func(t *testing.T) {
		input := changeset.ReserveDonsInput{
			Address:       "0x1234567890123456789012345678901234567890",
			ChainSelector: chainSelector,
			N:             0,
		}

		err := changeset.VerifyPreconditions(env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "N must be greater than 0")
	})

	t.Run("invalid chain selector fails", func(t *testing.T) {
		input := changeset.ReserveDonsInput{
			Address:       "0x1234567890123456789012345678901234567890",
			ChainSelector: 999999,
			N:             3,
		}

		err := changeset.VerifyPreconditions(env, input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain 999999 not found in environment")
	})
}
*/
