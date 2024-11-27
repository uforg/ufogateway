package gateway

import "net/http"

// cloneHeaderMap converts http.Header to a map[string][]string.
func cloneHeaderMap(h http.Header) map[string][]string {
	cloned := make(map[string][]string, len(h))
	for k, vv := range h {
		if len(vv) == 0 {
			cloned[k] = []string{}
			continue
		}
		vvCopy := append([]string(nil), vv...)
		cloned[k] = vvCopy
	}
	return cloned
}
