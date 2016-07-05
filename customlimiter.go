package ratelimit

import (
	"strings"

	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type CustomLimiter struct {
	Keys map[string]*rate.Limiter
}

func NewCustomLimiter() *CustomLimiter {

	customLimiter := CustomLimiter{}
	customLimiter.Keys = make(map[string]*rate.Limiter)
	return &customLimiter
}

func (c *CustomLimiter) Allow(keys []string, rule Rule) bool {

	keysJoined := strings.Join(keys, "|")
	if _, found := c.Keys[keysJoined]; !found {
		c.Keys[keysJoined] = rate.NewLimiter(rate.Limit(rule.Rate), rule.Burst)
	}

	return c.Keys[keysJoined].Allow()
}

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

	if "" == remoteIP {
		return sliceKeys
	} else {
		sliceKeys = append(sliceKeys, []string{remoteIP, res})
	}

	return sliceKeys
}
