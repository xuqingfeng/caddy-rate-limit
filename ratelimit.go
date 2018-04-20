package ratelimit

import (
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
	Methods   string
	Rate      int64
	Burst     int
	Unit      string
	Whitelist []string
	Resources []string
}

const (
	ignoreSymbol = "^"
)

var (
	caddyLimiter *CaddyLimiter
)

func init() {

	caddyLimiter = NewCaddyLimiter()
}

// ServeHTTP is the method handling every request
func (rl RateLimit) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	retryAfter := time.Duration(0)
	// get request ip address
	ipAddress, err := GetRemoteIP(r)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, rule := range rl.Rules {
		for _, res := range rule.Resources {

			// handle exception first
			if strings.HasPrefix(res, ignoreSymbol) {
				res = strings.TrimPrefix(res, ignoreSymbol)
				if httpserver.Path(r.URL.Path).Matches(res) {
					return rl.Next.ServeHTTP(w, r)
				}
			}

			if !httpserver.Path(r.URL.Path).Matches(res) {
				continue
			}

			// whitelist will apply to all rules
			if IsWhitelistIPAddress(ipAddress, whitelistIPNets) || !MatchMethod(rule.Methods, r.Method) {
				continue
			}

			sliceKeys := buildKeys(ipAddress, rule.Methods, res, r)
			for _, keys := range sliceKeys {
				ret := caddyLimiter.Allow(keys, rule)
				if !ret {
					retryAfter = caddyLimiter.RetryAfter(keys)
					w.Header().Add("X-RateLimit-RetryAfter", retryAfter.String())
					return http.StatusTooManyRequests, nil
				}
			}
		}
	}

	return rl.Next.ServeHTTP(w, r)
}
