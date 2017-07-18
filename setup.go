package ratelimit

import (
	"net"
	"strconv"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {

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
		case 3:
			// config block
			rule.Rate, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return rules, err
			}
			rule.Burst, err = strconv.Atoi(args[1])
			if err != nil {
				return rules, err
			}
			rule.Unit = args[2]
		case 4:
			// one line config
			rule.Resources = append(rule.Resources, args[0])
			rule.Rate, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return rules, err
			}
			rule.Burst, err = strconv.Atoi(args[2])
			if err != nil {
				return rules, err
			}
			rule.Unit = args[3]
		default:
			return rules, c.ArgErr()
		}

		for c.NextBlock() {
			val := c.Val()
			args = c.RemainingArgs()
			switch len(args) {
			case 0:
				// resource
				rule.Resources = append(rule.Resources, val)
			case 1:
				// whitelist
				if "whitelist" == val {
					// check if CIDR is valid
					_, _, err := net.ParseCIDR(args[0])
					if err != nil {
						return rules, err
					} else {
						rule.Whitelist = append(rule.Whitelist, args[0])
					}
				} else {
					return rules, c.Errf("expecting whitelist, got %s", val)
				}
			default:
				return rules, c.ArgErr()
			}
		}

		rules = append(rules, rule)
	}

	return rules, nil
}
