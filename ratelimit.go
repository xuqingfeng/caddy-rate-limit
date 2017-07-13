package ratelimit

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// RateLimit is an http.Handler that can limit request rate to specific paths or files
type RateLimit struct {
	Next  httpserver.Handler
	Rules []Rule
}

// Rule is a configuration for ratelimit
type Rule struct {
	Rate      int64
	Burst     int
	Whitelist []string
	Resources []string
	Unit      string
}

const (
	symbol = "^"
)

var (
	caddyLimiter    *CaddyLimiter
	whitelistIpNets []*net.IPNet
)

func init() {

	caddyLimiter = NewCaddyLimiter()
}

func (rl RateLimit) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	retryAfter := time.Duration(0)

	// handle exception first
	for _, rule := range rl.Rules {
		for _, res := range rule.Resources {
			if strings.HasPrefix(res, symbol) {
				res = strings.TrimPrefix(res, symbol)
				if httpserver.Path(r.URL.Path).Matches(res) {
					return rl.Next.ServeHTTP(w, r)
				}
			}
		}
		// get whitelist IPNet
		for _, s := range rule.Whitelist {
			_, ipNet, err := net.ParseCIDR(s)
			if err == nil {
				whitelistIpNets = append(whitelistIpNets, ipNet)
			}
		}
	}

	for _, rule := range rl.Rules {
		for _, res := range rule.Resources {
			if !httpserver.Path(r.URL.Path).Matches(res) {
				continue
			}

			// filter whitelist ips
			address, err := GetRemoteIP(r)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			if IsWhitelistIpAddress(address, whitelistIpNets) {
				continue
			}

			sliceKeys := buildKeys(res, r)
			for _, keys := range sliceKeys {
				ret := caddyLimiter.Allow(keys, rule)
				retryAfter = caddyLimiter.RetryAfter(keys)
				if !ret {
					w.Header().Add("X-RateLimit-RetryAfter", retryAfter.String())
					return http.StatusTooManyRequests, nil
				}
			}
		}
	}

	return rl.Next.ServeHTTP(w, r)
}
