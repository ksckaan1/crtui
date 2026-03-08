package registryclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortTags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "semantic versioning",
			input:    []string{"v1.0.0", "v2.0.0", "v1.5.0", "v1.0.1"},
			expected: []string{"v2.0.0", "v1.5.0", "v1.0.1", "v1.0.0"},
		},
		{
			name:     "latest first",
			input:    []string{"v1.0.0", "latest", "v2.0.0"},
			expected: []string{"latest", "v2.0.0", "v1.0.0"},
		},
		{
			name:     "alpha numeric mixed",
			input:    []string{"abc", "123", "xyz"},
			expected: []string{"abc", "xyz", "123"},
		},
		{
			name:     "pre-release versions",
			input:    []string{"v1.0.0-rc1", "v1.0.0-beta", "v1.0.0-alpha", "v1.0.0"},
			expected: []string{"v1.0.0", "v1.0.0-rc1", "v1.0.0-beta", "v1.0.0-alpha"},
		},
		{
			name:     "with hyphens",
			input:    []string{"v1.0.0", "v1.0.0-SNAPSHOT", "v1.0.0-alpha"},
			expected: []string{"v1.0.0", "v1.0.0-SNAPSHOT", "v1.0.0-alpha"},
		},
		{
			name:     "empty and latest",
			input:    []string{"v1.0.0", "latest", ""},
			expected: []string{"latest", "", "v1.0.0"},
		},
		{
			name:     "single version",
			input:    []string{"v1.0.0"},
			expected: []string{"v1.0.0"},
		},
		{
			name:     "empty list",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "all latest",
			input:    []string{"latest", "latest"},
			expected: []string{"latest", "latest"},
		},
		{
			name:     "date based versions",
			input:    []string{"2024.01.01", "2023.12.01", "2024.06.15"},
			expected: []string{"2024.06.15", "2024.01.01", "2023.12.01"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := make([]string, len(tt.input))
			copy(sorted, tt.input)
			sortTags(sorted)
			require.Equal(t, tt.expected, sorted)
		})
	}
}

func TestHasNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"v1.0.0", true},
		{"latest", false},
		{"abc123", true},
		{"abc", false},
		{"", false},
		{"v1.0.0-rc1", true},
		{"2024.01.01", true},
	}

	for _, tt := range tests {
		result := hasNumber(tt.input)
		require.Equal(t, tt.expected, result, "input: %s", tt.input)
	}
}
