package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGatewayToOriginPath(t *testing.T) {
	tests := []struct {
		name        string
		gatewayPath string
		endpoint    string
		want        string
	}{
		{
			name:        "empty paths",
			gatewayPath: "",
			endpoint:    "",
			want:        "",
		},
		{
			name:        "path with leading slashes",
			gatewayPath: "///test/path",
			endpoint:    "//test",
			want:        "path",
		},
		{
			name:        "path without leading slashes",
			gatewayPath: "test/path",
			endpoint:    "test",
			want:        "path",
		},
		{
			name:        "endpoint longer than path",
			gatewayPath: "/api",
			endpoint:    "/api/v1",
			want:        "api",
		},
		{
			name:        "path with multiple segments",
			gatewayPath: "/api/v1/users/123",
			endpoint:    "/api/v1",
			want:        "users/123",
		},
		{
			name:        "exact match path and endpoint",
			gatewayPath: "/api/v1",
			endpoint:    "/api/v1",
			want:        "",
		},
		{
			name:        "partial endpoint match",
			gatewayPath: "/api/v1/test",
			endpoint:    "/api",
			want:        "v1/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gatewayToOriginPath(tt.gatewayPath, tt.endpoint)
			assert.Equal(t, tt.want, got)
		})
	}
}
