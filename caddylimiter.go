package ratelimit

import (
	"net/http"
	"strings"
	"time"

	"github.com/xuqingfeng/caddy-rate-limit/libs/rate"
)

type CaddyLimiter struct {
	Keys map[string]*rate.Limiter
}

func NewCaddyLimiter() *CaddyLimiter {

	return &CaddyLimiter{
		Keys: make(map[string]*rate.Limiter),
	}
}

func (cl *CaddyLimiter) Allow(keys []string, rule Rule) bool {

	return cl.AllowN(keys, rule, 1)
}

func (cl *CaddyLimiter) AllowN(keys []string, rule Rule, n int) bool {

	keysJoined := strings.Join(keys, "|")
	if _, found := cl.Keys[keysJoined]; !found {
		cl.Keys[keysJoined] = rate.NewLimiter(rate.Limit(rule.Rate), rule.Burst, rule.Unit)
	}

	return cl.Keys[keysJoined].AllowN(time.Now(), n)
}

func buildKeys(res string, r *http.Request) [][]string {

	remoteIP := GetRemoteIP(r)
	sliceKeys := make([][]string, 0)

	if len(remoteIP) != 0 {
		sliceKeys = append(sliceKeys, []string{remoteIP, res})
	}

	return sliceKeys
}
