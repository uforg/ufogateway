package strutil

import "strings"

// RemoveAllTrailingSlashes removes all trailing slashes from a string.
func RemoveAllTrailingSlashes(s string) string {
	for {
		s = RemoveTrailingSlash(s)
		if !strings.HasSuffix(s, "/") {
			break
		}
	}
	return s
}
