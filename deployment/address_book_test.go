package deployment

import (
	"math/big"
	"sync"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
)

func TestTypeAndVersion_NewTypeAndVersion(t *testing.T) {
	contractType := ContractType("TestContract")
	version := semver.MustParse("1.0.0")

	tv1 := NewTypeAndVersion(contractType, *version)
	tv2 := NewTypeAndVersion(contractType, *version)
	tv3 := TypeAndVersion{
		Type:    "TestContract",
		Version: *version,
	}

	assert.True(t, tv1.Equal(tv2), "expected tv1 to be equal to tv2")
	assert.Equal(t, tv1, tv3, "expected tv1 to be equal to tv3")
}

func TestTypeAndVersion_String(t *testing.T) {
	contractType := ContractType("TestContract")
	version := semver.MustParse("1.0.0")

	tests := []struct {
		name     string
		tv       TypeAndVersion
		expected string
	}{
		{
			name: "Nil labels",
			tv: TypeAndVersion{
				Type:    contractType,
				Version: *version,
				Labels:  nil,
			},
			expected: "TestContract 1.0.0",
		},
		{
			name: "Empty labels",
			tv: TypeAndVersion{
				Type:    contractType,
				Version: *version,
				Labels:  make(LabelSet),
			},
			expected: "TestContract 1.0.0",
		},
		{
			name: "With labels",
			tv: TypeAndVersion{
				Type:    contractType,
				Version: *version,
				Labels:  NewLabelSet("alpha", "beta"),
			},
			expected: "TestContract 1.0.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.tv.String(), "unexpected string representation")
		})
	}
}

func TestTypeAndVersion_FullString(t *testing.T) {
	contractType := ContractType("TestContract")
	version := semver.MustParse("1.0.0")

	tests := []struct {
		name     string
		tv       TypeAndVersion
		expected string
	}{
		{
			name: "Nil labels",
			tv: TypeAndVersion{
				Type:    contractType,
				Version: *version,
				Labels:  nil,
			},
			expected: "TestContract 1.0.0",
		},
		{
			name: "Empty labels",
			tv: TypeAndVersion{
				Type:    contractType,
				Version: *version,
				Labels:  make(LabelSet),
			},
			expected: "TestContract 1.0.0",
		},
		{
			name: "With labels",
			tv: TypeAndVersion{
				Type:    contractType,
				Version: *version,
				Labels:  NewLabelSet("alpha", "beta"),
			},
			expected: "TestContract 1.0.0 alpha beta",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.tv.FullString(), "unexpected string representation")
		})
	}
}

func TestTypeAndVersion_DeepClone(t *testing.T) {
	tests := []struct {
		name      string
		input     TypeAndVersion
		mutate    func(tv *TypeAndVersion)
		wantEqual bool
	}{
		{
			name: "No labels",
			input: TypeAndVersion{
				Type:    "MyContract",
				Version: *semver.MustParse("1.2.3"),
				Labels:  nil,
			},
			mutate: func(tv *TypeAndVersion) {
				tv.Type = "Mutated"
				tv.Version = *semver.MustParse("9.9.9")
			},
			wantEqual: false,
		},
		{
			name: "With labels",
			input: TypeAndVersion{
				Type:    "AnotherContract",
				Version: *semver.MustParse("2.0.1"),
				Labels:  NewLabelSet("fast", "secure"),
			},
			mutate: func(tv *TypeAndVersion) {
				tv.Labels.Add("new-label")
			},
			wantEqual: false,
		},
		{
			name: "Empty label set",
			input: TypeAndVersion{
				Type:    "EmptyLabelContract",
				Version: *semver.MustParse("0.1.0"),
				Labels:  NewLabelSet(), // empty, but allocated
			},
			mutate: func(tv *TypeAndVersion) {
				tv.Labels.Add("test-label")
			},
			wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone
			clone := tt.input.DeepClone()

			// Before mutation, the clone should be Equal to the original
			assert.True(t, tt.input.Equal(clone),
				"DeepClone result should initially match the input")

			// Mutate the clone
			tt.mutate(&clone)

			// If wantEqual is false, the original should differ from the mutated clone
			if !tt.wantEqual {
				assert.False(t, tt.input.Equal(clone),
					"Mutating the clone should not affect the original if deep-cloned")
			} else {
				assert.True(t, tt.input.Equal(clone),
					"Mutating the clone incorrectly affected the original")
			}
		})
	}
}

