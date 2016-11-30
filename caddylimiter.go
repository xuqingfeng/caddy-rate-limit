package ratelimit

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
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
		switch rule.Unit {
		case "second":
			cl.Keys[keysJoined] = rate.NewLimiter(rate.Every(time.Second), rule.Burst)
		case "minute":
			cl.Keys[keysJoined] = rate.NewLimiter(rate.Every(time.Minute), rule.Burst)
		case "hour":
			cl.Keys[keysJoined] = rate.NewLimiter(rate.Every(time.Hour), rule.Burst)
		default:
			// Infinite
			cl.Keys[keysJoined] = rate.NewLimiter(rate.Every(0), rule.Burst)
		}
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
