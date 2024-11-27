package gateway

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// getRequestIP extracts the client's IP address from an HTTP request.
func getRequestIP(r *http.Request) (string, error) {
	// Check X-Forwarded-For header
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		splitIps := strings.Split(ips, ",")
		// Take the first IP (original client IP)
		clientIP := strings.TrimSpace(splitIps[0])
		netIP := net.ParseIP(clientIP)
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	// Fallback to RemoteAddr if no valid IP is found in X-Forwarded-For
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		// Normalize loopback address for consistency
		if netIP.IsLoopback() {
			return "127.0.0.1", nil
		}
		return netIP.String(), nil
	}

	return "", errors.New("IP not found")
}
