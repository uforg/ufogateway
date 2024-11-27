package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveTrailingSlash(t *testing.T) {
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
			name:     "string without slash",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "string with trailing slash",
			input:    "hello/",
			expected: "hello",
		},
		{
			name:     "string with multiple slashes",
			input:    "hello//",
			expected: "hello/",
		},
		{
			name:     "string with slash in middle",
			input:    "hello/world",
			expected: "hello/world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveTrailingSlash(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
