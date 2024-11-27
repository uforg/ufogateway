package gateway

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetRequestURL(t *testing.T) {
	table := []struct {
		name           string
		req            *http.Request
		route          Route
		wantGatewayURL string
		wantOriginURL  string
	}{
		{
			name: "simple request",
			req: &http.Request{
				Host: "localhost:8080",
				URL: &url.URL{
					Path: "/example",
				},
			},
			route: Route{
				Endpoint:  "/example",
				OriginURL: "https://example.com",
			},
			wantGatewayURL: "http://localhost:8080/example",
			wantOriginURL:  "https://example.com",
		},
		{
			name: "request with query parameters",
			req: &http.Request{
				Host: "localhost:8080",
				URL: &url.URL{
					Path:     "/api/users",
					RawQuery: "page=1&limit=10",
				},
			},
			route: Route{
				Endpoint:  "/api",
				OriginURL: "https://api.example.com",
			},
			wantGatewayURL: "http://localhost:8080/api/users?page=1&limit=10",
			wantOriginURL:  "https://api.example.com/users?page=1&limit=10",
		},
		{
			name: "request with URL fragment",
			req: &http.Request{
				Host: "localhost:8080",
				URL: &url.URL{
					Path:     "/docs",
					Fragment: "section-1",
				},
			},
			route: Route{
				Endpoint:  "/docs",
				OriginURL: "https://docs.example.com",
			},
			wantGatewayURL: "http://localhost:8080/docs#section-1",
			wantOriginURL:  "https://docs.example.com#section-1",
		},
		{
			name: "HTTPS scheme",
			req: &http.Request{
				Host: "secure.gateway.com",
				URL: &url.URL{
					Scheme: "https",
					Path:   "/secure/data",
				},
			},
			route: Route{
				Endpoint:  "/secure",
				OriginURL: "https://internal.example.com",
			},
			wantGatewayURL: "https://secure.gateway.com/secure/data",
			wantOriginURL:  "https://internal.example.com/data",
		},
		{
			name: "empty path",
			req: &http.Request{
				Host: "localhost:8080",
				URL: &url.URL{
					Path: "/",
				},
			},
			route: Route{
				Endpoint:  "/",
				OriginURL: "https://example.com",
			},
			wantGatewayURL: "http://localhost:8080",
			wantOriginURL:  "https://example.com",
		},
		{
			name: "keep origin trailing slash",
			req: &http.Request{
				Host: "localhost:8080",
				URL: &url.URL{
					Path: "/example",
				},
			},
			route: Route{
				Endpoint:  "/example",
				OriginURL: "https://example.com/",
			},
			wantGatewayURL: "http://localhost:8080/example",
			wantOriginURL:  "https://example.com/",
		},
		{
			name: "complex request with query and fragment",
			req: &http.Request{
				Host: "gateway.example.com",
				URL: &url.URL{
					Path:     "/api/v1/users/123/profile",
					RawQuery: "format=json&fields=name,email",
					Fragment: "personal-info",
				},
			},
			route: Route{
				Endpoint:  "/api/v1",
				OriginURL: "https://api.internal.com/v1",
			},
			wantGatewayURL: "http://gateway.example.com/api/v1/users/123/profile?format=json&fields=name,email#personal-info",
			wantOriginURL:  "https://api.internal.com/v1/users/123/profile?format=json&fields=name,email#personal-info",
		},
		{
			name: "destination URL with trailing slash",
			req: &http.Request{
				Host: "localhost:8080",
				URL: &url.URL{
					Path: "/api/resource",
				},
			},
			route: Route{
				Endpoint:  "/api",
				OriginURL: "https://api.example.com/",
			},
			wantGatewayURL: "http://localhost:8080/api/resource",
			wantOriginURL:  "https://api.example.com/resource",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotGateway, gotOrigin := getRequestURL(tt.req, tt.route)
			if gotGateway != tt.wantGatewayURL {
				t.Errorf("getRequestURL() gateway URL got = %v, want %v", gotGateway, tt.wantGatewayURL)
			}
			if gotOrigin != tt.wantOriginURL {
				t.Errorf("getRequestURL() origin URL got = %v, want %v", gotOrigin, tt.wantOriginURL)
			}
		})
	}
}
