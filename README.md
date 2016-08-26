## caddy-rate-limit
>a `rate limit` plugin for [caddy](https://caddyserver.com/)

[![Travis CI](https://img.shields.io/travis/xuqingfeng/caddy-rate-limit/master.svg?style=flat-square)](https://travis-ci.org/xuqingfeng/caddy-rate-limit)
[![Go Report Card](https://goreportcard.com/badge/github.com/xuqingfeng/caddy-rate-limit?style=flat-square)](https://goreportcard.com/report/github.com/xuqingfeng/caddy-rate-limit)

### Syntax

**Excessive requests will be terminated with an error 429 (Too Many Requests) !**

For single resource:

```
ratelimit path rate burst
```

- path is the file or directory to apply `rate limit`

- rate is the limited request in second (r/s) (eg. 1)

- burst is the maximum burst size client can exceed; burst >= rate (eg. 2)

For multiple resources:

```
ratelimit rate burst {
    resources
}
```

- resources is a list of files/directories to apply `rate limit`, one per line


### Examples

Limit clients to 2 requests per second (bursts of 3) to any resources in /r:

```
ratelimit /r 2 3
```

For the listed paths, limit clients to 2 requests per second (bursts of 2):

```
ratelimit 2 2 {
    /foo.html
    /dir
}
```

#### Inspired By

[http://nginx.org/en/docs/http/ngx_http_limit_req_module.html](http://nginx.org/en/docs/http/ngx_http_limit_req_module.html)

[https://github.com/didip/tollbooth](https://github.com/didip/tollbooth)
