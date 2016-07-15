package ratelimit

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// CustomLimiter holds a map of *rate.Limiter
type CustomLimiter struct {
	Keys map[string]*rate.Limiter
}

// NewCustomLimiter returns a new *CustomLimiter
func NewCustomLimiter() *CustomLimiter {

	customLimiter := CustomLimiter{
		Keys: make(map[string]*rate.Limiter),
	}
	return &customLimiter
}

// Allow is shorthand for AllowN(keys, rule, 1)
func (c *CustomLimiter) Allow(keys []string, rule Rule) bool {

	return c.AllowN(keys, rule, 1)
}

// AllowN check whether keys is allowed to pass based on rule and n
func (c *CustomLimiter) AllowN(keys []string, rule Rule, n int) bool {

	keysJoined := strings.Join(keys, "|")
	if _, found := c.Keys[keysJoined]; !found {
		c.Keys[keysJoined] = rate.NewLimiter(rate.Limit(rule.Rate), rule.Burst)
	}

	return c.Keys[keysJoined].AllowN(time.Now(), n)
}

func buildKeys(res string, r *http.Request) [][]string {

	remoteIP := GetRemoteIP(r)
	sliceKeys := make([][]string, 0)

	if len(remoteIP) != 0 {
		sliceKeys = append(sliceKeys, []string{remoteIP, res})
	}

	return sliceKeys
}
