package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneHeaderMap(t *testing.T) {
	tests := []struct {
		name     string
		input    http.Header
		expected map[string][]string
	}{
		{
			name: "Single header",
			input: http.Header{
				"Content-Type": {"application/json"},
			},
			expected: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "Multiple headers",
			input: http.Header{
				"Content-Type": {"application/json"},
				"Accept":       {"application/json", "text/plain"},
			},
			expected: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"application/json", "text/plain"},
			},
		},
		{
			name:     "Empty header",
			input:    http.Header{},
			expected: map[string][]string{},
		},
		{
			name:     "Nil header",
			input:    nil,
			expected: map[string][]string{},
		},
		{
			name: "Headers with empty values",
			input: http.Header{
				"Empty-Header":  {""},
				"Normal-Header": {"value"},
			},
			expected: map[string][]string{
				"Empty-Header":  {""},
				"Normal-Header": {"value"},
			},
		},
		{
			name: "Headers with multiple empty values",
			input: http.Header{
				"Multi-Empty": {"", "", ""},
			},
			expected: map[string][]string{
				"Multi-Empty": {"", "", ""},
			},
		},
		{
			name: "Case sensitivity preservation",
			input: http.Header{
				"Content-Type": {"application/json"},
				"CONTENT-TYPE": {"application/xml"},
				"content-type": {"text/plain"},
			},
			expected: map[string][]string{
				"Content-Type": {"application/json"},
				"CONTENT-TYPE": {"application/xml"},
				"content-type": {"text/plain"},
			},
		},
		{
			name: "Special characters in header names",
			input: http.Header{
				"X-Special-Header-!@#$": {"value"},
				"X-Header-With-空格":      {"value"},
				"X-Header-With-😀":       {"value"},
			},
			expected: map[string][]string{
				"X-Special-Header-!@#$": {"value"},
				"X-Header-With-空格":      {"value"},
				"X-Header-With-😀":       {"value"},
			},
		},
		{
			name: "Special characters in header values",
			input: http.Header{
				"X-Special": {"value with spaces", "value,with,commas", "value\nwith\nlinebreaks"},
				"X-Unicode": {"值", "사용자", "Café"},
				"X-Emoji":   {"👍 👎 🎉", "❤️ 💔"},
			},
			expected: map[string][]string{
				"X-Special": {"value with spaces", "value,with,commas", "value\nwith\nlinebreaks"},
				"X-Unicode": {"值", "사용자", "Café"},
				"X-Emoji":   {"👍 👎 🎉", "❤️ 💔"},
			},
		},
		{
			name: "Large number of values",
			input: http.Header{
				"X-Many-Values": {"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			},
			expected: map[string][]string{
				"X-Many-Values": {"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			},
		},
		{
			name: "Empty slice as value",
			input: http.Header{
				"Empty-Slice": {},
			},
			expected: map[string][]string{
				"Empty-Slice": {},
			},
		},
		{
			name: "Mixed content types",
			input: http.Header{
				"Accept": {
					"text/html",
					"application/xhtml+xml",
					"application/xml;q=0.9",
					"image/webp",
					"*/*;q=0.8",
				},
			},
			expected: map[string][]string{
				"Accept": {
					"text/html",
					"application/xhtml+xml",
					"application/xml;q=0.9",
					"image/webp",
					"*/*;q=0.8",
				},
			},
		},
		{
			name: "Common HTTP headers",
			input: http.Header{
				"Authorization":   {"Bearer token123"},
				"Cache-Control":   {"no-cache", "no-store"},
				"If-None-Match":   {"\"123456\"", "\"789\""},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			},
			expected: map[string][]string{
				"Authorization":   {"Bearer token123"},
				"Cache-Control":   {"no-cache", "no-store"},
				"If-None-Match":   {"\"123456\"", "\"789\""},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := cloneHeaderMap(tt.input)
			assert.Equal(t, tt.expected, cloned)

			// Verify the cloned map is independent of the original
			if tt.input != nil {
				for k, v := range tt.input {
					if len(v) > 0 {
						originalValue := v[0]
						tt.input[k][0] = "modified"
						assert.NotEqual(t, "modified", cloned[k][0], "Cloned map should not be affected by modifications to original")
						assert.Equal(t, originalValue, cloned[k][0], "Cloned value should maintain original value")
					}
				}
			}

			// Verify modifying the clone doesn't affect the original
			for k, v := range cloned {
				if len(v) > 0 {
					originalValue := v[0]
					cloned[k][0] = "modified-clone"
					if tt.input != nil && len(tt.input[k]) > 0 {
						assert.NotEqual(t, "modified-clone", tt.input[k][0], "Original map should not be affected by modifications to clone")
						assert.NotEqual(t, originalValue, cloned[k][0], "Modified clone should have different value")
					}
				}
			}
		})
	}
}
