package strutil

// RemoveLeadingSlash removes the leading slash from the input string.
func RemoveLeadingSlash(s string) string {
	if len(s) == 0 {
		return s
	}

	if s[0] == '/' {
		return s[1:]
	}

	return s
}
