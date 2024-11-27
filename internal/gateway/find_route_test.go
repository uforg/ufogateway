package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindRoute(t *testing.T) {
	tests := []struct {
		name     string
		routes   []Route
		path     string
		want     Route
		wantFind bool
	}{
		{
			name: "exact match",
			routes: []Route{
				{Endpoint: "/api/users"},
				{Endpoint: "/api/posts"},
			},
			path:     "/api/users",
			want:     Route{Endpoint: "/api/users"},
			wantFind: true,
		},
		{
			name: "prefix match",
			routes: []Route{
				{Endpoint: "/api"},
				{Endpoint: "/web"},
			},
			path:     "/api/users",
			want:     Route{Endpoint: "/api"},
			wantFind: true,
		},
		{
			name:     "empty routes",
			routes:   []Route{},
			path:     "/api/users",
			want:     Route{},
			wantFind: false,
		},
		{
			name: "no matching route",
			routes: []Route{
				{Endpoint: "/api"},
				{Endpoint: "/web"},
			},
			path:     "/admin",
			want:     Route{},
			wantFind: false,
		},
		{
			name: "multiple routes with prefix",
			routes: []Route{
				{Endpoint: "/api"},
				{Endpoint: "/api/v1"},
				{Endpoint: "/api/v2"},
			},
			path:     "/api/v1/users",
			want:     Route{Endpoint: "/api/v1"},
			wantFind: true,
		},
		{
			name: "empty path",
			routes: []Route{
				{Endpoint: "/api"},
				{Endpoint: ""},
			},
			path:     "",
			want:     Route{Endpoint: ""},
			wantFind: true,
		},
		{
			name: "root path",
			routes: []Route{
				{Endpoint: "/"},
				{Endpoint: "/api"},
			},
			path:     "/",
			want:     Route{Endpoint: "/"},
			wantFind: true,
		},
		{
			name: "case sensitive",
			routes: []Route{
				{Endpoint: "/API"},
				{Endpoint: "/api"},
			},
			path:     "/api/users",
			want:     Route{Endpoint: "/api"},
			wantFind: true,
		},
		{
			name: "deeply nested paths",
			routes: []Route{
				{Endpoint: "/api"},
				{Endpoint: "/api/v1"},
				{Endpoint: "/api/v1/users"},
				{Endpoint: "/api/v1/users/profiles"},
			},
			path:     "/api/v1/users/profiles/1234/settings",
			want:     Route{Endpoint: "/api/v1/users/profiles"},
			wantFind: true,
		},
		{
			name: "paths with special characters",
			routes: []Route{
				{Endpoint: "/api/users-data"},
				{Endpoint: "/api/users_data"},
				{Endpoint: "/api/users.data"},
			},
			path:     "/api/users-data/details",
			want:     Route{Endpoint: "/api/users-data"},
			wantFind: true,
		},
		{
			name: "paths with numbers",
			routes: []Route{
				{Endpoint: "/api/v1"},
				{Endpoint: "/api/v2"},
				{Endpoint: "/api/v10"},
			},
			path:     "/api/v10/users",
			want:     Route{Endpoint: "/api/v10"},
			wantFind: true,
		},
		{
			name: "conflicting nested paths",
			routes: []Route{
				{Endpoint: "/api/users"},
				{Endpoint: "/api/users/v1"},
				{Endpoint: "/api/users-v1"},
			},
			path:     "/api/users/v1/profiles",
			want:     Route{Endpoint: "/api/users/v1"},
			wantFind: true,
		},
		{
			name: "trailing slashes",
			routes: []Route{
				{Endpoint: "/api/"},
				{Endpoint: "/api/v1/"},
				{Endpoint: "/api/v1/users/"},
			},
			path:     "/api/v1/users/123",
			want:     Route{Endpoint: "/api/v1/users/"},
			wantFind: true,
		},
		{
			name: "mixed trailing slashes",
			routes: []Route{
				{Endpoint: "/api"},
				{Endpoint: "/api/v1/"},
			},
			path:     "/api/v1/users",
			want:     Route{Endpoint: "/api/v1/"},
			wantFind: true,
		},
		{
			name: "exact match over prefix",
			routes: []Route{
				{Endpoint: "/api/users"},
				{Endpoint: "/api/users/"},
				{Endpoint: "/api"},
			},
			path:     "/api/users",
			want:     Route{Endpoint: "/api/users"},
			wantFind: true,
		},
		{
			name: "root path with many options",
			routes: []Route{
				{Endpoint: "/"},
				{Endpoint: ""},
				{Endpoint: "/api"},
			},
			path:     "/any/path/here",
			want:     Route{Endpoint: "/"},
			wantFind: true,
		},
		{
			name: "multiple segments with same prefix",
			routes: []Route{
				{Endpoint: "/api/v1/users"},
				{Endpoint: "/api/v1/users-groups"},
				{Endpoint: "/api/v1/users-permissions"},
			},
			path:     "/api/v1/users-groups/123",
			want:     Route{Endpoint: "/api/v1/users-groups"},
			wantFind: true,
		},
		{
			name: "unicode paths",
			routes: []Route{
				{Endpoint: "/api/español"},
				{Endpoint: "/api/中文"},
				{Endpoint: "/api/русский"},
			},
			path:     "/api/español/usuarios",
			want:     Route{Endpoint: "/api/español"},
			wantFind: true,
		},
		{
			name: "very long paths",
			routes: []Route{
				{Endpoint: "/api/v1/very/long/path/with/many/segments"},
				{Endpoint: "/api/v1/very/long/path/with"},
				{Endpoint: "/api/v1/very"},
			},
			path:     "/api/v1/very/long/path/with/many/segments/and/more",
			want:     Route{Endpoint: "/api/v1/very/long/path/with/many/segments"},
			wantFind: true,
		},
		{
			name: "repeated segments",
			routes: []Route{
				{Endpoint: "/api/users"},
				{Endpoint: "/api/users/users"},
				{Endpoint: "/api/users/users/users"},
			},
			path:     "/api/users/users/users/details",
			want:     Route{Endpoint: "/api/users/users/users"},
			wantFind: true,
		},
		{
			name: "path parameter patterns",
			routes: []Route{
				{Endpoint: "/api/users/:id"},
				{Endpoint: "/api/users/profile"},
				{Endpoint: "/api/users/:id/profile"},
			},
			path:     "/api/users/:id/profile/settings",
			want:     Route{Endpoint: "/api/users/:id/profile"},
			wantFind: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := findRoute(tt.routes, tt.path)
			assert.Equal(t, tt.wantFind, found)
			assert.Equal(t, tt.want, got)
		})
	}
}
