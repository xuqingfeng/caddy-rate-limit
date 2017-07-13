package ratelimit

import (
	"net"
	"net/http"
)

// IsWhitelistIpAddress check whether an ip is in whitelist
func IsWhitelistIpAddress(address string, localIpNets []*net.IPNet) bool {

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

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return host, err
}
