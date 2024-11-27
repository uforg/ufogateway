package strutil

// RemoveTrailingSlash removes the trailing slash from a string.
func RemoveTrailingSlash(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[len(s)-1] == '/' {
		return s[:len(s)-1]
	}
	return s
}
