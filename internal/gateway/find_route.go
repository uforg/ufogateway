package gateway

import "strings"

// findRoute finds the appropriate route for a given path.
// It returns the matching route and a boolean indicating whether a match was found.
// When multiple routes match, it returns the most specific (longest) matching route.
func findRoute(routes []Route, path string) (Route, bool) {
	var bestMatch Route
	var bestLength int
	var found bool

	for _, route := range routes {
		if strings.HasPrefix(path, route.Endpoint) {
			currentLength := len(route.Endpoint)
			if !found || currentLength > bestLength {
				bestMatch = route
				bestLength = currentLength
				found = true
			}
		}
	}

	return bestMatch, found
}
