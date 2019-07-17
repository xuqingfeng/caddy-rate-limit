package main

import (
	"github.com/caddyserver/caddy/caddy/caddymain"

	_ "github.com/xuqingfeng/caddy-rate-limit"
)

func main() {
	caddymain.EnableTelemetry = false
	caddymain.Run()
}