func TestAddressBookMap_DeepCloneAddresses(t *testing.T) {
	// Prepare some TypeAndVersion items
	tvA := TypeAndVersion{
		Type:    "ContractA",
		Version: *semver.MustParse("1.0.0"),
		Labels:  NewLabelSet("labelA"),
	}
	tvB := TypeAndVersion{
		Type:    "ContractB",
		Version: *semver.MustParse("1.1.0"),
		Labels:  NewLabelSet("labelB1", "labelB2"),
	}

	// Build our sample input
	inputMap := map[uint64]map[string]TypeAndVersion{
		111: {
			"0x1234": tvA,
		},
		222: {
			"0xABCD": tvB,
		},
	}

	ab := NewMemoryAddressBookFromMap(inputMap)

	// Addresses() is supposed to return a deep clone
	clonedAddrs, err := ab.Addresses()
	require.NoError(t, err)

	// Now mutate something in the clone to see if the original is affected
	clonedAddrs[111]["0x1234"] = TypeAndVersion{
		Type:    "MutatedType",
		Version: *semver.MustParse("9.9.9"),
		Labels:  NewLabelSet("mutated"),
	}

	// Check original is not mutated
	originalAddrs, err := ab.Addresses()
	require.NoError(t, err)

	// The original 111 -> 0x1234 should remain tvA
	assert.Equal(t, tvA, originalAddrs[111]["0x1234"],
		"Mutating cloned addresses must not affect the original")

	// Also check that the `Labels` inside each TypeAndVersion are deeply cloned
	// For example, add a label to the clone's tvB
	cloneTvB := clonedAddrs[222]["0xABCD"]
	cloneTvB.Labels.Add("extra-label")
	clonedAddrs[222]["0xABCD"] = cloneTvB

	// Now see if the original's version is unchanged
	originalTvB := originalAddrs[222]["0xABCD"]
	assert.False(t, originalTvB.Labels.Contains("extra-label"),
		"Original TypeAndVersion's Labels should not reflect changes to the clone")

	// Optionally, ensure the rest of the original is still correct
	assert.Equal(t, tvB, originalTvB,
		"Original TypeAndVersion for 222 -> 0xABCD should remain unchanged")
}

func TestAddressBook_Save(t *testing.T) {
	ab := NewMemoryAddressBook()
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()

	err := ab.Save(chainsel.TEST_90000001.Selector, addr1, onRamp100)
	require.NoError(t, err)

	// Invalid address
	err = ab.Save(chainsel.TEST_90000001.Selector, "asdlfkj", onRamp100)
	require.ErrorIs(t, err, ErrInvalidAddress)

	// Valid chain but not present.
	_, err = ab.AddressesForChain(chainsel.TEST_90000002.Selector)
	require.ErrorIs(t, err, ErrChainNotFound)

	// Invalid selector
	err = ab.Save(0, addr1, onRamp100)
	require.ErrorIs(t, err, ErrInvalidChainSelector)

	// Duplicate
	err = ab.Save(chainsel.TEST_90000001.Selector, addr1, onRamp100)
	require.Error(t, err)

	// Zero address
	err = ab.Save(chainsel.TEST_90000001.Selector, common.HexToAddress("0x0").Hex(), onRamp100)
	require.Error(t, err)

	// Zero address but non evm chain
	err = NewMemoryAddressBook().Save(chainsel.APTOS_MAINNET.Selector, common.HexToAddress("0x0").Hex(), onRamp100)
	require.NoError(t, err)

	// Distinct address same TV will not
	err = ab.Save(chainsel.TEST_90000001.Selector, addr2, onRamp100)
	require.NoError(t, err)
	// Same address different chain will not error
	err = ab.Save(chainsel.TEST_90000002.Selector, addr1, onRamp100)
	require.NoError(t, err)
	// We can save different versions of the same contract
	err = ab.Save(chainsel.TEST_90000002.Selector, addr2, onRamp110)
	require.NoError(t, err)

	addresses, err := ab.Addresses()
	require.NoError(t, err)
	assert.Equal(t, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp100,
			addr2: onRamp110,
		},
	}, addresses)
}

