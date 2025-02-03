package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLabelSet(t *testing.T) {
	t.Run("no labels", func(t *testing.T) {
		ls := NewLabelSet()
		assert.Empty(t, ls, "expected empty set")
	})

	t.Run("some labels", func(t *testing.T) {
		ls := NewLabelSet("foo", "bar")
		assert.Len(t, ls, 2)
		assert.True(t, ls.Contains("foo"))
		assert.True(t, ls.Contains("bar"))
		assert.False(t, ls.Contains("baz"))
	})
}

func TestLabelSet_Add(t *testing.T) {
	ls := NewLabelSet("initial")
	ls.Add("new")

	assert.True(t, ls.Contains("initial"), "expected 'initial' in set")
	assert.True(t, ls.Contains("new"), "expected 'new' in set")
	assert.Len(t, ls, 2, "expected 2 distinct labels in set")

	// Add duplicate "new" again; size should remain 2
	ls.Add("new")
	assert.Len(t, ls, 2, "expected size to remain 2 after adding a duplicate")

	// Add multiple labels at once
	ls.Add("label1", "label2", "label3")
	assert.Len(t, ls, 5, "expected 5 distinct labels in set") // 2 previous + 3 new
	assert.True(t, ls.Contains("label1"))
	assert.True(t, ls.Contains("label2"))
	assert.True(t, ls.Contains("label3"))
}

func TestLabelSet_Remove(t *testing.T) {
	ls := NewLabelSet("remove_me", "keep", "label1", "label2", "label3")
	ls.Remove("remove_me")

	assert.False(t, ls.Contains("remove_me"), "expected 'remove_me' to be removed")
	assert.True(t, ls.Contains("keep"), "expected 'keep' to remain")
	assert.True(t, ls.Contains("label1"), "expected 'label1' to remain")
	assert.True(t, ls.Contains("label2"), "expected 'label2' to remain")
	assert.True(t, ls.Contains("label3"), "expected 'label3' to remain")
	assert.Len(t, ls, 4, "expected set size to be 4 after removal")

	// Removing a non-existent item shouldn't change the size
	ls.Remove("non_existent")
	assert.Len(t, ls, 4, "expected size to remain 4 after removing a non-existent item")

	// Remove multiple labels at once
	ls.Remove("label2", "label4")

	assert.Len(t, ls, 3, "expected 3 distinct labels in set after removal") // keep, label1, label3
	assert.True(t, ls.Contains("keep"))
	assert.True(t, ls.Contains("label1"))
	assert.False(t, ls.Contains("label2"))
	assert.True(t, ls.Contains("label3"))
	assert.False(t, ls.Contains("label4"))
}

func TestLabelSet_Contains(t *testing.T) {
	ls := NewLabelSet("foo", "bar")

	assert.True(t, ls.Contains("foo"))
	assert.True(t, ls.Contains("bar"))
	assert.False(t, ls.Contains("baz"))
}

// TestLabelSet_String tests the String() method of the LabelSet type.
func TestLabelSet_String(t *testing.T) {
	tests := []struct {
		name     string
		labels   LabelSet
		expected string
	}{
		{
			name:     "Empty LabelSet",
			labels:   NewLabelSet(),
			expected: "",
		},
		{
			name:     "Single label",
			labels:   NewLabelSet("alpha"),
			expected: "alpha",
		},
		{
			name:     "Multiple labels in random order",
			labels:   NewLabelSet("beta", "gamma", "alpha"),
			expected: "alpha beta gamma",
		},
		{
			name:     "Labels with special characters",
			labels:   NewLabelSet("beta", "gamma!", "@alpha"),
			expected: "@alpha beta gamma!",
		},
		{
			name:     "Labels with spaces",
			labels:   NewLabelSet("beta", "gamma delta", "alpha"),
			expected: "alpha beta gamma delta",
		},
		{
			name:     "Labels added in different orders",
			labels:   NewLabelSet("delta", "beta", "alpha"),
			expected: "alpha beta delta",
		},
		{
			name:     "Labels with duplicate additions",
			labels:   NewLabelSet("alpha", "beta", "alpha", "gamma", "beta"),
			expected: "alpha beta gamma",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := tt.labels.String()
			assert.Equal(t, tt.expected, result, "LabelSet.String() should return the expected sorted string")
		})
	}
}

