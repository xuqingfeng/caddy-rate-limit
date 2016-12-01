package ratelimit

import (
	"net"
	"net/http"
)

// IsLocalIpAddress check whether a ip belongs to private network
func IsLocalIpAddress(address string, localIpNets []*net.IPNet) bool {

	ip := net.ParseIP(address)
	if ip != nil {
		for _, ipNet := range localIpNets {
			if ipNet.Contains(ip) {
				return true
			}
		}
	}

	return false
}

// GetRemoteIP returns the ip of requester
// Don't care if the ip is real or not
func GetRemoteIP(r *http.Request) (string, error) {

	// Don't work for mac local
	//remoteAddress := r.RemoteAddr
	//
	//idx := strings.LastIndex(remoteAddress, ":")
	//if idx == -1 {
	//	return remoteAddress
	//}
	//return remoteAddress[:idx]

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return host, err
}