func TestAddressBook_Merge(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	a1 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
	})
	a2 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	})
	require.NoError(t, a1.Merge(a2))

	addresses, err := a1.Addresses()
	require.NoError(t, err)
	assert.Equal(t, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	}, addresses)

	// Merge with conflicting addresses should error
	a3 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
	})
	require.Error(t, a1.Merge(a3))
	// a1 should not have changed
	addresses, err = a1.Addresses()
	require.NoError(t, err)
	assert.Equal(t, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	}, addresses)
}

func TestAddressBook_Remove(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	addr3 := common.HexToAddress("0x3").String()

	baseAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
			addr3: onRamp110,
		},
	})

	copyOfBaseAB := NewMemoryAddressBookFromMap(baseAB.cloneAddresses(baseAB.addressesByChain))

	// this address book shouldn't be removed (state of baseAB not changed, error thrown)
	failAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr3: onRamp100, // doesn't exist in TEST_90000001.Selector
		},
	})
	require.Error(t, baseAB.Remove(failAB))
	require.EqualValues(t, baseAB, copyOfBaseAB)

	// this Address book should be removed without error
	successAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000002.Selector: {
			addr3: onRamp100,
		},
		chainsel.TEST_90000001.Selector: {
			addr2: onRamp100,
		},
	})

	expectingAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110},
	})

	require.NoError(t, baseAB.Remove(successAB))
	require.EqualValues(t, baseAB, expectingAB)
}

func TestAddressBook_ConcurrencyAndDeadlock(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)

	baseAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			common.BigToAddress(big.NewInt(1)).String(): onRamp100,
		},
	})

	// concurrent writes
	var i int64
	wg := sync.WaitGroup{}
	for i = 2; i < 1000; i++ {
		wg.Add(1)
		go func(input int64) {
			assert.NoError(t, baseAB.Save(
				chainsel.TEST_90000001.Selector,
				common.BigToAddress(big.NewInt(input)).String(),
				onRamp100,
			))
			wg.Done()
		}(i)
	}

	// concurrent reads
	for i = 0; i < 100; i++ {
		wg.Add(1)
		go func(input int64) {
			addresses, err := baseAB.Addresses()
			if !assert.NoError(t, err) {
				return
			}
			for chainSelector, chainAddresses := range addresses {
				// concurrent read chainAddresses from Addresses() method
				for address := range chainAddresses {
					addresses[chainSelector][address] = onRamp110
				}

				// concurrent read chainAddresses from AddressesForChain() method
				chainAddresses, err = baseAB.AddressesForChain(chainSelector)
				if assert.NoError(t, err) {
					for address := range chainAddresses {
						_ = addresses[chainSelector][address]
					}
				}
			}
			wg.Done()
		}(i)
	}

	// concurrent merges, starts from 1001 to avoid address conflicts
	for i = 1001; i < 1100; i++ {
		wg.Add(1)
		go func(input int64) {
			// concurrent merge
			additionalAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
				chainsel.TEST_90000002.Selector: {
					common.BigToAddress(big.NewInt(input)).String(): onRamp100,
				},
			})
			assert.NoError(t, baseAB.Merge(additionalAB))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestAddressesContainBundle(t *testing.T) {
	t.Parallel()

	// Define some TypeAndVersion values
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	onRamp120 := NewTypeAndVersion("OnRamp", Version1_2_0)

	// Create one with labels
	onRamp100WithLabels := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp100WithLabels.AddLabel("sa")
	onRamp100WithLabels.AddLabel("staging")

	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	addr3 := common.HexToAddress("0x3").String()

	tests := []struct {
		name       string
		addrs      map[string]TypeAndVersion // input address map
		wantTypes  []TypeAndVersion          // the "bundle" we want
		wantErr    bool
		wantErrMsg string
		wantResult bool // expected boolean return when no error
	}{
		{
			name: "More than one instance => error",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100,
				addr2: onRamp100, // duplicate
			},
			wantTypes: []TypeAndVersion{onRamp100},
			wantErr:   true,
			// an example substring check:
			wantErrMsg: "found more than one instance of contract",
		},
		{
			name: "No instance => result false, no error",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp110,
				addr2: onRamp110,
			},
			wantTypes:  []TypeAndVersion{onRamp100},
			wantErr:    false,
			wantResult: false,
		},
		{
			name: "2 elements => success",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100,
				addr2: onRamp110,
				addr3: onRamp120,
			},
			wantTypes:  []TypeAndVersion{onRamp100, onRamp110},
			wantErr:    false,
			wantResult: true,
		},
		{
			name: "Mismatched labels => false",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100, // no labels
			},
			wantTypes:  []TypeAndVersion{onRamp100WithLabels},
			wantErr:    false,
			wantResult: false, // label mismatch => not found
		},
		{
			name: "Exact label match => success",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100WithLabels,
			},
			wantTypes:  []TypeAndVersion{onRamp100WithLabels},
			wantErr:    false,
			wantResult: true,
		},
		{
			name: "Duplicate labeled => error",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100WithLabels,
				addr2: onRamp100WithLabels, // same type/version/labels => duplicate
			},
			wantTypes:  []TypeAndVersion{onRamp100WithLabels},
			wantErr:    true,
			wantErrMsg: "more than one instance of contract",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotResult, gotErr := AddressesContainBundle(tt.addrs, tt.wantTypes)

			if tt.wantErr {
				require.Error(t, gotErr, "expected an error but got none")
				if tt.wantErrMsg != "" {
					require.Contains(t, gotErr.Error(), tt.wantErrMsg)
				}
				return
			}
			require.NoError(t, gotErr, "did not expect an error but got one")
			assert.Equal(t, tt.wantResult, gotResult,
				"expected result %v but got %v", tt.wantResult, gotResult)
		})
	}
}

func TestTypeAndVersionFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		input              string
		wantErr            bool
		wantType           ContractType
		wantVersion        semver.Version
		wantTypeAndVersion string
	}{
		{
			name:               "valid - no labels",
			input:              "CallProxy 1.0.0",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantTypeAndVersion: "CallProxy 1.0.0",
		},
		{
			name:               "valid - multiple labels, normal spacing",
			input:              "CallProxy 1.0.0 SA staging",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantTypeAndVersion: "CallProxy 1.0.0",
		},
		{
			name:               "valid - multiple labels, extra spacing",
			input:              "   CallProxy     1.0.0    SA    staging   ",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantTypeAndVersion: "CallProxy 1.0.0",
		},
		{
			name:    "invalid - not enough parts",
			input:   "CallProxy",
			wantErr: true,
		},
		{
			name:    "invalid - version not parseable",
			input:   "CallProxy notASemver",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotTV, gotErr := TypeAndVersionFromString(tt.input)
			if tt.wantErr {
				require.Error(t, gotErr, "expected error but got none")
				return
			}
			require.NoError(t, gotErr, "did not expect an error but got one")

			// Check ContractType
			require.Equal(t, tt.wantType, gotTV.Type, "incorrect contract type")

			// Check Version
			require.Equal(t, tt.wantVersion.String(), gotTV.Version.String(), "incorrect version")

			// Check labels
			require.Equal(t, LabelSet(nil), gotTV.Labels, "labels mismatch")

			// Check type and version
			require.Equal(t, tt.wantTypeAndVersion, gotTV.String(), "type and version mismatch")
		})
	}
}

func TestTypeAndVersionFromFullString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		input              string
		wantErr            bool
		wantType           ContractType
		wantVersion        semver.Version
		wantLabels         LabelSet
		wantTypeAndVersion string
	}{
		{
			name:               "valid - no labels",
			input:              "CallProxy 1.0.0",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantLabels:         nil,
			wantTypeAndVersion: "CallProxy 1.0.0",
		},
		{
			name:               "valid - multiple labels, normal spacing",
			input:              "CallProxy 1.0.0 SA staging",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantLabels:         NewLabelSet("SA", "staging"),
			wantTypeAndVersion: "CallProxy 1.0.0 SA staging",
		},
		{
			name:               "valid - multiple labels, extra spacing",
			input:              "   CallProxy     1.0.0    SA    staging   ",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantLabels:         NewLabelSet("SA", "staging"),
			wantTypeAndVersion: "CallProxy 1.0.0 SA staging",
		},
		{
			name:    "invalid - not enough parts",
			input:   "CallProxy",
			wantErr: true,
		},
		{
			name:    "invalid - version not parseable",
			input:   "CallProxy notASemver",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotTV, gotErr := TypeAndVersionFromFullString(tt.input)
			if tt.wantErr {
				require.Error(t, gotErr, "expected error but got none")
				return
			}
			require.NoError(t, gotErr, "did not expect an error but got one")

			// Check ContractType
			require.Equal(t, tt.wantType, gotTV.Type, "incorrect contract type")

			// Check Version
			require.Equal(t, tt.wantVersion.String(), gotTV.Version.String(), "incorrect version")

			// Check labels
			require.Equal(t, tt.wantLabels, gotTV.Labels, "labels mismatch")

			// Check full type + version + labels
			require.Equal(t, tt.wantTypeAndVersion, gotTV.FullString(), "type and version mismatch")
		})
	}
}

