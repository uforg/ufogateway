package gateway

import (
	"fmt"
	"net/http"

	"github.com/uforg/ufogateway/internal/util/strutil"
)

// getRequestURL returns the URL for the incoming request given an http.Request and Route.
//
// Returns:
//   - string: The URL of the request to the gateway.
//   - string: The URL of the request to the origin server.
func getRequestURL(req *http.Request, route Route) (string, string) {
	schemeStr := "http"
	if req.URL.Scheme != "" {
		schemeStr = req.URL.Scheme
	}

	gatewayPath := strutil.RemoveAllLeadingSlashes(req.URL.Path)
	originPath := gatewayToOriginPath(gatewayPath, route.Endpoint)

	queryStr := req.URL.RawQuery
	hasQuery := queryStr != ""

	fragmentStr := req.URL.Fragment
	hasFragment := fragmentStr != ""

	gatewayURL := fmt.Sprintf("%s://%s", schemeStr, req.Host)
	if gatewayPath != "" {
		gatewayURL = strutil.RemoveAllTrailingSlashes(gatewayURL) + "/" + gatewayPath
	}
	if hasQuery {
		gatewayURL += "?" + queryStr
	}
	if hasFragment {
		gatewayURL += "#" + fragmentStr
	}

	originURL := route.OriginURL
	if originPath != "" {
		originURL = strutil.RemoveAllTrailingSlashes(originURL) + "/" + originPath
	}
	if hasQuery {
		originURL += "?" + queryStr
	}
	if hasFragment {
		originURL += "#" + fragmentStr
	}

	return gatewayURL, originURL
}
