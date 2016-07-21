package ratelimit

import (
	"net/http"
	"strings"
)

// GetRemoteIP returns the ip of requester
// Doesn't care about the ip is real or not
func GetRemoteIP(r *http.Request) string {

	remoteAddress := r.RemoteAddr

	idx := strings.LastIndex(remoteAddress, ":")
	if idx == -1 {
		return remoteAddress
	}
	return remoteAddress[:idx]
}
