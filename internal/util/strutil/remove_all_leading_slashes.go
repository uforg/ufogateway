package strutil

import "strings"

// RemoveAllLeadingSlashes removes all leading slashes from the input string.
func RemoveAllLeadingSlashes(s string) string {
	for {
		s = RemoveLeadingSlash(s)
		if !strings.HasPrefix(s, "/") {
			break
		}
	}

	return s
}
