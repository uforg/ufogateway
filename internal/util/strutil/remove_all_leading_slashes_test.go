package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveAllLe6adingSlashes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no leading slashes",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "single leading slash",
			input:    "/hello",
			expected: "hello",
		},
		{
			name:     "multiple leading slashes",
			input:    "///hello",
			expected: "hello",
		},
		{
			name:     "slashes in middle and end",
			input:    "/path/to/something/",
			expected: "path/to/something/",
		},
		{
			name:     "only slashes",
			input:    "////",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveAllLeadingSlashes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