func TestTypeAndVersion_AddLabels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		initialLabels []string
		toAdd         []string
		wantContains  []string
		wantLen       int
	}{
		{
			name:          "add single labels to empty set",
			initialLabels: nil,
			toAdd:         []string{"foo"},
			wantContains:  []string{"foo"},
			wantLen:       1,
		},
		{
			name:          "add multiple labels to existing set",
			initialLabels: []string{"alpha"},
			toAdd:         []string{"beta", "gamma"},
			wantContains:  []string{"alpha", "beta", "gamma"},
			wantLen:       3,
		},
		{
			name:          "add duplicate labels",
			initialLabels: []string{"dup"},
			toAdd:         []string{"dup", "dup", "new"},
			wantContains:  []string{"dup", "new"},
			wantLen:       2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Construct a TypeAndVersion with any initial labels
			tv := TypeAndVersion{
				Type:    "CallProxy",
				Version: Version1_0_0,
				Labels:  NewLabelSet(tt.initialLabels...),
			}

			// Call AddLabel for each item in toAdd
			for _, label := range tt.toAdd {
				tv.AddLabel(label)
			}

			// Check final labels length
			require.Len(t, tv.Labels, tt.wantLen, "labels size mismatch")

			// Check that expected labels is present
			for _, md := range tt.wantContains {
				require.True(t, tv.Labels.Contains(md),
					"expected labels %q was not found in tv.Labels", md)
			}
		})
	}
}

func TestTypeAndVersion_RemoveLabel(t *testing.T) {
	contractType := ContractType("TestContract")
	version := semver.MustParse("1.0.0")

	tests := []struct {
		name          string
		initialLabels []string
		toRemove      []string
		wantLabels    LabelSet
	}{
		{
			name:          "Remove from nil labels",
			initialLabels: nil,
			toRemove:      []string{"alpha"},
			wantLabels:    NewLabelSet(),
		},
		{
			name:          "Remove from empty labels",
			initialLabels: []string{},
			toRemove:      []string{"alpha"},
			wantLabels:    NewLabelSet(),
		},
		{
			name:          "Remove existing label",
			initialLabels: []string{"alpha", "beta"},
			toRemove:      []string{"alpha"},
			wantLabels:    NewLabelSet("beta"),
		},
		{
			name:          "Remove non-existing label",
			initialLabels: []string{"alpha"},
			toRemove:      []string{"beta"},
			wantLabels:    NewLabelSet("alpha"),
		},
		{
			name:          "Remove multiple labels",
			initialLabels: []string{"alpha", "beta", "gamma"},
			toRemove:      []string{"alpha", "gamma"},
			wantLabels:    NewLabelSet("beta"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tv := TypeAndVersion{
				Type:    contractType,
				Version: *version,
			}

			if tt.initialLabels != nil {
				tv.Labels = NewLabelSet(tt.initialLabels...)
			}

			tv.RemoveLabel(tt.toRemove...)

			assert.True(t, tt.wantLabels.Equal(tv.Labels), "unexpected labels after removal")
		})
	}
}

func TestTypeAndVersion_LabelsString(t *testing.T) {
	contractType := ContractType("TestContract")
	version := semver.MustParse("1.0.0")

	tests := []struct {
		name     string
		labels   []string
		expected string
	}{
		{
			name:     "Nil labels",
			labels:   nil,
			expected: "",
		},
		{
			name:     "Empty labels",
			labels:   []string{},
			expected: "",
		},
		{
			name:     "Single label",
			labels:   []string{"alpha"},
			expected: "alpha",
		},
		{
			name:     "Multiple labels",
			labels:   []string{"alpha", "beta"},
			expected: "alpha beta",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tv := TypeAndVersion{
				Type:    contractType,
				Version: *version,
				Labels:  NewLabelSet(tt.labels...),
			}

			assert.Equal(t, tt.expected, tv.LabelsString(), "unexpected labels string")
		})
	}
}
