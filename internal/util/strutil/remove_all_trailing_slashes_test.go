package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveAllTrailingSlashes(t *testing.T) {
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
			name:     "string without slashes",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "string with single trailing slash",
			input:    "hello/",
			expected: "hello",
		},
		{
			name:     "string with multiple trailing slashes",
			input:    "hello///",
			expected: "hello",
		},
		{
			name:     "string with slashes in middle",
			input:    "hello/world/test",
			expected: "hello/world/test",
		},
		{
			name:     "string with slashes in middle and trailing",
			input:    "hello/world/test///",
			expected: "hello/world/test",
		},
		{
			name:     "only slashes",
			input:    "////",
			expected: "",
		},
		{
			name:     "string with mixed path separators",
			input:    "hello/world\\test///",
			expected: "hello/world\\test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveAllTrailingSlashes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
