package ratelimit

import (
	"log"
	"os"
    "strconv"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// todo: rm
var logger *log.Logger

func init() {

	logger = log.New(os.Stderr, "ratelimit: ", log.LstdFlags)
	caddy.RegisterPlugin("ratelimit", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {

	cfg := httpserver.GetConfig(c)

	rules, err := rateLimitParse(c)
	if err != nil {
		return err
	}

	rateLimit := RateLimit{Rules: rules}
	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		rateLimit.Next = next
		return rateLimit
	})

	return nil
}

func rateLimitParse(c *caddy.Controller) (rules []Rule, err error) {

	for c.Next() {
		var rule Rule

		args := c.RemainingArgs()
		switch len(args) {
		case 2:
			rule.Rate, err = strconv.ParseFloat(args[0], 64)
			if err != nil {
				return rules, err
			}
			rule.Burst, err = strconv.Atoi(args[1])
			if err != nil {
				return rules, err
			}
			for c.NextBlock() {
				rule.Resources = append(rule.Resources, c.Val())
				if c.NextArg() {
					return rules, c.Errf("Expecting only one resource per line (extra '%s')", c.Val())
				}
			}
		case 3:
			rule.Resources = append(rule.Resources, args[0])
			rule.Rate, err = strconv.ParseFloat(args[1], 64)
			if err != nil {
				return rules, err
			}
			rule.Burst, err = strconv.Atoi(args[2])
			if err != nil {
				return rules, err
			}
		default:
			return rules, c.ArgErr()
		}

		rules = append(rules, rule)
	}

	// ? no return; error
	return rules, nil
}
