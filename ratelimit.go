package ratelimit

import (
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	"net"
)

// RateLimit is an http.Handler that can limit request rate to specific paths or files
type RateLimit struct {
	Next  httpserver.Handler
	Rules []Rule
}

// Rule is a configuration for ratelimit
type Rule struct {
	Rate      float64
	Burst     int
	Resources []string
}

var (
	customLimiter *CustomLimiter
	localIpNets   []*net.IPNet
)

func init() {

	customLimiter = NewCustomLimiter()
	localCIDRs := []string{
		"127.0.0.1/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
	}
	for _, s := range localCIDRs {
		_, ipNet, err := net.ParseCIDR(s)
		if err == nil {
			localIpNets = append(localIpNets, ipNet)
		}
	}
}

func (rl RateLimit) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	for _, rule := range rl.Rules {
		for _, res := range rule.Resources {
			if !httpserver.Path(r.URL.Path).Matches(res) {
				continue
			}

			// filter local ip address
			if IsLocalIpAddress(r.RemoteAddr, localIpNets) {
				continue
			}

			sliceKeys := buildKeys(res, r)
			for _, keys := range sliceKeys {
				ret := customLimiter.Allow(keys, rule)
				if !ret {
					return http.StatusTooManyRequests, nil
				}
			}
		}
	}

	return rl.Next.ServeHTTP(w, r)
}
