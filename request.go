package ratelimit

import (
	"net"
	"net/http"
)

// IsWhitelistIPAddress check whether an ip is in whitelist
func IsWhitelistIPAddress(address string, localIPNets []*net.IPNet) bool {

	ip := net.ParseIP(address)
	if ip != nil {
		for _, ipNet := range localIPNets {
			if ipNet.Contains(ip) {
				return true
			}
		}
	}

	return false
}

// GetRemoteIP returns the ip of requester
// Doesn't care if the ip is real or not
func GetRemoteIP(r *http.Request) (string, error) {

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return host, err
}
