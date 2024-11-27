package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveLeadingSlash(t *testing.T) {
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
			name:     "no leading slash",
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
			input:    "//hello",
			expected: "/hello",
		},
		{
			name:     "slash in middle",
			input:    "hello/world",
			expected: "hello/world",
		},
		{
			name:     "only slash",
			input:    "/",
			expected: "",
		},
		{
			name:     "slash with space",
			input:    "/ hello",
			expected: " hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveLeadingSlash(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
