## caddy-rate-limit
>a `rate limit` plugin for [caddy](https://caddyserver.com/)

[![Travis CI](https://img.shields.io/travis/xuqingfeng/caddy-rate-limit/master.svg?style=flat-square)](https://travis-ci.org/xuqingfeng/caddy-rate-limit)
[![Go Report Card](https://goreportcard.com/badge/github.com/xuqingfeng/caddy-rate-limit?style=flat-square)](https://goreportcard.com/report/github.com/xuqingfeng/caddy-rate-limit)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/xuqingfeng/caddy-rate-limit)

### Syntax

**Excessive requests will be terminated with an error 429 (Too Many Requests)! And `X-RateLimit-RetryAfter` header will be returned.**

For single resource:

```
ratelimit methods path rate burst unit
```

- methods are the requests methods it tries to match (comma separately)

- path is the file or directory to apply `rate limit`

- rate is the limited request in every time unit (r/s, r/m, r/h) (e.g. 1)

- burst is the maximum burst size client can exceed; burst >= rate (e.g. 2)

- unit is the time interval (currently support: `second`, `minute`, `hour`)

For multiple resources:

```
ratelimit methods rate burst unit {
    whitelist CIDR
    resources
}
```

- whitelist is the keyword for whitelist your trusted ips, [CIDR](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) is the IP range you don't want to perform `rate limit` and it will apply to all rules.
- resources is a list of files/directories to apply `rate limit`, one per line.

**Note:** If you don't want to apply `rate limit` on some special resources, add `^` in front of the path.


### Examples

Limit clients to 2 requests per second (bursts of 3) to any methods and any resources in /r:

```
ratelimit * /r 2 3 second
```

Don't perform `rate limit` if requests come from **1.2.3.4** or **192.168.1.0/30(192.168.1.0 ~ 192.168.1.3)**, for the listed paths, limit clients to 2 requests per minute (bursts of 2) if the request method is **GET** or **POST** and always ignore `/dir/app.js`:

```
ratelimit get,post 2 2 minute {
    whitelist 1.2.3.4/32
    whitelist 192.168.1.0/30
    /foo.html
    /api
    ^/dir/app.js
}
```

### Download

`curl https://getcaddy.com | bash -s personal http.ratelimit`

### Docker

```bash
docker run -d -p 2016:2016 -v `pwd`/Caddyfile:/go/src/github.com/xuqingfeng/caddy-rate-limit/Caddyfile --name ratelimit xuqingfeng/caddy-rate-limit
```

---

**Inspired by**

[http://nginx.org/en/docs/http/ngx_http_limit_req_module.html](http://nginx.org/en/docs/http/ngx_http_limit_req_module.html)

[https://github.com/didip/tollbooth](https://github.com/didip/tollbooth)
