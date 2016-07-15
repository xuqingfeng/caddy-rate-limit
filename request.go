package ratelimit

import (
	"net/http"
	"strings"
)

// GetRemoteIP returns the ip of requester
// priority: X-Forwarded-For > X-Real-Ip > RemoteAddress
func GetRemoteIP(r *http.Request) string {

	remoteAddress := r.RemoteAddr
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	xRealIP := r.Header.Get("X-Real-Ip")

	if len(xForwardedFor) == 0 && len(xRealIP) == 0 {
		idx := strings.LastIndex(remoteAddress, ":")
		if idx == -1 {
			return remoteAddress
		}
		return remoteAddress[:idx]
	} else if len(xForwardedFor) != 0 {
		// list of ip addresses separated by comma
		for _, addr := range strings.Split(xForwardedFor, ",") {
			addr = strings.TrimSpace(addr)
			if len(addr) != 0 {
				return addr
			}
		}
	}

	return xRealIP
}
