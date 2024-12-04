package gateway

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriter(t *testing.T) {
	tests := []struct {
		name     string
		writeStr string
	}{
		{
			name:     "empty response",
			writeStr: "",
		},
		{
			name:     "simple string",
			writeStr: "Hello, World!",
		},
		{
			name:     "multi-write response",
			writeStr: "Part1Part2Part3",
		},
		{
			name:     "response with special characters",
			writeStr: "Hello, 世界! 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writer := newResponseWriter(rec)

			n, err := writer.Write([]byte(tt.writeStr))

			assert.NoError(t, err)
			assert.Equal(t, len(tt.writeStr), n)

			assert.Equal(t, tt.writeStr, string(writer.getBody()))
			assert.Equal(t, tt.writeStr, rec.Body.String())
		})
	}
}
