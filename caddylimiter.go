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
	sync.Mutex
}

func NewCaddyLimiter() *CaddyLimiter {

	return &CaddyLimiter{
		Keys: make(map[string]*rate.Limiter),
	}
}

// Allow is just a shortcut for AllowN
func (cl *CaddyLimiter) Allow(keys []string, rule Rule) bool {

	return cl.AllowN(keys, rule, 1)
}

// AllowN check if n count are allowed for a specific key
func (cl *CaddyLimiter) AllowN(keys []string, rule Rule, n int) bool {

	keysJoined := strings.Join(keys, "|")

	cl.Lock()

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

	cl.Unlock()

	return cl.Keys[keysJoined].AllowN(time.Now(), n)
}

// RetryAfter return a helper message for client
func (cl *CaddyLimiter) RetryAfter(keys []string) time.Duration {

	keysJoined := strings.Join(keys, "|")
	reserve := cl.Keys[keysJoined].Reserve()
	defer reserve.Cancel()

	if reserve.OK() {
		retryAfter := reserve.Delay()
		return retryAfter
	}

	return rate.InfDuration
}

// buildKeys combine client ip and resource
func buildKeys(res string, r *http.Request) [][]string {

	remoteIP, _ := GetRemoteIP(r)
	sliceKeys := make([][]string, 0)

	if len(remoteIP) != 0 {
		sliceKeys = append(sliceKeys, []string{remoteIP, res})
	}

	return sliceKeys
}
