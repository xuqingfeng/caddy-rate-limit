package ratelimit

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type CaddyLimiter struct {
	Keys map[string]*rate.Limiter
	mu   sync.Mutex
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
	cl.mu.Lock()
	if _, found := cl.Keys[keysJoined]; !found {

		switch rule.Unit {
		case "second":
			cl.Keys[keysJoined] = rate.NewLimiter(rate.Limit(rule.Rate)/rate.Limit(time.Second.Seconds()), rule.Burst)
		case "minute":
			cl.Keys[keysJoined] = rate.NewLimiter(rate.Limit(rule.Rate)/rate.Limit(time.Minute.Seconds()), rule.Burst)
		case "hour":
			cl.Keys[keysJoined] = rate.NewLimiter(rate.Limit(rule.Rate)/rate.Limit(time.Hour.Seconds()), rule.Burst)
		default:
			// Infinite
			cl.Keys[keysJoined] = rate.NewLimiter(rate.Inf, rule.Burst)
		}
	}
	cl.mu.Unlock()

	return cl.Keys[keysJoined].AllowN(time.Now(), n)
}

func (cl *CaddyLimiter) RetryAfter(keys []string) time.Duration {

	keysJoined := strings.Join(keys, "|")
	reserve := cl.Keys[keysJoined].Reserve()
	if reserve.OK() {
		retryAfter := reserve.Delay()
		reserve.Cancel()
		return retryAfter
	}

	reserve.Cancel()
	return rate.InfDuration
}

func buildKeys(res string, r *http.Request) [][]string {

	remoteIP, _ := GetRemoteIP(r)
	sliceKeys := make([][]string, 0)

	if len(remoteIP) != 0 {
		sliceKeys = append(sliceKeys, []string{remoteIP, res})
	}

	return sliceKeys
}