func TestLabelSet_Equal(t *testing.T) {
	tests := []struct {
		name     string
		set1     LabelSet
		set2     LabelSet
		expected bool
	}{
		{
			name:     "Both sets empty",
			set1:     NewLabelSet(),
			set2:     NewLabelSet(),
			expected: true,
		},
		{
			name:     "First set empty, second set non-empty",
			set1:     NewLabelSet(),
			set2:     NewLabelSet("label1"),
			expected: false,
		},
		{
			name:     "First set non-empty, second set empty",
			set1:     NewLabelSet("label1"),
			set2:     NewLabelSet(),
			expected: false,
		},
		{
			name:     "Identical sets with single label",
			set1:     NewLabelSet("label1"),
			set2:     NewLabelSet("label1"),
			expected: true,
		},
		{
			name:     "Identical sets with multiple labels",
			set1:     NewLabelSet("label1", "label2", "label3"),
			set2:     NewLabelSet("label3", "label2", "label1"), // Different order
			expected: true,
		},
		{
			name:     "Different sets, same size",
			set1:     NewLabelSet("label1", "label2", "label3"),
			set2:     NewLabelSet("label1", "label2", "label4"),
			expected: false,
		},
		{
			name:     "Different sets, different sizes",
			set1:     NewLabelSet("label1", "label2"),
			set2:     NewLabelSet("label1", "label2", "label3"),
			expected: false,
		},
		{
			name:     "Subset sets",
			set1:     NewLabelSet("label1", "label2"),
			set2:     NewLabelSet("label1", "label2", "label3"),
			expected: false,
		},
		{
			name:     "Disjoint sets",
			set1:     NewLabelSet("label1", "label2"),
			set2:     NewLabelSet("label3", "label4"),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.set1.Equal(tt.set2)
			assert.Equal(t, tt.expected, result, "Equal(%v, %v) should be %v", tt.set1, tt.set2, tt.expected)
		})
	}
}

func TestLabelSet_DeepClone(t *testing.T) {
	tests := []struct {
		name      string
		input     LabelSet
		mutate    func(ls LabelSet)
		wantEqual bool
	}{
		{
			name: "Empty label set",
			// An empty-but-initialized map
			input: NewLabelSet(),
			mutate: func(ls LabelSet) {
				ls.Add("test-label")
			},
			wantEqual: false,
		},
		{
			name:  "Non-empty label set",
			input: NewLabelSet("fast", "secure"),
			mutate: func(ls LabelSet) {
				ls.Remove("secure") // remove an existing label
				ls.Add("extra")     // add a new label
			},
			wantEqual: false,
		},
		{
			name:  "Clone but no mutation",
			input: NewLabelSet("unchanged"),
			mutate: func(ls LabelSet) {
				// no change
			},
			// If we do no mutation, they should remain equal
			wantEqual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.input
			clone := original.DeepClone()

			// Clones should initially match the original (especially for non-nil).
			// For nil, clone should also be nil, so they're logically "equal" if you treat nil sets as empty.
			assert.True(t, original.Equal(clone),
				"DeepClone result should match input before mutation")

			// Mutate the clone
			tt.mutate(clone)

			// After mutation, check if they should remain equal or differ.
			if tt.wantEqual {
				assert.True(t, original.Equal(clone),
					"Unexpected difference between original and clone")
			} else {
				assert.False(t, original.Equal(clone),
					"Expected original and clone to differ after mutation, but they are equal")
			}

			// Optionally, also check that mutating the original won't affect the clone.
			if original != nil {
				original.Add("original-mutated")
				// They should differ if we just mutated original (unless they were already different).
				// We'll just do a quick check that the new label isn't in the clone.
				if original.Contains("original-mutated") && clone != nil {
					assert.False(t, clone.Contains("original-mutated"),
						"Clone should not contain a label added to the original after deep clone")
				}
			}
		})
	}
}
