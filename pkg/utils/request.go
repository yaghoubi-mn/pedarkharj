package utils

import "net/http"

// returns "" if ip not found
func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}

func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}
