package utils

import (
	"net/http"
	"strings"
)

// returns "" if ip not found
func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	// check localhost
	if ip[:5] == "[::1]" {
		ip = "127.0.0.1"
	}

	// remove port from ip
	ip = strings.Split(ip, ":")[0]

	return ip
}

func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}
