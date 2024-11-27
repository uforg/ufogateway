package gateway

import (
	"strings"

	"github.com/uforg/ufogateway/internal/util/strutil"
)

// gatewayToOriginPath transforms the path of a request made
// to the gateway to the path of the request that will be made
// to the origin server.
func gatewayToOriginPath(gatewayPath, endpoint string) string {
	cleanPath := strutil.RemoveAllLeadingSlashes(gatewayPath)
	cleanEndpoint := strutil.RemoveAllLeadingSlashes(endpoint)
	return strutil.RemoveAllLeadingSlashes(strings.TrimPrefix(cleanPath, cleanEndpoint))
}
