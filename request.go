package ratelimit

import (
	"net/http"
	"strings"
)

func GetRemoteIP(r *http.Request) string {

	// remote address
	remoteAddress := r.RemoteAddr
	// x-forwarded-for
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	// real ip
	xRealIP := r.Header.Get("X-Real-Ip")

	if len(remoteAddress) != 0 {
		idx := strings.LastIndex(remoteAddress, ":")
		if -1 == idx {
			return remoteAddress
		}
		return remoteAddress[:idx]
	} else if len(xForwardedFor) != 0 {
		// ip list separated by comma
		// get first match ?
		for _, addr := range strings.Split(xForwardedFor, ",") {
			addr = strings.TrimSpace(addr)
			if len(addr) != 0 {
				return addr
			}
		}
		return xForwardedFor
	} else if len(xRealIP) != 0 {

		return xRealIP
	}

	return ""
}
